package router

import (
	"time"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/gofiber/template/html/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rrune/alleytrack/pkg/data"
	"github.com/rrune/alleytrack/pkg/models"
)

type routes struct {
	Alleycat *models.Alleycat
	DB       data.Data
}

func Init(alleycat *models.Alleycat, data data.Data) {
	r := routes{alleycat, data}

	engine := html.New("./web/templates", ".html")
	engine.AddFunc(
		"formatDate", func(timestamp time.Time) string {
			loc, _ := time.LoadLocation("Europe/Berlin")
			return timestamp.In(loc).Format("15:04:05")
		},
	)
	engine.Reload(true)

	app := fiber.New(fiber.Config{
		Views:       engine,
		ProxyHeader: fiber.HeaderXForwardedFor,
	})

	// jwt middleware for participants
	pAuth := jwtware.New(jwtware.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Redirect("/login?path=" + c.Path())
		},
		TokenLookup:   "cookie:JWT",
		SigningKey:    []byte(r.Alleycat.Config.JwtKey),
		SigningMethod: "HS256",
	})

	// jwt middleware for admin
	aAuth := jwtware.New(jwtware.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Redirect("/adminlogin?path=" + c.Path())
		},
		TokenLookup:   "cookie:JWT",
		SigningKey:    []byte(r.Alleycat.Config.JwtKey),
		SigningMethod: "HS256",
	})

	app.Static("/static", "./web/public/static")

	// views
	app.Get("/", r.Index)
	app.Get("/signup", r.SignUp)
	app.Get("/login", r.Login)
	app.Get("/logout", r.HandleLogout)

	// api
	app.Get("/isNumberTaken/:number", r.IsNumerTaken)
	app.Post("/signup", r.HandleSignUp)
	app.Post("/login", r.HandleLogin)

	// admin
	app.Get("/adminlogin", r.AdminLogin)
	app.Post("/adminlogin", r.HandleAdminLogin)
	admin := app.Group("/admin", aAuth)
	admin.Use(func(c *fiber.Ctx) error {
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)

		if !claims["admin"].(bool) {
			c.SendStatus(fiber.StatusUnauthorized)
		}
		return c.Next()
	})
	admin.Get("/", r.Admin)
	admin.Get("/switchEnabled", r.HandleSwitchEnabled)
	admin.Get("/participant/:number", r.Participant)
	admin.Post("/participant/:number", r.HandleParticipant)
	admin.Get("/removeCheckpoint/:number/:checkpoint", r.HandleRemoveCheckpoint)
	admin.Get("/removeParticipant/:number", r.HandleRemoveParticipant)

	// manifest
	manifest := app.Group("/manifest", pAuth)
	manifest.Get("/", r.Manifest)

	// checkpoint endpoints
	for _, c := range r.Alleycat.Manifest {
		if c.Text {
			app.Get("/"+c.Link, pAuth, r.TextCheckpoint)
			app.Post("/"+c.Link, pAuth, r.CompleteTextCheckpoint)
		} else {
			app.Get("/"+c.Link, pAuth, r.CompleteSimpleCheckpoint)
		}
	}

	app.Listen(":" + alleycat.Config.Port)
}

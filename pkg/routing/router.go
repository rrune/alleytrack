package router

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/gofiber/template/html/v2"
	"github.com/rrune/alleytrack/pkg/data"
	"github.com/rrune/alleytrack/pkg/models"
)

type routes struct {
	Config models.Config
	DB     data.Data
}

func Init(conf models.Config, data data.Data) {
	r := routes{conf, data}

	engine := html.New("./web/templates", ".html")
	engine.Reload(true)

	app := fiber.New(fiber.Config{
		Views:       engine,
		ProxyHeader: fiber.HeaderXForwardedFor,
	})

	auth := jwtware.New(jwtware.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Redirect("/login?path=" + c.Path())
		},
		TokenLookup:   "cookie:JWT",
		SigningKey:    []byte(r.Config.JwtKey),
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

	for _, c := range r.Config.Manifest {
		if c.Text {
			app.Get("/"+c.Link, auth, r.TextCheckpoint)
			app.Post("/"+c.Link, auth, r.CompleteTextCheckpoint)
		} else {
			app.Get("/"+c.Link, auth, r.CompleteSimpleCheckpoint)
		}
	}

	manifest := app.Group("/manifest", auth)
	manifest.Get("/", r.Manifest)

	app.Listen(":" + conf.Port)
}

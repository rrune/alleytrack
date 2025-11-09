package router

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rrune/alleytrack/internal/util"
)

func (r routes) Index(c *fiber.Ctx) error {
	jwtCookie := c.Cookies("JWT", "")
	if jwtCookie != "" {
		return c.Redirect("/manifest")
	}

	return c.Render("index", fiber.Map{
		"Title":       r.Alleycat.Config.LogoText,
		"CSS":         "index",
		"LogoText":    r.Alleycat.Config.LogoText,
		"WelcomeText": r.Alleycat.WelcomeText,
	})
}

func (r routes) SignUp(c *fiber.Ctx) error {
	return c.Render("signup", fiber.Map{
		"Title":    "Sign Up",
		"CSS":      "signup",
		"LogoText": r.Alleycat.Config.LogoText,
	})
}

func (r routes) Login(c *fiber.Ctx) error {
	return c.Render("login", fiber.Map{
		"Title":    "Login",
		"CSS":      "signup",
		"LogoText": r.Alleycat.Config.LogoText,
		"Msg":      c.Query("msg", ""),
	})
}

func (r routes) Manifest(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	number := claims["number"].(string)
	admin := claims["admin"].(bool)
	if admin {
		return c.Redirect("/admin")
	}

	num, err := strconv.Atoi(number)
	if util.CheckWLogs(err) {
		return c.SendStatus(500)
	}

	p, exist, err := r.DB.Participants.GetByNumber(num)
	if !exist {
		return c.Redirect("/logout")
	}
	if util.CheckWLogs(err) {
		return c.SendStatus(500)
	}

	completed, err := r.DB.ParticipantsCheckpoints.GetCompleted(num)
	if util.CheckWLogs(err) {
		return c.SendStatus(500)
	}

	unlocked, err := r.DB.CheckpointDependencies.GetAvailableByNumber(num)
	if util.CheckWLogs(err) {
		return c.SendStatus(500)
	}

	return c.Render("manifest", fiber.Map{
		"Title":       "Manifest",
		"CSS":         "manifest",
		"LogoText":    "Manifest",
		"Name":        p.Name,
		"Number":      number,
		"Completed":   completed,
		"Unlocked":    unlocked,
		"WelcomeText": r.Alleycat.WelcomeText,
	})
}

func (r routes) TextCheckpoint(c *fiber.Ctx) error {
	// isolate link used
	spl := strings.Split(c.OriginalURL(), "/")
	link := spl[len(spl)-1]

	cp, _, err := r.DB.Checkpoints.GetByLink(link)
	if util.CheckWLogs(err) {
		return c.SendStatus(500)
	}

	return c.Render("checkpoint", fiber.Map{
		"Title": cp.Location,
		"CSS":   "checkpoint",
		"Cp":    cp,
	})
}

func (r routes) HelpList(c *fiber.Ctx) error {
	chs, err := r.DB.Checkpoints.GetAll()
	if util.CheckWLogs(err) {
		return c.SendStatus(500)
	}

	return c.Render("helpList", fiber.Map{
		"Title":       "Help",
		"Checkpoints": chs,
	})
}

func (r routes) Help(c *fiber.Ctx) error {
	link := c.Params("link", "")

	return c.Render("help", fiber.Map{
		"Title": "Help",
		"Link":  link,
	})
}

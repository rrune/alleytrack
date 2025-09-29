package router

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rrune/alleytrack/internal/util"
	"github.com/rrune/alleytrack/pkg/models"
)

func (r routes) Index(c *fiber.Ctx) error {
	jwt := c.Cookies("JWT", "")
	if jwt != "" {
		c.Redirect("/manifest")
	}

	return c.Render("index", fiber.Map{
		"Title":       r.Config.LogoText,
		"CSS":         "index",
		"LogoText":    r.Config.LogoText,
		"WelcomeText": r.Config.WelcomeText,
	})
}

func (r routes) SignUp(c *fiber.Ctx) error {
	return c.Render("signup", fiber.Map{
		"Title":    "Sign Up",
		"CSS":      "signup",
		"LogoText": r.Config.LogoText,
	})
}

func (r routes) Login(c *fiber.Ctx) error {
	return c.Render("login", fiber.Map{
		"Title":    "Login",
		"CSS":      "signup",
		"LogoText": r.Config.LogoText,
		"Msg":      c.Query("msg", ""),
	})
}

func (r routes) Manifest(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	number := claims["number"].(string)

	num, err := strconv.Atoi(number)
	if util.CheckWLogs(err) {
		c.SendStatus(500)
	}

	p, _, err := r.DB.GetParicipantFromNumber(num)
	if util.CheckWLogs(err) {
		c.SendStatus(500)
	}
	completed := []models.Checkpoint{}
	unlocked := []models.Checkpoint{}
	// loop through every checkpoint
	for _, c := range r.Config.Manifest {
		// if the id is in the participants list, add the checkpoint to completed
		_, ok := p.Checkpoints[c.ID]
		if ok {

			// if text input checkpoint, inject answer
			if c.Text {
				c.Content = p.Checkpoints[c.ID].Content
			}

			completed = append(completed, c)

			// if its not, check if the requirements are met
		} else {
			// loop through every requirement and check if its met
			met := true
			for _, id := range c.Requirements {
				_, ok := p.Checkpoints[id]
				if !ok {
					met = false
				}
			}
			if met {
				unlocked = append(unlocked, c)
			}
		}
	}

	return c.Render("manifest", fiber.Map{
		"Title":     "Manifest",
		"CSS":       "manifest",
		"LogoText":  "Manifest",
		"Name":      p.Name,
		"Number":    number,
		"Completed": completed,
		"Unlocked":  unlocked,
	})
}

func (r routes) TextCheckpoint(c *fiber.Ctx) error {
	// isolate link used
	spl := strings.Split(c.OriginalURL(), "/")
	link := spl[len(spl)-1]

	var cp models.Checkpoint
	for _, ch := range r.Config.Manifest {
		if ch.Link == link {
			cp = ch
		}
	}

	return c.Render("checkpoint", fiber.Map{
		"Title": cp.Location,
		"CSS":   "checkpoint",
		"Cp":    cp,
	})
}

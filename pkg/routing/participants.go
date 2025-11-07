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

	p, exist, err := r.DB.GetParicipantFromNumber(num)
	if !exist {
		return c.Redirect("/logout")
	}
	if util.CheckWLogs(err) {
		return c.SendStatus(500)
	}
	completed := []models.Checkpoint{}
	unlocked := []models.Checkpoint{}

	if r.Alleycat.Config.Enabled {
		// loop through every checkpoint
		for _, c := range r.Alleycat.Manifest {
			// if the id is in the participants list, add the checkpoint to completed
			_, ok := p.Checkpoints[c.ID]
			if ok {

				// if text input checkpoint, inject answer
				if c.Text {
					c.Content = p.Checkpoints[c.ID].Content
				}

				// inject time
				c.Time = p.Checkpoints[c.ID].Time

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

	var cp models.Checkpoint
	for _, ch := range r.Alleycat.Manifest {
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

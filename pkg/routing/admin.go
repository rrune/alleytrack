package router

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rrune/alleytrack/internal/util"
)

func (r routes) AdminLogin(c *fiber.Ctx) error {
	return c.Render("adminlogin", fiber.Map{
		"Title":    "Login",
		"CSS":      "signup",
		"LogoText": r.Config.LogoText,
		"Msg":      c.Query("msg", ""),
	})
}

func (r routes) Admin(c *fiber.Ctx) error {
	ps, err := r.DB.GetAllParticipants()
	if util.CheckWLogs(err) {
		return c.SendStatus(500)
	}

	return c.Render("admin", fiber.Map{
		"Title":    "Admin Panel",
		"CSS":      "admin",
		"LogoText": "Admin Panel",
		"Ps":       ps,
	})
}

func (r routes) Participant(c *fiber.Ctx) error {
	number := c.Params("number", "")

	num, err := strconv.Atoi(number)
	if util.CheckWLogs(err) {
		return c.Render("response", fiber.Map{
			"Title": "NaN",
			"Text":  "Only valid numbers",
		})
	}

	p, exists, err := r.DB.GetParicipantFromNumber(num)
	if util.CheckWLogs(err) {
		return c.SendStatus(500)
	}
	if !exists {
		return c.Render("response", fiber.Map{
			"Title": "Unknown Number",
			"Text":  "That Number is not tied to a participant",
		})
	}

	for k, v := range p.Checkpoints {
		v.ID = k
		p.Checkpoints[k] = v
	}

	return c.Render("participant", fiber.Map{
		"Title":       number,
		"CSS":         "participant",
		"P":           p,
		"Checkpoints": p.Checkpoints,
	})
}

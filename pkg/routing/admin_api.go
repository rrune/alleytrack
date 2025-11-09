package router

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rrune/alleytrack/internal/util"
	"github.com/rrune/alleytrack/pkg/config"
)

func (r routes) HandleAdminLogin(c *fiber.Ctx) error {
	CallbackPath := c.Query("path", "")
	if CallbackPath == "" {
		CallbackPath = "/"
	}

	// get form values
	password := c.FormValue("password")

	// if number and name match
	if password == r.Alleycat.Config.AdminPassword {
		err := setJwtCookie(c, r.Alleycat.Config.JwtKey, "", "", true)
		if util.CheckWLogs(err) {
			return c.SendStatus(500)
		}

		return c.Redirect(CallbackPath)
	}

	// if they dont
	return c.Redirect("/adminlogin?msg=Incorrect password&path=" + CallbackPath)
}

func (r routes) HandleParticipant(c *fiber.Ctx) error {
	number := c.Params("number", "")
	if number == "" {
		return c.Render("response", fiber.Map{
			"Title": "Missing Data",
			"Text":  "Missing Data",
		})
	}

	num, err := strconv.Atoi(number)
	if util.CheckWLogs(err) {
		return c.SendStatus(500)
	}

	name := c.FormValue("name")
	outoftown := c.FormValue("outoftown") == "on"
	flinta := c.FormValue("flinta") == "on"

	p, exists, err := r.DB.Participants.GetByNumber(num)
	if !exists {
		return c.Render("response", fiber.Map{
			"Title": "Number does not exist",
			"Text":  "Number does not exist",
		})
	}
	if util.CheckWLogs(err) {
		return c.SendStatus(500)
	}

	p.Name = name
	p.OutOfTown = outoftown
	p.Flinta = flinta

	_, err = r.DB.Participants.UpdateByNumber(p.Number, p)
	if util.CheckWLogs(err) {
		return c.SendStatus(500)
	}

	return c.Redirect("/admin")
}

func (r routes) HandleRemoveCheckpoint(c *fiber.Ctx) error {
	number := c.Params("number", "")
	checkpoint := c.Params("checkpoint", "")

	if number == "" || checkpoint == "" {
		return c.Render("response", fiber.Map{
			"Title": "Missing Data",
			"Text":  "Missing Data",
		})
	}

	num, err := strconv.Atoi(number)
	if util.CheckWLogs(err) {
		return c.SendStatus(500)
	}
	ch, err := strconv.Atoi(checkpoint)
	if util.CheckWLogs(err) {
		return c.SendStatus(500)
	}

	_, err = r.DB.ParticipantsCheckpoints.Remove(num, ch)
	if util.CheckWLogs(err) {
		return c.SendStatus(500)
	}

	return c.Redirect("/admin/participant/" + number)
}

func (r routes) HandleRemoveParticipant(c *fiber.Ctx) error {
	number := c.Params("number", "")
	if number == "" {
		return c.Render("response", fiber.Map{
			"Title": "Missing Data",
			"Text":  "Missing Data",
		})
	}
	num, err := strconv.Atoi(number)
	if util.CheckWLogs(err) {
		return c.SendStatus(500)
	}

	exists, err := r.DB.Participants.RemoveByNumber(num)
	if util.CheckWLogs(err) {
		return c.SendStatus(500)
	}
	if !exists {
		return c.Render("response", fiber.Map{
			"Title": "Unsuccessful",
			"Text":  "Unsuccessful",
		})
	}

	return c.Redirect("/admin")
}

func (r routes) HandleSwitchEnabled(c *fiber.Ctx) error {
	err := config.SwitchEnabledInConfig(r.Alleycat)
	if util.CheckWLogs(err) {
		return c.SendStatus(500)
	}

	return c.Redirect("/admin")
}

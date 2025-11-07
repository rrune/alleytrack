package router

import (
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rrune/alleytrack/internal/util"
	"github.com/rrune/alleytrack/pkg/models"
)

func (r routes) HandleSignUp(c *fiber.Ctx) error {
	// get form values
	name := c.FormValue("name")
	number := c.FormValue("number")
	outoftown := c.FormValue("outoftown") == "on"
	flinta := c.FormValue("flinta") == "on"

	// check if number is already taken
	taken, err := r.DB.IsNumberTaken(number)
	if util.CheckWLogs(err) {
		return c.SendStatus(500)
	}
	if taken {
		return c.Render("response", fiber.Map{
			"Title": "Taken",
			"Text":  "Number is already taken. Sorry",
		})
	}

	// convert number to int
	num, err := strconv.Atoi(number)
	if util.CheckWLogs(err) {
		return c.SendStatus(500)
	}

	// add participant to database
	r.DB.NewParticipant(models.Participant{
		Name:      name,
		Number:    num,
		OutOfTown: outoftown,
		Flinta:    flinta,
	})

	err = setJwtCookie(c, r.Alleycat.Config.JwtKey, name, number, false)
	if util.CheckWLogs(err) {
		return c.SendStatus(500)
	}

	return c.Redirect("/")
}

func (r routes) HandleLogin(c *fiber.Ctx) error {
	CallbackPath := c.Query("path", "")
	if CallbackPath == "" {
		CallbackPath = "/"
	}

	// get form values
	name := c.FormValue("name")
	number := c.FormValue("number")

	// number to int
	num, err := strconv.Atoi(number)
	if util.CheckWLogs(err) {
		return c.SendStatus(500)
	}

	// check if number exists
	p, exists, err := r.DB.GetParicipantFromNumber(num)
	if !exists {
		return c.Redirect("/login?msg=Number doesnt exist&path=" + CallbackPath)
	}
	if util.CheckWLogs(err) {
		return c.SendStatus(500)
	}

	// if number and name match
	if strings.EqualFold(p.Name, name) {
		err = setJwtCookie(c, r.Alleycat.Config.JwtKey, name, number, false)
		if util.CheckWLogs(err) {
			return c.SendStatus(500)
		}

		return c.Redirect(CallbackPath)
	}

	// if they dont
	return c.Redirect("/login?msg=That name does not belong to that number&path=" + CallbackPath)
}

func (r routes) HandleLogout(c *fiber.Ctx) error {
	c.ClearCookie("JWT")
	return c.Redirect("/")
}

func (r routes) IsNumerTaken(c *fiber.Ctx) error {
	number := c.Params("number", "")
	if number == "" {
		return c.SendStatus(400)
	}

	taken, err := r.DB.IsNumberTaken(number)
	if util.CheckWLogs(err) {
		return c.SendStatus(500)
	}
	if taken {
		return c.SendString("true")
	}
	return c.SendString("false")
}

func (r routes) CompleteSimpleCheckpoint(c *fiber.Ctx) error {
	return r.completeCheckpoint(c, "")
}

func (r routes) CompleteTextCheckpoint(c *fiber.Ctx) error {
	content := c.FormValue("answer")
	return r.completeCheckpoint(c, content)
}

func (r routes) completeCheckpoint(c *fiber.Ctx, content string) error {
	// isolate link used
	spl := strings.Split(c.OriginalURL(), "/")
	link := spl[len(spl)-1]

	// get user
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	number := claims["number"].(string)

	num, err := strconv.Atoi(number)
	if util.CheckWLogs(err) {
		return c.SendStatus(500)
	}

	// get participant
	p, _, err := r.DB.GetParicipantFromNumber(num)
	if util.CheckWLogs(err) {
		return c.SendStatus(500)
	}

	if len(p.Checkpoints) == 0 {
		p.Checkpoints = make(map[int]models.ParticipantCheckpoint)
	}

	// add the checkpoint to the completed ones
	for _, ch := range r.Alleycat.Manifest {
		if ch.Link == link {
			p.Checkpoints[ch.ID] = models.ParticipantCheckpoint{Time: time.Now(), Content: content}
			err = r.DB.UpdateCheckpoints(p)
			if util.CheckWLogs(err) {
				return c.SendStatus(500)
			}
			return c.Render("response", fiber.Map{
				"Title": "Checkpoint completed",
				"Text":  "Succesfully completed Checkpoint. Congratulations!",
			})
		}
	}
	return c.Render("response", fiber.Map{
		"Title": "Something went wrong...",
		"Text":  "Hmmm",
	})
}

func setJwtCookie(c *fiber.Ctx, jwtKey string, name string, number string, admin bool) error {
	claims := jwt.MapClaims{
		"name":   name,
		"number": number,
		"admin":  admin,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(jwtKey))
	if util.CheckWLogs(err) {
		return err
	}
	cookie := fiber.Cookie{
		Name:  "JWT",
		Value: t,

		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	return nil
}

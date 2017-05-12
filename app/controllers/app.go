package controllers

import (
	"revel_start/app/routes"

	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
	"revel_start/app/models"
	"github.com/Masterminds/squirrel"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) Hello(myName string) revel.Result {
	c.Validation.Required(myName).Message("Your name is required!")
	c.Validation.MinSize(myName, 3).Message("Your name is not long enough!")

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(App.Index)
	}

	return c.Render(myName)
}

type user struct {
	HashedPassword []byte
}

func (c App) Login(username, password string, remember bool) revel.Result {

	hash := []byte(`qwe`)

	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err == nil {
		c.Session["user"] = username
		if remember {
			c.Session.SetDefaultExpiration()
		} else {
			c.Session.SetNoExpiration()
		}
		c.Flash.Success("Welcome, " + username)
		return c.Redirect(routes.App.Index())
	}

	c.Flash.Out["username"] = username
	c.Flash.Error("Login failed")
	return c.Redirect(routes.App.Index())
}
func (c App) Register() revel.Result {
	return c.Render()
}

func (c App) SaveUser(user models.User, verifyPassword string) revel.Result {
	c.Validation.Required(verifyPassword)
	c.Validation.Required(verifyPassword == user.Password).
		Message("Password does not match")
	user.Validate(c.Validation)

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.App.Register())
	}

	user.HashedPassword, _ = bcrypt.GenerateFromPassword(
		[]byte(user.Password), bcrypt.DefaultCost)
	qb :=squirrel.Insert("squirrel.user").Values("login")

	//err := c.Txn.Insert(&user)
	if err != nil {
		panic(err)
	}

	c.Session["user"] = user.Username
	c.Flash.Success("Welcome, " + user.Name)
	return c.Redirect(routes.Hotels.Index())
}


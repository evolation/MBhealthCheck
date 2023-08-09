//go:build js && wasm

package main

import (
	"log"
	"net/http"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type login struct {
	app.Compo

	Username string
	Password string
	Message  string
}

func (l *login) Render() app.UI {
	return app.Div().Body(
		app.H1().Text("Login"),
		app.Input().
			Type("text").
			Value(l.Username).
			OnChange(l.OnUsernameChange),
		app.Input().
			Type("password").
			Value(l.Password).
			OnChange(l.OnPasswordChange),
		app.Button().
			Text("Login").
			OnClick(l.OnLoginClick),
		app.P().Text(l.Message),
	)
}

func (l *login) OnUsernameChange(ctx app.Context, e app.Event) {
	l.Username = ctx.JSSrc().Get("value").String()
	l.Update()
}

func (l *login) OnPasswordChange(ctx app.Context, e app.Event) {
	l.Password = ctx.JSSrc().Get("value").String()
	l.Update()
}

func (l *login) OnLoginClick(ctx app.Context, e app.Event) {
	// TODO: Handle login logic here. For simplicity, we just check hardcoded values.
	if l.Username == "admin" && l.Password == "password" {
		ctx.Navigate("/table") // Navigate to table page
	} else {
		l.Message = "Invalid credentials."
	}
	l.Update()
}

func main() {
	app.Route("/", &login{})
	app.Route("/table", &tablePage{})
	app.RunWhenOnBrowser()

	// HTTP routing:
	http.Handle("/", &app.Handler{
		Name:        "PaulLogin",
		Description: "Login Page for End Users",
	})

	if err := http.ListenAndServe(":7999", nil); err != nil {
		log.Fatal(err)
	}

}

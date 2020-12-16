package main

import (
	"fmt"
	"html/template"
	"net/http"

	"log"

	"github.com/go-chi/chi"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func main() {
	key := "vlDxjmHJX80vOuHa5THxfCsR" // Replace with your SESSION_SECRET or similar
	maxAge := 86400 * 30              // 30 days
	isProd := false                   // Set to true when serving over https
	clientid := "379756554270-olm9ma6g4dru3lil2cse84eaeimpj0u2.apps.googleusercontent.com"
	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true // HttpOnly should always be enabled
	store.Options.Secure = isProd

	gothic.Store = store
	goth.UseProviders(google.New(
		clientid,
		key,
		"http://localhost:4000/auth/callback?provider=google", "email", "profile"))

	p := chi.NewRouter()

	p.Get("/auth/callback", func(res http.ResponseWriter, req *http.Request) {

		user, err := gothic.CompleteUserAuth(res, req)
		if err != nil {
			fmt.Fprintln(res, err)
			return
		}

		t, _ := template.ParseFiles("templates/success.html")
		t.Execute(res, user)
	})
	p.Get("/auth", gothic.BeginAuthHandler)

	p.Get("/", func(res http.ResponseWriter, req *http.Request) {
		t, _ := template.ParseFiles("templates/index.html")
		t.Execute(res, false)
	})
	log.Println("listening on localhost:4000")
	log.Fatal(http.ListenAndServe(":4000", p))
}

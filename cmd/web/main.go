package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"log"

	"github.com/go-chi/chi"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

// func home(w http.ResponseWriter, r *http.Request) {

// 	files := []string{
// 		"./ui/html/home.page.tmpl",
// 		"./ui/html/base.layout.tmpl",
// 	}

// 	ts, err := template.ParseFiles(files...)
// 	if err != nil {
// 		log.Println(err.Error())
// 		http.Error(w, "Internal Server Error", 500)
// 		return
// 	}
// 	err = ts.Execute(w, nil)
// 	if err != nil {
// 		log.Println(err.Error())
// 		http.Error(w, "Internal Server Error", 500)
// 	}
// 	//w.Write([]byte("HEllo from test"))
// }

func main() {
	key := "DoN4QZCXaa3TJfr4BJZMQZNo" // Replace with your SESSION_SECRET or similar
	maxAge := 86400 * 30              // 30 days
	isProd := false                   // Set to true when serving over https

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true // HttpOnly should always be enabled
	store.Options.Secure = isProd

	gothic.Store = store

	// goth.UseProviders(
	// 	google.New("263741611747-2bgmmh2vnbjvt02c3m8s30ujbb76obgf.apps.googleusercontent.com", "DoN4QZCXaa3TJfr4BJZMQZNo", "http://localhost:4000/auth/google/callback", "email", "profile"),
	// )
	goth.UseProviders(google.New(
		os.Getenv("263741611747-2bgmmh2vnbjvt02c3m8s30ujbb76obgf.apps.googleusercontent.com"),
		os.Getenv("DoN4QZCXaa3TJfr4BJZMQZNo"),
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

	// p.Get("/auth", func(res http.ResponseWriter, req *http.Request) {
	// 	gothic.BeginAuthHandler(res, req)
	// })
	p.Get("/auth", gothic.BeginAuthHandler)

	p.Get("/", func(res http.ResponseWriter, req *http.Request) {
		t, _ := template.ParseFiles("templates/index.html")
		t.Execute(res, false)
	})
	log.Println("listening on localhost:4000")
	log.Fatal(http.ListenAndServe(":4000", p))
}

// mux := http.NewServeMux()
// mux.HandleFunc("/", home)

// log.Println("Starting server on :4000")
// err := http.ListenAndServe(":4000", mux)
// log.Fatal(err)

package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type User struct {
	Email    string
	Password string
}

var users = map[string]User{}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles("templates/" + tmpl + ".html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
}

func landingPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index", nil)
}

func signInPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		user, exists := users[email]
		if exists && user.Password == password {
			fmt.Fprintf(w, "✅ Signed in successfully! Welcome, %s", email)
		} else {
			fmt.Fprintf(w, "❌ Invalid credentials")
		}
		return
	}
	renderTemplate(w, "signin", nil)
}

func signUpPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		if _, exists := users[email]; exists {
			fmt.Fprintf(w, "⚠️ User already exists.")
			return
		}

		users[email] = User{Email: email, Password: password}
		fmt.Fprintf(w, "✅ Signed up successfully! You can now sign in.")
		return
	}
	renderTemplate(w, "signup", nil)
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", landingPage)
	http.HandleFunc("/signin", signInPage)
	http.HandleFunc("/signup", signUpPage)

	fmt.Println("🌐 Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

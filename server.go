package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/alexedwards/scs"
)

// HomePage : home page information; most notably posts, maybe more later
type HomePage struct {
	Posts []Post `json:"posts"`
}

var sessionManager *scs.SessionManager

func loadHTML(file string) (string, error) {
	body, err := ioutil.ReadFile("web/" + file + ".html")
	if err != nil {
		return "", err
	}
	return string(body[:]), nil
}

func allPostsHandler(w http.ResponseWriter, r *http.Request) {
	html, err := loadHTML("index")
	check(err)

	var homePage HomePage
	homePage.Posts = posts

	t, err := template.New("home").Parse(html)
	if err != nil {
		fmt.Fprintf(w, "Error occured")
		fmt.Println(err)
		return
	}

	t.Execute(w, homePage)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	errPage, err := loadHTML("notPost")
	check(err)

	postID, err := strconv.Atoi(r.URL.Path[len("/post/"):])
	if err != nil {
		fmt.Fprintf(w, errPage)
		return
	}

	post := getPost(postID)
	if post == nil {
		fmt.Fprintf(w, errPage)
		return
	}

	page, err := loadHTML("post")
	if err != nil {
		fmt.Fprintf(w, errPage)
		return
	}

	t, err := template.New("post").Parse(page)
	if err != nil {
		fmt.Fprintf(w, errPage)
		return
	}

	t.Execute(w, post)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	admin := sessionManager.GetBool(r.Context(), "admin")
	if admin {
		http.Redirect(w, r, "../", 302)
	}

	page, err := loadHTML("login")
	check(err)

	if r.Method == "GET" {
		fmt.Fprintf(w, page)
	} else if r.Method == "POST" {
		err := r.ParseForm()
		check(err)

		d := r.Form
		fmt.Print(d)
		pass := strings.Join(d["password"], "")
		if pass == getConfig()["password"] {
			fmt.Fprintf(w, "lit")

		} else {
			fmt.Fprintf(w, "not lit")
		}
	}

}

func newPostHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	admin := sessionManager.GetBool(r.Context(), "admin")
	if admin {
		http.Redirect(w, r, "../", 302)
	}

	page, err := loadHTML("newPost")
	check(err)

	if r.Method == "GET" {
		fmt.Fprintf(w, page)
	} else if r.Method == "POST" {
		err := r.ParseForm()
		check(err)

		d := r.Form
		fmt.Print(d)
		if admin {

			pass := strings.Join(d["password"], "")
			if pass == getConfig()["password"] {
				fmt.Fprintf(w, "lit")

			} else {
				fmt.Fprintf(w, "not lit")
			}

			//post(db)
		}
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	page, err := loadHTML("index")
	check(err)

	fmt.Fprintf(w, page)
}

func serve(port int, db *sql.DB) {
	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour

	mux := http.NewServeMux()

	mux.HandleFunc("/", allPostsHandler)
	mux.HandleFunc("/post/", postHandler)
	mux.HandleFunc("/login/", loginHandler)
	mux.HandleFunc("/newPost/", func(w http.ResponseWriter, r *http.Request) {
		newPostHandler(w, r, db)
	})
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	fmt.Printf("Setting up port: %d. https://127.0.0.1:%[1]d", port)
	err := http.ListenAndServeTLS(":443", "cert.pem", "key.pem", sessionManager.LoadAndSave(mux))
	log.Fatal(err)
}

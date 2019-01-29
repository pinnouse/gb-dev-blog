package main

import (
	"strconv"
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
	"html/template"
)

// HomePage : home page information; most notably posts, maybe more later
type HomePage struct {
	Posts []Post `json:"posts"`
}

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

func handler(w http.ResponseWriter, r *http.Request) {
	page, err := loadHTML("index")
	check(err)

	fmt.Fprintf(w, page)
}

func serve(port int) {
	fmt.Printf("Serving on port %d ", port)

	http.HandleFunc("/", allPostsHandler)
	http.HandleFunc("/post/", postHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}
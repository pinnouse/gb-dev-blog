package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"time"

	_ "github.com/lib/pq"
)

var posts []Post

// Post : Post information
type Post struct {
	ID       int           `json:"id"`
	Title    string        `json:"title"`
	ImageURL string        `json:"imageURL"`
	Date     time.Time     `json:"date"`
	Content  template.HTML `json:"content"`
}

func checkTable(db *sql.DB) bool {
	fmt.Println("Checking if table already exists")
	_, err := db.Query("SELECT 1 FROM posts LIMIT 1")
	if err == nil {
		fmt.Println("Table exists")
		return true
	}

	_, err = db.Exec(`CREATE TABLE posts ( 
		id SERIAL PRIMARY KEY,
		title varchar(255),
		imageURL varchar(255),
		date timestamp,
		content text)`)
	if err != nil {
		panic(err)
	}

	fmt.Println("Created new table")
	return true
}

func connect(username string, password string, dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s/postsdb?sslmode=disable", username, password, dbURL))
	if err != nil {
		panic(err)
	}

	return db, err
}

func post(db *sql.DB, title string, imageURL string, content string) bool {
	stmt, err := db.Prepare("INSERT INTO posts (title, imageURL, date, content) values (?,?,?,?)")
	if err != nil {
		return false
	}

	t := time.Now()
	_, err = stmt.Exec(title, imageURL, t, content)
	if err != nil {
		return false
	}

	return true
}

func getPosts(db *sql.DB) *[]Post {
	rows, err := db.Query("SELECT * FROM posts")
	check(err)

	posts = nil
	for rows.Next() {
		var p Post
		err = rows.Scan(&p.ID, &p.Title, &p.ImageURL, &p.Date, &p.Content)
		check(err)
		posts = append(posts, p)
	}

	rows.Close()

	return &posts
}

func getPost(id int) *Post {
	for _, p := range posts {
		if p.ID == id {
			return &p
		}
	}
	return nil
}

package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	setDB()

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", index)
	mux.HandleFunc("/post/id", post)
	mux.HandleFunc("/commentedPosts", commentedPosts)
	mux.HandleFunc("/dashBoard", dashBoard)
	mux.HandleFunc("/myPosts", myPosts)
	mux.HandleFunc("/newPost", newPost)
	mux.HandleFunc("/likedPosts", likedPosts)
	mux.HandleFunc("/dislikedPosts", dislikedPosts)
	// mux.HandleFunc("/editComment", editComment)
	// mux.HandleFunc("/editPost", editPost)
	mux.HandleFunc("/error", showError)

	// Handle forms
	mux.HandleFunc("/auth", auth)
	mux.HandleFunc("/registration", registration)
	mux.HandleFunc("/logout", logout)
	mux.HandleFunc("/post/like/id", likePost)
	mux.HandleFunc("/post/dislike/id", dislikePost)
	mux.HandleFunc("/comment/like/", likeComment)
	mux.HandleFunc("/comment/dislike/", dislikeComment)
	mux.HandleFunc("/comment", createComment)
	mux.HandleFunc("/createPost", createPost)

	// Create a custom server with a timeout
	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
	}

	fmt.Println("\nStarting server at http://127.0.0.1:8080/")
	fmt.Println("Quit the server with CONTROL-C.")

	// Start the server
	log.Fatal(server.ListenAndServe())
}

var database *sql.DB

func setDB() {

	file, err := os.Create("database.db")
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
	database, _ = sql.Open("sqlite3", "database.db")
	createUsersTable(database)
	createThreadsTable(database)
	createPostsTable(database)
	createCommentsTable(database)
	createCommentsReactionsTable(database)
	createPostsReactionsTable(database)

	addUser(database, "test", "test@gmail.com", "1234")
	addUser(database, "Lasso-less Cowboy", "cowboy@gmail.com", "blacksheep")
	addUser(database, "SnapHappy", "snaphappyphotographer@email.com", "ZoomZoomZap")
	addUser(database, "RodeoQueen", "rodeoqueen@email.com", "kingtomyqueen")

	addThread(database, "Ranch", 1)
	addThread(database, "Dogs", 1)
	addThread(database, "Other", 1)

	addPost(database, title1, image1, post1, threads1, 2, 2, 1)
	addPost(database, title2, image2, post2, threads2, 3, 3, 2)
	addPost(database, title3, image3, post3, threads3, 2, 2, 0)
	addPost(database, title4, image4, post4, threads4, 4, 7, 1)

	addComment(database, comment1_1, 1, 2, 1, 0)
	addComment(database, comment1_2, 1, 3, 2, 0)
	addComment(database, comment1_3, 1, 4, 0, 0)
	addComment(database, comment2_1, 2, 2, 2, 1)
	addComment(database, comment2_2, 2, 3, 0, 0)
	addComment(database, comment3_1, 3, 4, 3, 0)
	addComment(database, comment4_1, 4, 2, 0, 1)
	addComment(database, comment4_2, 4, 3, 2, 2)

}

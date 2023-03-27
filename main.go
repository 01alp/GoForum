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
	mux.HandleFunc("/editComment", editComment)
	mux.HandleFunc("/editPost", editPost)
	mux.HandleFunc("/error", showError)

	// Handle forms
	mux.HandleFunc("/auth", auth)
	mux.HandleFunc("/registration", registration)
	mux.HandleFunc("/logout", logout)
	mux.HandleFunc("/post/like/id", like)
	mux.HandleFunc("/post/dislike/id", dislike)

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
	addUser(database, "test2", "test2@gmail.com", "1234")
	addUser(database, "test123", "test123@gmail.com", "varwa123")

	addThread(database, "Cat", 1)
	addThread(database, "Kitten", 1)
	addThread(database, "Kittens", 1)

	var threads1 = []string{"Cat"}
	addPost(database, "Some smart thoughts", "blablabla", threads1, 1)
	var threads2 = []string{"Kitten"}
	addPost(database, "Some smarter thoughts 1", "blablabla", threads2, 2)
	var threads3 = []string{"Cat", "Kitten"}
	addPost(database, "Some smartest thoughts 2", "blablabla", threads3, 1)
	var threads4 = []string{"Kittens"}
	addPost(database, "Some smartest thoughts 3", "blablablaaa", threads4, 1)
	var threads5 = []string{"Kittens"}
	addPost(database, "Some smartest thoughts 4", "blablablaaa", threads5, 1)
	var threads6 = []string{"Kittens"}
	addPost(database, "Some smartest thoughts 5", "blablablaaa", threads6, 1)
	var threads7 = []string{"Kittens"}
	addPost(database, "Some smartest thoughts 6", "blablablaaa", threads7, 1)

	addComment(database, "Hello brave adventurer", 1, 2)
	addComment(database, "Hello GOOD adventurer", 1, 2)
	addComment(database, "Hello GOOD adventurer", 1, 1)
	addComment(database, "Hello GOOD adventurer", 2, 2)

	// fetchUsers(database)
	// fmt.Println("-----------------------")
	// fetchUserByEmail(database, "test@gmail.com")

	// fmt.Println("-----------------------")
	// fetchPostsByThread(database, "Kitten")
	// fetchPostsByThread(database, "Cat")

	// fmt.Println("-----------------------")
	// fetchPostsByUser(database, 1)

	// fmt.Println("-----------------------")
	// fetchCommentsByPost(database, 1)

	// fmt.Println("-----------------------")
	// fetchCommentsByUser(database, 1)

	// reaction test
	// addPostsReactions(database, -1, 1, 1)
	// addPostsReactions(database, 0, 2, 2)
	fmt.Println()
	// fetchReactionByUserAndPost(database, 1, 1)
	// fetchReactionByUserAndPost(database, 1, 2)
	// fetchReactionByUserAndPost(database, 2, 2)

	// addCommentsReactions(database, -1, 1, 1)
	// addCommentsReactions(database, 0, 2, 2)
	// fmt.Println()
	// fetchReactionByUserAndComment(database, 1, 1)
	// fetchReactionByUserAndComment(database, 1, 2)
	// fetchReactionByUserAndComment(database, 2, 2)
}

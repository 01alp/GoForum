package main

import (
	"fmt"
	"log"
	"strings"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Id        int
	Email     string
	Username  string
	Password  string
	Timestamp string
}

type Post struct {
	Id        int
	Title     string
	Content   string
	Thread    string
	UserId    int
	Likes     int
	Dislikes  int
	Timestamp string
	Comments  []Comment
	User      User

	UserReaction int
}

type Comment struct {
	Id        int
	Content   string
	PostId    int
	UserId    int
	Likes     int
	Dislikes  int
	Timestamp string
	User      User
}

type Reaction struct {
	Id     int
	Value  int
	UserId int
	UnitId int
}

// CRUD: 4 main db commands usually --> create, read, update, delete (we will need just create and read probably)

// users
// -------------------------------------------------------------------------------------

func createUsersTable(db *sql.DB) {
	users_table := `CREATE TABLE users (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        "Username" TEXT UNIQUE,
        "Email" TEXT UNIQUE,
        "Password" TEXT,
        timestamp TEXT DEFAULT CURRENT_TIMESTAMP);`
	query, err := db.Prepare(users_table)
	if err != nil {
		log.Fatal(err)
	}
	query.Exec()
	fmt.Println("Table for users created successfully!")
}

func addUser(db *sql.DB, Username string, Email string, Password string) {
	records := `INSERT INTO users(Username, Email, Password) VALUES (?, ?, ?)`
	query, err := db.Prepare(records)
	if err != nil {
		log.Print(err)
	}
	_, err = query.Exec(Username, Email, Password)
	if err != nil {
		log.Print(err)
	}
}

func fetchUserByEmail(db *sql.DB, email string) User {
	var user User
	db.QueryRow("SELECT * FROM users WHERE email=?", email).Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.Timestamp)
	return user
}

func fetchUserByUsername(db *sql.DB, username string) User {
	var user User
	db.QueryRow("SELECT * FROM users WHERE username=?", username).Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.Timestamp)
	return user
}

func fetchUserById(db *sql.DB, id int) User {
	var user User
	db.QueryRow("SELECT * FROM users WHERE id=?", id).Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.Timestamp)
	return user
}

// func fetchUsers(db *sql.DB) {
//     record, err := db.Query("SELECT * FROM users")
//     if err != nil {
//         log.Fatal(err)
//     }
//     defer record.Close()
// for record.Next() {
//         var id int
//         var Username string
//         var Email string
//         var Password string
//         record.Scan(&id, &Username, &Email, &Password)
//         fmt.Printf("User: %d %s %s %s \n", id, Username, Email, Password)
//     }
// }

// threads (categories)
// -------------------------------------------------------------------------------------

func createThreadsTable(db *sql.DB) {
	threads_table := `CREATE TABLE threads (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        "Subject" TEXT UNIQUE,
        "User_id" INTEGER,
        timestamp DATETIME DEFAULT CURRENT_TIMESTAMP);`
	query, err := db.Prepare(threads_table)
	if err != nil {
		log.Fatal(err)
	}
	query.Exec()
	fmt.Println("Table for threads created successfully!")
}

func addThread(db *sql.DB, Subject string, User_id int) {
	records := `INSERT INTO threads(Subject, User_id) VALUES (?, ?)`
	query, err := db.Prepare(records)
	if err != nil {
		log.Fatal(err)
	}
	_, err = query.Exec(Subject, User_id)
	if err != nil {
		log.Fatal(err)
	}
}

func fetchAllThreads(db *sql.DB) []string {
	record, err := db.Query("SELECT Subject FROM threads")
	if err != nil {
		log.Fatal(err)
	}
	defer record.Close()

	var threads []string
	for record.Next() {
		var thread string
		err = record.Scan(&thread)
		if err != nil {
			log.Println(err)
		}
		threads = append(threads, thread)
	}
	fmt.Println(threads)
	return threads

}

// posts
// -------------------------------------------------------------------------------------

func createPostsTable(db *sql.DB) {
	posts_table := `CREATE TABLE posts (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"Title" TEXT,
        "Content" TEXT,
        "Subject" TEXT,
        "User_id" INTEGER,
		"Likes"	INTEGER DEFAULT 0,
		"Dislikes" INTEGER DEFAULT 0,
        timestamp TEXT DEFAULT(strftime('%Y.%m.%d %H:%M', 'now')));`
	query, err := db.Prepare(posts_table)
	if err != nil {
		log.Fatal(err)
	}
	query.Exec()
	fmt.Println("Table for posts created successfully!")
}

func addPost(db *sql.DB, Title string, Content string, Subject []string, User_id int) {
	records := `INSERT INTO posts(Title, Content, Subject, User_id) VALUES (?, ?, ?, ?)`
	query, err := db.Prepare(records)
	if err != nil {
		log.Fatal(err)
	}
	_, err = query.Exec(Title, Content, strings.Join(Subject, ", "), User_id)
	if err != nil {
		log.Fatal(err)
	}
}

// func updatePost(db *sql.DB, id int, Title string, Content string, Subject string, Likes int, Dislikes int) {
// 	db.Exec("UPDATE posts SET title = ?, content = ?, subject = ?, likes = ?, dislikes = ? WHERE id = ?", Title, Content, Subject, Likes, Dislikes, id)
// }

func updatePostLikes(db *sql.DB, Likes int, id int) {
	db.Exec("UPDATE posts SET likes = ? WHERE id = ?", Likes, id)
}

func updatePostDislikes(db *sql.DB, Dislikes int, id int) {
	db.Exec("UPDATE posts SET dislikes = ? WHERE id = ?", Dislikes, id)
}

func fetchPostsByThread(db *sql.DB, subject string) []Post {
	record, err := db.Query("SELECT * FROM posts WHERE Subject=? OR Subject LIKE  ? OR Subject LIKE ? OR Subject LIKE ?", subject, subject+",%", "%, "+subject+",%", "%, "+subject)
	if err != nil {
		log.Fatal(err)
	}
	defer record.Close()

	var posts []Post
	for record.Next() {
		var post Post
		err = record.Scan(&post.Id, &post.Title, &post.Content, &post.Thread, &post.UserId, &post.Likes, &post.Dislikes, &post.Timestamp)
		if err != nil {
			log.Println(err)
		}
		posts = append(posts, post)
	}

	for _, p := range posts {
		fmt.Printf("Post by thread %s: Original thread: %s; Content: %s; User_id: %d; Likes: %d; Dislikes: %d; Timestamp: %s \n", subject, p.Thread, p.Content, p.UserId, p.Likes, p.Dislikes, p.Timestamp)
	}
	return posts
}

func fetchPostsByUser(db *sql.DB, user_id int) []Post {
	record, err := db.Query("SELECT * FROM posts WHERE user_id=?", user_id)
	if err != nil {
		log.Fatal(err)
	}
	defer record.Close()

	var posts []Post
	for record.Next() {
		var post Post
		err = record.Scan(&post.Id, &post.Title, &post.Content, &post.Thread, &post.UserId, &post.Likes, &post.Dislikes, &post.Timestamp)
		if err != nil {
			log.Println(err)
		}
		posts = append(posts, post)
	}
	for _, p := range posts {
		fmt.Printf("Posts by User %d: id: %d; Content: %s; Subject: %s; User_id: %d, Likes: %d; Dislikes: %d; Timestamp: %s \n", user_id, p.Id, p.Thread, p.Content, p.UserId, p.Likes, p.Dislikes, p.Timestamp)
	}

	return posts
}

func fetchAllPosts(db *sql.DB) []Post {
	record, err := db.Query("SELECT * FROM posts")
	if err != nil {
		log.Fatal(err)
	}
	defer record.Close()

	var posts []Post
	for record.Next() {
		var post Post
		err = record.Scan(&post.Id, &post.Title, &post.Content, &post.Thread, &post.UserId, &post.Likes, &post.Dislikes, &post.Timestamp)
		if err != nil {
			log.Println(err)
		}
		posts = append(posts, post)
	}
	for _, p := range posts {
		fmt.Printf("Posts by User %d: id: %d; Content: %s; Subject: %s, Likes: %d; Dislikes: %d; Timestamp: %s \n", p.UserId, p.Id, p.Content, p.Thread, p.Likes, p.Dislikes, p.Timestamp)
	}
	return posts
}

func fetchPostByID(db *sql.DB, id int) Post {
	record, err := db.Query("SELECT * FROM posts WHERE id=?", id)
	if err != nil {
		log.Fatal(err)
	}
	defer record.Close()

	var post Post
	for record.Next() {
		err = record.Scan(&post.Id, &post.Title, &post.Content, &post.Thread, &post.UserId, &post.Likes, &post.Dislikes, &post.Timestamp)
		if err != nil {
			log.Println(err)
		}
	}
	return post
}

// comments
// -------------------------------------------------------------------------------------

func createCommentsTable(db *sql.DB) {
	posts_table := `CREATE TABLE comments (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        "Content" TEXT,
        "Post_id" INTEGER,
        "User_id" INTEGER,
		"Likes"	INTEGER DEFAULT 0,
		"Dislikes" INTEGER DEFAULT 0,
        timestamp TEXT DEFAULT(strftime('%Y.%m.%d %H:%M', 'now')));`
	query, err := db.Prepare(posts_table)
	if err != nil {
		log.Fatal(err)
	}
	query.Exec()
	fmt.Println("Table for comments created successfully!")
}

func addComment(db *sql.DB, Content string, Post_id int, User_id int) {
	records := `INSERT INTO comments(Content, Post_id, User_id) VALUES (?, ?, ?)`
	query, err := db.Prepare(records)
	if err != nil {
		log.Fatal(err)
	}
	_, err = query.Exec(Content, Post_id, User_id)
	if err != nil {
		log.Fatal(err)
	}
}

func fetchCommentsByPost(db *sql.DB, post_id int) []Comment {
	record, err := db.Query("SELECT * FROM comments WHERE post_id=?", post_id)
	if err != nil {
		log.Fatal(err)
	}
	defer record.Close()

	var comments []Comment
	for record.Next() {
		var comment Comment
		err = record.Scan(&comment.Id, &comment.Content, &comment.PostId, &comment.UserId, &comment.Likes, &comment.Dislikes, &comment.Timestamp)
		if err != nil {
			log.Println(err)
		}
		comments = append(comments, comment)

	}

	for _, c := range comments {
		fmt.Printf("Comments of post %d: %v \n", c.PostId, c.Content)
	}

	// // try to append comments to the specif post, not working
	// allPosts := Posts{}
	// for _, p := range allPosts.posts {
	// 	if p.Id == post_id {
	// 		p.Comments = append(p.Comments, comments...)
	// 		fmt.Println("I appended comments", comments, "to post", p.Id)
	// 		break
	// 	}
	// }
	return comments
}

// for advanced features project
func fetchCommentsByUser(db *sql.DB, user_id int) []int {
	record, err := db.Query("SELECT Post_id FROM comments WHERE user_id=?", user_id)
	if err != nil {
		log.Fatal(err)
	}
	defer record.Close()

	var postIDs []int
	for record.Next() {
		var id int
		err = record.Scan(&id)
		if err != nil {
			log.Println(err)
		}

		if isUnique(id, postIDs) {
			postIDs = append(postIDs, id)
		}
	}
	return postIDs
}

func isUnique(id int, ids []int) bool {
	for _, i := range ids {
		if i == id {
			return false
		}
	}
	return true
}

//	reaction tables
//
// -------------------------------------------------------------------------------------
func createCommentsReactionsTable(db *sql.DB) {
	posts_table := `CREATE TABLE commentsReactions (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        "Reaction" INTEGER,
        "User_id" INTEGER,
        "Comment_id" INTEGER)`
	query, err := db.Prepare(posts_table)
	if err != nil {
		log.Fatal(err)
	}
	query.Exec()
	fmt.Println("Table for comments reactions created successfully!")
}

func createPostsReactionsTable(db *sql.DB) {
	posts_table := `CREATE TABLE postsReactions (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        "Reaction" INTEGER,
        "User_id" INTEGER,
        "Post_id" INTEGER)`
	query, err := db.Prepare(posts_table)
	if err != nil {
		log.Fatal(err)
	}
	query.Exec()
	fmt.Println("Table for posts reactions created successfully!")
}

func addCommentsReactions(db *sql.DB, Reaction int, User_id int, Comment_id int) {
	records := `INSERT INTO commentsReactions(Reaction, User_id, Comment_id) VALUES (?, ?, ?)`
	query, err := db.Prepare(records)
	if err != nil {
		log.Fatal(err)
	}
	_, err = query.Exec(Reaction, User_id, Comment_id)
	if err != nil {
		log.Fatal(err)
	}
}

func addPostsReactions(db *sql.DB, Reaction int, User_id int, Post_id int) {
	records := `INSERT INTO postsReactions(Reaction, User_id, Post_id) VALUES (?, ?, ?)`
	query, err := db.Prepare(records)
	if err != nil {
		log.Fatal(err)
	}
	_, err = query.Exec(Reaction, User_id, Post_id)
	if err != nil {
		log.Fatal(err)
	}
}

func fetchReactionByUserAndId(db *sql.DB, table string, user_id int, post_id int) Reaction {
	var reaction Reaction

	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=? AND post_id=?", table)
	db.QueryRow(query, user_id, post_id).Scan(&reaction.Id, &reaction.Value, &reaction.UserId, &reaction.UnitId)

	if reaction.Id == 0 {
		fmt.Println("No reaction found")
		return Reaction{}
	}

	return reaction
}

func deleteRow(db *sql.DB, table string, id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", table)
	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

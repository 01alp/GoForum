package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func likePost(w http.ResponseWriter, r *http.Request) {
	post_id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	posts := fetchAllPosts(database)
	user := welcome(w, r).User
	allLikes := posts[post_id-1].Likes
	allDislikes := posts[post_id-1].Dislikes

	reactionsPosts := "postsReactions"
	postsTable := "posts"

	if user.Id == 0 { // if not login
		http.Redirect(w, r, "/?modal=true", http.StatusSeeOther)
	} else {
		if post_id > 0 && post_id <= len(posts) {
			fetch := fetchReactionByUserAndId(database, reactionsPosts, user.Id, post_id)
			if fetch.Value != 1 { // if not like
				if fetch.Value == -1 {
					// delete dislike in frontend and backend
					deleteRow(database, reactionsPosts, fetch.Id)
					updateTableDislikes(database, postsTable, allDislikes-1, post_id)
				}
				// add like in frontend and backend
				updateTableLikes(database, postsTable, allLikes+1, post_id)
				addPostsReactions(database, 1, user.Id, post_id)
			} else {
				// delete like in frontend and backend
				updateTableLikes(database, postsTable, allLikes-1, post_id)
				deleteRow(database, reactionsPosts, fetch.Id)
			}
		}
		c, err := r.Cookie("last_page")
		if err != nil {
			fmt.Println("cookie err", err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		fmt.Println("REDIRECT TO", c.Value)
		http.Redirect(w, r, c.Value, http.StatusSeeOther)
	}
}

func dislikePost(w http.ResponseWriter, r *http.Request) {
	post_id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	posts := fetchAllPosts(database)
	user := welcome(w, r).User
	allLikes := posts[post_id-1].Likes
	allDislikes := posts[post_id-1].Dislikes

	reactionsPosts := "postsReactions"
	postsTable := "posts"

	if user.Id == 0 { // if not login
		http.Redirect(w, r, "/?modal=true", http.StatusSeeOther)
	} else {
		if post_id > 0 && post_id <= len(posts) {
			fetch := fetchReactionByUserAndId(database, reactionsPosts, user.Id, post_id)
			if fetch.Value != -1 { // if not dislike
				if fetch.Value == 1 {
					// delete like in frontend and backend
					deleteRow(database, reactionsPosts, fetch.Id)
					updateTableLikes(database, postsTable, allLikes-1, post_id)
				}
				// add dislike in frontend and backend
				updateTableDislikes(database, postsTable, allDislikes+1, post_id)
				addPostsReactions(database, -1, user.Id, post_id)
			} else {
				// delete dislike in frontend and backend
				updateTableDislikes(database, postsTable, allDislikes-1, post_id)
				deleteRow(database, reactionsPosts, fetch.Id)
			}
		}
		c, err := r.Cookie("last_page")
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		http.Redirect(w, r, c.Value, http.StatusSeeOther)
	}
}

func likeComment(w http.ResponseWriter, r *http.Request) {
	comment_id, _ := strconv.Atoi(r.URL.Query().Get("comment_id"))
	user := welcome(w, r).User
	allLikes := fetchCommentByID(database, comment_id).Likes
	allDislikes := fetchCommentByID(database, comment_id).Dislikes

	reactionsComments := "commentsReactions"
	commentsTable := "comments"

	if user.Id == 0 { // if not login
		http.Redirect(w, r, "/?modal=true", http.StatusSeeOther)
	} else {
		// if comment_id > 0 && comment_id <= len(comments) {
		fetch := fetchReactionByUserAndId(database, reactionsComments, user.Id, comment_id)
		if fetch.Value != 1 { // if not like
			if fetch.Value == -1 {
				// delete dislike in frontend and backend
				deleteRow(database, reactionsComments, fetch.Id)
				updateTableDislikes(database, commentsTable, allDislikes-1, comment_id)
			}
			// add like in frontend and backend
			updateTableLikes(database, commentsTable, allLikes+1, comment_id)
			addCommentsReactions(database, 1, user.Id, comment_id)
		} else {
			// delete like in frontend and backend
			updateTableLikes(database, commentsTable, allLikes-1, comment_id)
			deleteRow(database, reactionsComments, fetch.Id)
		}
		c, err := r.Cookie("last_page")
		if err != nil {
			fmt.Println("cookie err", err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		fmt.Println("REDIRECT TO", c.Value)
		http.Redirect(w, r, c.Value, http.StatusSeeOther)
		// }
	}
}

func dislikeComment(w http.ResponseWriter, r *http.Request) {
	comment_id, _ := strconv.Atoi(r.URL.Query().Get("comment_id"))
	user := welcome(w, r).User
	allLikes := fetchCommentByID(database, comment_id).Likes
	allDislikes := fetchCommentByID(database, comment_id).Dislikes

	reactionsComments := "commentsReactions"
	commentsTable := "comments"

	if user.Id == 0 { // if not login
		http.Redirect(w, r, "/?modal=true", http.StatusSeeOther)
	} else {
		// if comment_id > 0 && comment_id <= len(comments) {
		fetch := fetchReactionByUserAndId(database, reactionsComments, user.Id, comment_id)
		if fetch.Value != -1 { // if not dislike
			if fetch.Value == 1 {
				// delete like in frontend and backend
				deleteRow(database, reactionsComments, fetch.Id)
				updateTableLikes(database, commentsTable, allLikes-1, comment_id)
			}
			// add dislike in frontend and backend
			updateTableDislikes(database, commentsTable, allDislikes+1, comment_id)
			addCommentsReactions(database, -1, user.Id, comment_id)
		} else {
			// delete dislike in frontend and backend
			updateTableDislikes(database, commentsTable, allDislikes-1, comment_id)
			deleteRow(database, reactionsComments, fetch.Id)
		}
		c, err := r.Cookie("last_page")
		if err != nil {
			fmt.Println("cookie err", err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		fmt.Println("REDIRECT TO", c.Value)
		http.Redirect(w, r, c.Value, http.StatusSeeOther)
		// }
	}
}

func createPost(w http.ResponseWriter, r *http.Request) {
	data := welcome(w, r)
	r.ParseForm()

	msg := &Message{
		Threads: r.Form["threads"],
	}

	if !msg.ValidateThreads() {
		data := Data{Message: msg, Post: Post{Title: r.FormValue("title"), Content: r.FormValue("content")}, Threads: fetchAllThreads(database)}
		fmt.Println(data.Post)
		tmpl, _ := template.ParseFiles("static/template/newPost.html", "static/template/base.html")
		tmpl.Execute(w, data)
		return
	}

	fmt.Println("new post content", r.FormValue("content"))
	addPost(database, r.FormValue("title"), r.FormValue("content"), r.Form["threads"], data.User.Id, 0, 0)
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func createComment(w http.ResponseWriter, r *http.Request) {
	data := welcome(w, r)

	// msg := &Message{}

	// if !msg.ValidateComment() {
	// 	data := Data{Message: msg, Post: Post{Title: r.FormValue("title"), Content: r.FormValue("content")}, Threads: fetchAllThreads(database)}
	// 	fmt.Println(data.Post)
	// 	tmpl, _ := template.ParseFiles("static/template/newPost.html", "static/template/base.html")
	// 	tmpl.Execute(w, data)
	// 	return
	// }
	if data.User.Id == 0 { // if not login
		http.Redirect(w, r, "/?modal=true", http.StatusSeeOther)
	}

	fmt.Println("new comment content", r.FormValue("content"))
	post, _ := strconv.Atoi(r.FormValue("id"))
	addComment(database, r.FormValue("comment"), post, data.User.Id, 0, 0)

	c, err := r.Cookie("last_page")
	if err != nil {
		fmt.Println("cookie err", err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	fmt.Println("REDIRECT TO", c.Value)
	http.Redirect(w, r, c.Value, http.StatusSeeOther)
}

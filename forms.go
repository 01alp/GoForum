package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func like(w http.ResponseWriter, r *http.Request) { // for posts
	post_id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	posts := fetchAllPosts(database)
	user := welcome(w, r).User
	allLikes := posts[post_id-1].Likes
	allDislikes := posts[post_id-1].Dislikes
	postsTable := "postsReactions"

	if user.Id == 0 { // if not login
		http.Redirect(w, r, "/?modal=true", http.StatusSeeOther)
	} else {
		if post_id > 0 && post_id <= len(posts) {
			fetch := fetchReactionByUserAndId(database, postsTable, user.Id, post_id)
			if fetch.Value != 1 { // if not like
				if fetch.Value == -1 {
					// delete dislike in frontend and backend
					deleteRow(database, postsTable, fetch.Id)
					updatePostDislikes(database, allDislikes-1, post_id)
				}
				// add like in frontend and backend
				updatePostLikes(database, allLikes+1, post_id)
				addPostsReactions(database, 1, user.Id, post_id)
			} else {
				// delete like in frontend and backend
				updatePostLikes(database, allLikes-1, post_id)
				deleteRow(database, postsTable, fetch.Id)
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

func dislike(w http.ResponseWriter, r *http.Request) { // for posts
	post_id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	posts := fetchAllPosts(database)
	user := welcome(w, r).User
	allLikes := posts[post_id-1].Likes
	allDislikes := posts[post_id-1].Dislikes
	postsTable := "postsReactions"

	if user.Id == 0 { // if not login
		http.Redirect(w, r, "/?modal=true", http.StatusSeeOther)
	} else {
		if post_id > 0 && post_id <= len(posts) {
			fetch := fetchReactionByUserAndId(database, postsTable, user.Id, post_id)
			if fetch.Value != -1 { // if not dislike
				if fetch.Value == 1 {
					// delete like in frontend and backend
					deleteRow(database, postsTable, fetch.Id)
					updatePostLikes(database, allLikes-1, post_id)
				}
				// add dislike in frontend and backend
				updatePostDislikes(database, allDislikes+1, post_id)
				addPostsReactions(database, -1, user.Id, post_id)
			} else {
				// delete dislike in frontend and backend
				updatePostDislikes(database, allDislikes-1, post_id)
				deleteRow(database, postsTable, fetch.Id)
			}
		}
		c, err := r.Cookie("last_page")
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		http.Redirect(w, r, c.Value, http.StatusSeeOther)
	}
}

func createPost(w http.ResponseWriter, r *http.Request) {
	data := welcome(w, r)
	r.ParseForm()
	addPost(database, r.FormValue("title"), r.FormValue("content"), r.Form["threads"], data.User.Id)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

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
	allLikes := posts[post_id].Likes
	allDislikes := posts[post_id].Dislikes

	if user.Id == 0 { // if not login
		http.Redirect(w, r, "/?modal=true", http.StatusSeeOther)
	} else {
		if post_id > 0 && post_id <= len(posts) {
			fetch := fetchReactionByUserAndPost(database, user.Id, post_id)
			if fetch.Value != 1 { // if not like
				if fetch.Value == -1 {
					deleteRow(database, "postsReactions", fetch.Id)
					updatePostDislikes(database, allDislikes, post_id)
				}
				// add like in backend
				fmt.Println("Like added")
				updatePostLikes(database, allLikes+1, post_id)
				addPostsReactions(database, 1, user.Id, post_id)
				fmt.Println(fetchReactionByUserAndPost(database, user.Id, post_id))
				// TODO add like value in frontend
			} else {
				// delete like
				fmt.Println("Like deleted")
				updatePostLikes(database, allLikes, post_id)
				deleteRow(database, "postsReactions", fetch.Id)
				fmt.Println(fetchReactionByUserAndPost(database, user.Id, post_id))
			}
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		c, err := r.Cookie("last_page")
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		http.Redirect(w, r, c.Value, http.StatusSeeOther)
	}
}

func dislike(w http.ResponseWriter, r *http.Request) { // for posts
	post_id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	posts := fetchAllPosts(database)
	user := welcome(w, r).User
	allLikes := posts[post_id].Likes
	allDislikes := posts[post_id].Dislikes

	if user.Id == 0 { // if not login
		http.Redirect(w, r, "/?modal=true", http.StatusSeeOther)
	} else {
		if post_id > 0 && post_id <= len(posts) {
			fetch := fetchReactionByUserAndPost(database, user.Id, post_id)
			if fetch.Value != -1 { // if not dislike
				if fetch.Value == 1 {
					deleteRow(database, "postsReactions", fetch.Id)
					updatePostLikes(database, allLikes, post_id)
				}
				// add dislike in backend
				fmt.Println("Dislike added")
				updatePostDislikes(database, allDislikes+1, post_id)
				addPostsReactions(database, -1, user.Id, post_id)
				fmt.Println(fetchReactionByUserAndPost(database, user.Id, post_id))
				// TODO add dislike value in frontend
			} else {
				// delete dislike
				fmt.Println("Dislike deleted")
				updatePostDislikes(database, allDislikes, post_id)
				deleteRow(database, "postsReactions", fetch.Id)
				fmt.Println(fetchReactionByUserAndPost(database, user.Id, post_id))
			}
			http.Redirect(w, r, "/", http.StatusSeeOther)
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

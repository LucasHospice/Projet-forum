package main

import (
	pckg "Forum/Forum"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var (
	key   = []byte("inserer clé")
	store = sessions.NewCookieStore(key)
)

func main() {
	rr := mux.NewRouter()

	db := pckg.InitDatabase("forum.db")
	defer db.Close()

	rr.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./pages/404.html")
	})

	// PAS SUPPRIMER: DECOMMENTER POUR GENERER TABLES EXEMPLE //

	// userId1, _ := pckg.Create(db, "user", pckg.User{}, "akhy_deter", pckg.Encrypt("mdp"), "aeze@gmail.com", "6314134235235", strconv.Itoa(rand.Intn(9-1)+1), "1")
	// userId2, _ := pckg.Create(db, "user", pckg.User{}, "fifi_grognon", pckg.Encrypt("mdp"), "aeqze@gmail.com", "64235235", strconv.Itoa(rand.Intn(9-1)+1), "2")
	// userId3, _ := pckg.Create(db, "user", pckg.User{}, "koba", pckg.Encrypt("mdp"), "aedqze@gmail.com", "6488235235", strconv.Itoa(rand.Intn(9-1)+1), "2")
	// userId4, _ := pckg.Create(db, "user", pckg.User{}, "la_dd", pckg.Encrypt("mdp"), "aeqzed@gmail.com", "64235244435", strconv.Itoa(rand.Intn(9-1)+1), "3")

	// pckg.Create(db, "category", pckg.Category{}, "Santé", "pink")
	// pckg.Create(db, "category", pckg.Category{}, "Nostalgie", "purple")

	// pqrentPostId, _ := pckg.Create(db, "post", pckg.Post{}, "1 1 1 1 1 1 1 1 1 1 1", 1, "Je suis 1", 1, nil, userId1, "44/44", 0)
	// pckg.Create(db, "post", pckg.Post{}, "2 2 2 2 ", 0, "Je suis 2", 1, pqrentPostId, userId2, "15/13", 0)
	// postId3, _ := pckg.Create(db, "post", pckg.Post{}, "6 6 6 6 6 6 ", 1, "Je suis 666", 1, pqrentPostId, userId3, "15/13", 0)
	// pckg.Create(db, "post", pckg.Post{}, "3 3 3 3", 1, "Je suis 3", 1, nil, userId4, "9312", 0)
	// pckg.Create(db, "post", pckg.Post{}, "4 4 4 4", 0, "Je suis 4", 1, postId3, userId1, "9312", 0)

	// pqrentPostId2, _ := pckg.Create(db, "post", pckg.Post{}, "11 11 11 11", 1, "Je suis 11", 2, nil, userId1, "25/43", 0)
	// parent22, _ := pckg.Create(db, "post", pckg.Post{}, "22 22 22", 0, "Je suis 22", 1, pqrentPostId2, userId2, "35/96", 0)
	// pckg.Create(db, "post", pckg.Post{}, "22 22 22", 0, "Je suis 22", 1, parent22, userId2, "35/96", 0)
	// pckg.Create(db, "post", pckg.Post{}, "33 33 33", 0, "Je suis 33", 1, pqrentPostId2, userId1, "14/04", 0)

	// fs := http.FileServer(http.Dir("static/"))
	// http.Handle("/static/", http.StripPrefix("/static/", fs))
	fileServer := http.FileServer(http.Dir("./static"))
	rr.PathPrefix("/static").Handler(http.StripPrefix("/static", fileServer))

	userList := pckg.Get(db, "user", "user")
	userData := pckg.GetUserRows(userList)
	for i := 0; i < len(userData); i++ {
		rr.HandleFunc("/user/"+strconv.Itoa(userData[i].ID), func(w http.ResponseWriter, r *http.Request) {
			tmpl, _ := template.ParseFiles("./pages/user.html", "./templates/menu.html")
			tmpl.Execute(w, pckg.GetUserId(r))
		})
	}

	postList := pckg.Get(db, "post", "child")
	postArray := pckg.GetPostRows(postList)
	for i := 0; i < len(postArray); i++ {
		rr.HandleFunc("/topic/"+strconv.Itoa(postArray[i].ID), func(w http.ResponseWriter, r *http.Request) {
			tmpl, _ := template.ParseFiles("./pages/topic.html", "./templates/menu.html")
			tmpl.Execute(w, pckg.GetUserId(r))
		})
	}

	rr.HandleFunc("/", pckg.HandleHome)

	rr.HandleFunc("/signin", func(w http.ResponseWriter, r *http.Request) {
		pckg.HandleSignin(w, r, db)
	})

	rr.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		pckg.HandleLogin(db, w, r)
	})

	rr.HandleFunc("/logout", pckg.HandleLogout)

	rr.HandleFunc("/profil/{userID}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userId := vars["userID"]
		fmt.Print(userId)
		http.ServeFile(w, r, "./pages/user.html")
	})

	rr.HandleFunc("/mon-compte", func(w http.ResponseWriter, r *http.Request) {
		tmpl, _ := template.ParseFiles("./pages/mon-compte.html", "./templates/menu.html", "./templates/footer.html")
		tmpl.Execute(w, nil)
	})

	rr.HandleFunc("/post/{postID}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userId := vars["postID"]
		fmt.Print(userId)
		http.ServeFile(w, r, "./pages/topic.html")
	})

	rr.HandleFunc("/Apropos", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./pages/aproposde.html", "./templates/footer.html", "./templates/logo.html", "./templates/menu.html"))
		if r.Method != http.MethodPost {
			tmpl.Execute(w, nil)
			return
		}
	})

	rr.HandleFunc("/Cgu", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./pages/cgu.html", "./templates/footer.html", "./templates/logo.html", "./templates/menu.html"))
		if r.Method != http.MethodPost {
			tmpl.Execute(w, nil)
			return
		}
	})

	rr.HandleFunc("/Support", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./pages/support.html", "./templates/footer.html", "./templates/logo.html", "./templates/menu.html"))
		if r.Method != http.MethodPost {
			tmpl.Execute(w, nil)
			return
		}
		r.ParseForm()
		// objet := r.Form.Get("objetsupport")
		// corps := r.Form.Get("corpssupport")
		// result := objet + "\n" + corps
		// fmt.Println(result)
		http.Redirect(w, r, "/", http.StatusFound)
		//joindre le ticket a la bdd pour l'afficher dans le profil du superadmin
	})

	rr.HandleFunc("/Equipe", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./pages/equipe.html", "./templates/footer.html", "./templates/logo.html", "./templates/menu.html"))
		if r.Method != http.MethodPost {
			tmpl.Execute(w, nil)
			return
		}
	})

	rr.HandleFunc("/IsUpvoted", func(w http.ResponseWriter, r *http.Request) {
		var params pckg.IsUpvotedParams
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &params)
		if pckg.IsUpvoted(db, params.PostId, params.UserId) {
			w.Write([]byte("upvote"))
		} else {
			w.Write([]byte("upvoteUp"))
		}
	})

	rr.HandleFunc("/UpdateVote", func(w http.ResponseWriter, r *http.Request) {
		if pckg.IsLogged(r) {

			var params pckg.UpdateVoteParams
			body, _ := ioutil.ReadAll(r.Body)
			json.Unmarshal(body, &params)
			voteNumber := pckg.GetVoteById(db, params.PostId)
			notVoted := pckg.IsUpvoted(db, params.PostId, params.Id)
			if notVoted {
				voteNumber++
			} else {
				voteNumber--
			}
			w.Write([]byte(pckg.UpdateVotes(db, params.Table, strconv.Itoa(voteNumber), params.Field, params.PostId, params.Where, params.Id, notVoted)))
		} else {
			w.Write([]byte("isntLogged"))
		}
	})

	rr.HandleFunc("/new-post", func(w http.ResponseWriter, r *http.Request) {
		pckg.NewPost(w, r, db)
	})

	rr.HandleFunc("/addNewPost", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./pages/new-post.html")
	})

	// routes API

	rr.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		userList = pckg.Get(db, "user", "user")
		userData = pckg.GetUserRows(userList)
		json, _ := json.Marshal(userData)
		w.Write(json)
	})

	rr.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
		postList = pckg.Get(db, "post", "child")
		postArray = pckg.GetPostRows(postList)
		json, _ := json.Marshal(postArray)
		w.Write(json)
	})

	// rr.HandleFunc("/topics", func(w http.ResponseWriter, r *http.Request) {
	// 	topicList := pckg.Get(db, "post", "topic")
	// 	topicArray := pckg.GetPostRows(topicList)
	// 	json, _ := json.Marshal(topicArray)
	// 	w.Write(json)
	// })

	rr.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		topicList := pckg.Get(db, "category", "")
		a := pckg.GetCategoryRows(topicList)
		json, _ := json.Marshal(a)
		w.Write(json)
	})

	rr.HandleFunc("/upvote", func(w http.ResponseWriter, r *http.Request) {
		upvoteList := pckg.Get(db, "upvote", "")
		a := pckg.GetUpvoteRows(upvoteList)
		json, _ := json.Marshal(a)
		w.Write(json)
	})

	rr.HandleFunc("/cookies-data", func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "cookie-forum")
		auth := session.Values["authenticated"]
		var data pckg.UserData = pckg.UserData{}

		if auth != nil {
			// fmt.Println("-")
			// fmt.Println(auth.(string))
			// fmt.Println	("-")
			json.Unmarshal([]byte(auth.(string)), &data)
			data2 := pckg.UserDataConvert{Pseudo: data.Pseudo[0], Password: data.Password[0], UserID: data.UserID[0]}
			cookie, _ := json.Marshal(data2)
			w.Write(cookie)
		}
	})

	http.ListenAndServe(":8080", rr)
}

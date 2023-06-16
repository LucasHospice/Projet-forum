package Forum

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func FormatDate() string {

	// p := fmt.Println
	now := time.Now()
	date := now.Format("01/02/2006")
	// p(date)
	return date
}

func NewPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	var data UserData = UserData{}

	session, _ := store.Get(r, "cookie-forum")
	auth := session.Values["authenticated"]

	data2 := UserDataConvert{}

	if auth != nil {
		json.Unmarshal([]byte(auth.(string)), &data)
		data2 = UserDataConvert{data.Pseudo[0], data.Password[0], data.UserID[0]}
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", 500)
	}

	date := FormatDate()
	userID, _ := strconv.Atoi(data2.UserID)

	title := r.Form.Get("title")
	content := r.Form.Get("post-body")
	category := 1
	// category := r.Form.Get("category")

	fmt.Println("user:", userID, "titre:", title, "content:", content, "category:", category)
	// on recupere le form, et on le fout dans la fonction create Post.

	newt, _ := Create(db, "post", Post{}, content, 1, title, category, nil, userID, date, 0)
	res := strconv.Itoa(int(newt))

	http.ServeFile(w, r, "./pages/topic/"+res)
}

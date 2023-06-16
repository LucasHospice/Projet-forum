package Forum

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"text/template"

	"github.com/gorilla/sessions"
)

var (
	key   = []byte("inserer clé")
	store = sessions.NewCookieStore(key)
)

type UserData struct {
	Pseudo   [1]string `json:"pseudo"`
	Password [1]string `json:"password"`
	UserID   [1]string `json:"user_id"`
}

type signinParams struct {
	Pseudo    string `json:"Pseudo"`
	Mail      string `json:"Email"`
	Number    string `json:"Number"`
	Password  string `json:"Password"`
	Password2 string `json:"Password2"`
}

type loginParams struct {
	Pseudo   string `json:"Pseudo"`
	Password string `json:"Password"`
}

type UserDataConvert struct {
	Pseudo   string `json:"pseudo"`
	Password string `json:"password"`
	UserID   string `json:"user_id"`
}

func HandleHome(w http.ResponseWriter, r *http.Request) {
	var data UserData = UserData{}

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	session, _ := store.Get(r, "cookie-forum")
	auth := session.Values["authenticated"]

	tmpl, _ := template.ParseFiles("./pages/accueil.html", "./templates/menu.html", "./templates/filtre.html")
	data2 := UserDataConvert{}

	if auth != nil {
		// fmt.Println("-")
		// fmt.Println(auth.(string))
		// fmt.Println("-")
		json.Unmarshal([]byte(auth.(string)), &data)
		data2 = UserDataConvert{data.Pseudo[0], data.Password[0], data.UserID[0]}
	}
	tmpl.Execute(w, data2)
}

func HandleSignin(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var params signinParams
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &params)
	fmt.Println(params)

	// if err := r.ParseForm(); err != nil {
	// 	http.Error(w, "Error parsing form", 500)
	// }

	if r.URL.Path != "/signin" {
		http.NotFound(w, r)
		return
	}

	if params.Password != params.Password2 {
		w.Write([]byte("Erreur: Mots de passes non identique"))
		return
	}
	if len(params.Password) < 6 {
		w.Write([]byte("Erreur: Le mot de passe doit contenir au moins 6 caractères"))
		return
	}
	if _, err := strconv.Atoi(params.Number); err != nil || len(params.Number) != 10 {
		w.Write([]byte("Erreur: Le numéro de téléphone est invalide"))
		return
	}
	if params.Pseudo == "" || params.Password == "" || (params.Mail == "" && params.Number == "") {
		w.Write([]byte("Erreur: Veuillez renseigner tous les champs obligatoire"))
		return
	}

	_, err := Create(db, "user", User{}, params.Pseudo, Encrypt(params.Password), params.Mail, params.Number, strconv.Itoa(rand.Intn(9-1)+1), "1")
	if err != nil {
		fmt.Println("error on user creation " + err.Error())
	}
}

func HandleLogin(db *sql.DB, w http.ResponseWriter, r *http.Request) {

	// if err := r.ParseForm(); err != nil {
	// 	http.Error(w, "Error parsing form", 500)
	// }

	// pseudo := r.Form.Get("pseudo")
	// password := Encrypt(r.Form.Get("password"))

	var params loginParams
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &params)
	fmt.Println(params)
	if r.URL.Path != "/login" {
		http.NotFound(w, r)
		return
	}

	fmt.Println(params.Pseudo + " zz " + params.Password)
	loginQuery := db.QueryRow("SELECT pseudo, password, id FROM user WHERE Pseudo=? AND Password=? ", params.Pseudo, Encrypt(params.Password))
	fmt.Println(loginQuery)
	u := UserData{}
	connexion := loginQuery.Scan(&u.Pseudo[0], &u.Password[0], &u.UserID[0])

	if connexion != nil {
		fmt.Println(connexion)
		fmt.Println("error: Wrong password or username. Please try again.")
		w.Write([]byte("Erreur: Pseudo ou mot de passe erroné"))
	} else {
		fmt.Println("utilisateur trouvé dans la bdd")
		// tmpl, _ := template.ParseFiles("./pages/accueil.html")

		session, err := store.Get(r, "cookie-forum")

		if err != nil {
			fmt.Println(err.Error())
		}

		res, err := json.Marshal(u)
		if err != nil {
			fmt.Println(err.Error())
		}
		session.Values["authenticated"] = string(res)
		session.Save(r, w)
		// http.Redirect(w, r, "/", http.StatusFound)
	}
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/logout" {
		http.NotFound(w, r)
		return
	}

	session, _ := store.Get(r, "cookie-forum")
	session.Values["authenticated"] = nil
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}

func Encrypt(pwd string) string {
	salt := "Jessica Alba"
	hasher := md5.New()
	hasher.Write([]byte(pwd + salt))
	a := hex.EncodeToString(hasher.Sum(nil))
	return a
}

func IfExists(db *sql.DB, target string, table string, field string) {

	req := "SELECT * FROM " + table + " WHERE " + field + " LIKE " + "'%" + target + "%'"
	res, err := db.Query(req)
	res.Scan(&target, &table, &field)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	// fmt.Printf("%v", res)
}

func GetUserId(r *http.Request) string {
	var data UserData = UserData{}
	session, _ := store.Get(r, "cookie-forum")
	auth := session.Values["authenticated"]
	if auth != nil {
		json.Unmarshal([]byte(auth.(string)), &data)
	}
	return data.UserID[0]
}

func IsLogged(r *http.Request) bool {
	var data UserData = UserData{}
	session, _ := store.Get(r, "cookie-forum")
	auth := session.Values["authenticated"]

	if auth != nil {
		json.Unmarshal([]byte(auth.(string)), &data)
		fmt.Println(data)
		return true
	} else {
		return false
	}
}

// func checkRegister(db *sql.DB, pseudo string, mail string, number string) bool {
// 	var dbPseudo string
// 	var dbMail string
// 	var dbNumber string
// 	var UserExists bool
// 	// query := db.QueryRow("SELECT pseudo, mail, number FROM user WHERE pseudo=?", pseudo, email, number).Scan(&dbPseudo, &dbEmail, &dbNumber)
// 	query := db.QueryRow("SELECT pseudo, mail, number FROM user WHERE pseudo=?", pseudo, mail, number).Scan(&dbPseudo, &dbMail, &dbNumber)
// 	if dbPseudo == "" {
// 		fmt.Println("user can be created:")
// 		UserExists = false
// 		fmt.Println("ok", query)
// 		return UserExists
// 	} else {
// 		fmt.Println("user already found ! ", dbPseudo, dbNumber, dbNumber)
// 		fmt.Println("Credentials already exists", dbPseudo, dbMail, dbNumber, "please login")
// 		UserExists = true
// 		fmt.Println("nope", query)
// 		return UserExists
// 	}
// }

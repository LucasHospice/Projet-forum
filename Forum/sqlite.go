package Forum

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type Category struct {
	ID    int
	Name  string
	Color string
}

type User struct {
	ID         int
	Pseudo     string
	Password   string
	Mail       string
	Number     string
	ProfilePic string
	Level      string
}

type Post struct {
	ID           int
	Content      string
	IsTopic      int
	Title        sql.NullString
	Category     sql.NullInt64
	ParentPostId sql.NullInt64
	UserId       int
	Date         string
	UpVote       int
}

type UpdateVoteParams struct {
	Table  string
	Id     string
	Field  string
	Value  string
	Where  string
	PostId string
	UserId string
}

type IsUpvotedParams struct {
	UserId string
	PostId string
}

type Upvote struct {
	ID     string
	UserId string
	PostId string
}

func GetUserRows(rows *sql.Rows) []User {
	final := make([]User, 0)
	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.Pseudo, &u.Password, &u.Mail, &u.Number, &u.ProfilePic, &u.Level)
		if err != nil {
			log.Fatal(err)
		}
		final = append(final, u)
	}
	return final
}

func GetPostRows(rows *sql.Rows) []Post {
	final := make([]Post, 0)
	for rows.Next() {
		var u Post
		err := rows.Scan(&u.ID, &u.Content, &u.IsTopic, &u.Title, &u.Category, &u.ParentPostId, &u.UserId, &u.Date, &u.UpVote)
		if err != nil {
			log.Fatal(err)
		}
		final = append(final, u)
		// fmt.Println(u)
	}
	return final
}

func GetCategoryRows(rows *sql.Rows) []Category {
	final := make([]Category, 0)
	for rows.Next() {
		// fmt.Println(rows)
		var u Category
		err := rows.Scan(&u.ID, &u.Name, &u.Color)
		if err != nil {
			log.Fatal(err)
		}
		final = append(final, u)
		// fmt.Println(u)
	}
	return final
}

func GetUpvoteRows(rows *sql.Rows) []Upvote {
	final := make([]Upvote, 0)
	for rows.Next() {
		// fmt.Println(rows)
		var u Upvote
		err := rows.Scan(&u.ID, &u.UserId, &u.PostId)
		if err != nil {
			log.Fatal(err)
		}
		final = append(final, u)
		// fmt.Println(u)
	}
	return final
}

func InitDatabase(database string) *sql.DB {
	db, err := sql.Open("sqlite3", database)

	if err != nil {
		log.Fatal(err)
	}

	sqlStmnt := `
				PRAGMA foreign_keys = ON;
				CREATE TABLE IF NOT EXISTS user (
					ID INTEGER PRIMARY KEY AUTOINCREMENT,
					Pseudo STRING NOT NULL UNIQUE,
					Password STRING NOT NULL,
					Mail STRING UNIQUE,
					Number STRING UNIQUE,
					ProfilePic STRING,
					Level STRING
				);
				CREATE TABLE IF NOT EXISTS category (
					ID INTEGER PRIMARY KEY AUTOINCREMENT,
					Name STRING NOT NULL,
					Color STRING NOT NULL
				);
				CREATE TABLE IF NOT EXISTS post (
					ID INTEGER PRIMARY KEY AUTOINCREMENT,
					Content STRING NOT NULL,
					IsTopic INTEGER NOT NULL,
					Title STRING,
					Category INTEGER NOT NULL,
					ParentPostId INTEGER,
					UserId INTEGER NOT NULL ,
					Date STRING NOT NULL,
					UpVote STRING NOT NULL,
					FOREIGN KEY (UserId) REFERENCES user(ID) ,
					FOREIGN KEY (ParentPostId) REFERENCES post(ID),
					FOREIGN KEY (Category) REFERENCES category(ID)
				);
				CREATE TABLE IF NOT EXISTS upvote (
					ID INTEGER PRIMARY KEY AUTOINCREMENT,
					UserID INTEGER NOT NULL,
					PostID INTEGER NOT NULL,
					FOREIGN KEY (UserID) REFERENCES user(ID),
					FOREIGN KEY (PostID) REFERENCES post(ID)
				)
				`
	_, err = db.Exec(sqlStmnt)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func ParseTable(model interface{}, table string) (a string) {

	result := ""
	e := reflect.ValueOf(model)
	for i := 0; i < e.NumField(); i++ {
		varName := e.Type().Field(i).Name
		if varName != "ID" {
			result += string(varName) + ", "
		}
	}
	result = result[:len(result)-2]
	return result
}

func Create(db *sql.DB, table string, model interface{}, t ...interface{}) (int64, error) {

	result := "INSERT INTO " + table + " (" + ParseTable(model, table) + ")" + " VALUES " + "("
	for i := 0; i < len(t); i++ {
		result += "?, "
	}
	result = result[:len(result)-2]
	result += ")"

	request, err := db.Exec(result, t...)

	if err != nil {
		fmt.Println(err)
		return -1, nil
	}

	return request.LastInsertId()
}

func Get(db *sql.DB, table string, mode string) *sql.Rows {
	query := "SELECT * FROM " + table
	if mode == "topic" {
		query += " WHERE isTopic = 1"
	}
	result, _ := db.Query(query)
	return result
}

func DeletePostById(db *sql.DB, id string) (int64, error) {
	result, _ := db.Exec(`DELETE FROM post WHERE postId = ?`, id)
	return result.LastInsertId()
}

func UpdateVotes(db *sql.DB, table string, value string, field string, value2 string, field2 string, userId string, notVoted bool) string {
	_, err := db.Exec("UPDATE " + table + " SET " + field + " = " + value + " WHERE " + field2 + " = " + value2 + ";")
	if err != nil {
		fmt.Println("Update Post upvote error:")
	}
	if notVoted {
		_, err2 := db.Exec("INSERT INTO upvote (UserID, PostID) VALUES (?, ?)", userId, value2)
		if err2 != nil {
			fmt.Println("Insert upvote error")
		}
	} else {
		_, err2 := db.Exec("DELETE FROM upvote WHERE UserID = "+userId+" AND PostID = "+value2+";", userId, value2)
		if err2 != nil {
			fmt.Println("Delete upvote error")
		}
	}
	return value
}

func IsUpvoted(db *sql.DB, postId string, userId string) bool {
	query := "SELECT ID FROM upvote WHERE UserID = " + userId + " AND PostID = " + postId + ";"
	rows := db.QueryRow(query)
	var u string
	rows.Scan(&u)
	if u != "" {
		return false
	} else {
		return true
	}
}

func GetVoteById(db *sql.DB, id string) int {
	query := "SELECT UpVote FROM post WHERE ID =" + id + ";"
	var number string
	db.QueryRow(query).Scan(&number)
	result, _ := strconv.Atoi(number)
	return result
}

// reriter à

// func CreateUser(db *sql.DB, pseudo string, password string, mail string, number int, profilePic string) (int64, error) {
// 	result, _ := db.Exec(`INSERT INTO user (pseudo, password, mail, number, profilePic) VALUES (?,?,?,?,?)`, pseudo, password, mail, number, profilePic)
// 	return result.LastInsertId()
// }

// func CreateTopic(db *sql.DB, content string, userId int, isTopic int, titre string, categorie string) (int64, error) {
// 	result, _ := db.Exec(`INSERT INTO post (content, userId, isTopic, title, category) VALUES (?,?,?,?,?)`, content, userId, isTopic, titre, categorie)
// 	return result.LastInsertId()
// }

// func CreatePost(db *sql.DB, content string, userId int, isTopic int, parentPostId int) (int64, error) {
// 	result, _ := db.Exec(`INSERT INTO post (content, userId, isTopic, parentPostId) VALUES (?,?,?,?)`, content, userId, isTopic, parentPostId)
// 	return result.LastInsertId()
// }

// func GetTopic(db *sql.DB, table string) *sql.Rows {
// 	query := "SELECT * FROM " + table + " WHERE isTopic = 1"
// 	result, _ := db.Query(query)
// 	return result
// }

// func GetTable(db *sql.DB, table string) *sql.Rows {
// 	query := "SELECT * FROM " + table
// 	result, _ := db.Query(query)
// 	return result
// }

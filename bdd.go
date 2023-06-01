package back

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func OpenBDD() *sql.DB {
	database, bdderr := sql.Open("sqlite3", "./BDD.db")
	if bdderr != nil {

		log.Fatal(bdderr.Error())

	}
	return database
}

var users []user
var posts []post
var comments []comment
var likes []like

func InitBDD() {
	database := OpenBDD()
	defer database.Close()
	tmp := `
	CREATE TABLE IF NOT EXISTS "user" (
		"id"					INTEGER NOT NULL UNIQUE,
		"age"					INTEGER NOT NULL,
		"firstname_user"		VARCHAR(20) NOT NULL,
		"lastname_user"			VARCHAR(30) NOT NULL,
		"email_user"			VARCHAR(50) NOT NULL UNIQUE,
		"password_hashed_user"	VARCHAR(45) NOT NULL,
		"pseudo_user"			ARCHAR(20) NOT NULL UNIQUE,
		PRIMARY KEY("id" AUTOINCREMENT)
		
	);

	CREATE TABLE IF NOT EXISTS "post" (
		"id"				INTEGER NOT NULL UNIQUE,
		"id_user" 		 	INTEGER NOT NULL UNIQUE REFERENCES user(id),
		"title_post" 		VARCHAR(50) NOT NULL,
		"content_post" 		LONGTEXT NOT NULL,
		PRIMARY KEY("id" AUTOINCREMENT)

	);

	CREATE TABLE IF NOT EXISTS "comment" (
		"id"	INTEGER NOT NULL UNIQUE,
		"id_post"  INTEGER NOT NULL UNIQUE REFERENCES post(id),
		"id_user"  INTEGER NOT NULL UNIQUE REFERENCES user(id),
		"content_comment" 	LONGTEXT NOT NULL,
		PRIMARY KEY("id" AUTOINCREMENT)
	);

	CREATE TABLE IF NOT EXISTS "like" (
		"id"	INTEGER NOT NULL UNIQUE,
		"id_post"  INTEGER NOT NULL UNIQUE REFERENCES post(id),
		"effet"   BOOLEAN, 
		PRIMARY KEY("id" AUTOINCREMENT)
	);
	`

	_, bdderr := database.Exec(tmp)
	if bdderr != nil {
		log.Fatal(bdderr.Error())
	}

}

func GetAllUsers() {
	var user user
	var id int = 0
	var age int = 0
	var firstname string = ""
	var lastname string = ""
	var email string = ""
	var password string = ""
	var pseudo string = ""
	users = append(users, user)
	database := OpenBDD()
	rows, err := database.Query("SELECT * FROM user")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&id, &age, &firstname, &lastname, &email, &password, &pseudo)
		if err != nil {
			log.Fatal(err)
		}
		user.id = id
		user.age = age
		user.firstname_user = firstname
		user.lastname_user = lastname
		user.email_user = email
		user.password_hashed_user = password
		user.pseudo_user = pseudo
		users = append(users, user)
	}
}

func GetAlPosts() {
	var post post
	var id int = 0
	var id_user int = 0
	var title_post string = ""
	var content_post string = ""
	posts = append(posts, post)
	database := OpenBDD()
	rows, err := database.Query("SELECT * FROM post")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&id, &id_user, &title_post, &content_post)
		if err != nil {
			log.Fatal(err)
		}
		post.id = id
		post.id_user = id_user
		post.title_post = title_post
		post.content_post = content_post
		posts = append(posts, post)
	}
}
func GetAlComments() {
	var comment comment
	var id int
	var id_post int
	var id_user int
	var content_comment string
	database := OpenBDD()
	rows, err := database.Query("SELECT * FROM comment")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&id, &id_post, &id_user, &content_comment)
		if err != nil {
			log.Fatal(err)
		}
		comment.id = id
		comment.id_post = id_post
		comment.id_user = id_user
		comment.content_comment = content_comment
		comments = append(comments, comment)

	}
}
func GetAlLikes() {
	var like like
	var id int
	var id_post int
	var effect bool
	database := OpenBDD()
	rows, err := database.Query("SELECT * FROM like")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&id, &id_post, &effect)
		if err != nil {
			log.Fatal(err)
		}
		like.id = id
		like.id_post = id_post
		like.effet = effect
		likes = append(likes, like)
	}
}

func AddUser(age int, firstname string, lastname string, email string, password string, pseudo string) {
	database := OpenBDD()
	password = HashPassword(password)
	fmt.Println(password)
	if age < 13 {
		log.Fatal("grandis et tu pourras parler")
	}
	statement, BDDerr := database.Prepare(`INSERT INTO user(age, firstname_user, lastname_user, email_user, password_hashed_user, pseudo_user) VALUES(?,?,?,?,?,?);`)
	if BDDerr != nil {
		log.Fatal("ERROR 500 : error Prepare new user \n ", BDDerr)
	}
	_, BDDerr = statement.Exec(strconv.Itoa(age), firstname, lastname, email, password, pseudo)
	if BDDerr != nil {
		log.Fatal("ERROR 500 : error exec new user \n ", BDDerr)
	}
	defer database.Close()
}

func HashPassword(password string) string {
	bytes, hashingErr := bcrypt.GenerateFromPassword([]byte(password), 14)
	if hashingErr != nil {
		log.Fatal(hashingErr)
	}
	return string(bytes)
}

func CheckPasswordHash(password string, hash string) bool {
	hashingErr := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return hashingErr == nil
}

package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

//User is a user
type User struct {
	gorm.Model
	Name  string
	Email string
}

func main() {

	//env取得
	err := godotenv.Load(fmt.Sprintf("./%s.env", os.Getenv("GO_ENV")))
	if err != nil {
		//env読めなかった時の処理
	}

	db := sqlConnect()
	db.AutoMigrate(&User{})
	defer db.Close()

	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")

	router.GET("/", func(ctx *gin.Context) {
		db := sqlConnect()
		var users []User
		db.Order("created_at asc").Find(&users)
		defer db.Close()

		ctx.HTML(200, "index.html", gin.H{
			"users": users,
		})
	})

	router.POST("/new", func(ctx *gin.Context) {
		db := sqlConnect()
		name := ctx.PostForm("name")
		email := ctx.PostForm("email")
		fmt.Println("create user " + name + "with email " + email)
		db.Create(&User{Name: name, Email: email})
		defer db.Close()
		ctx.Redirect(302, "/")
	})

	router.POST("/delete/:id", func(ctx *gin.Context) {
		db := sqlConnect()
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("id is not a number")
		}
		var user User
		db.First(&user, id)
		db.Delete(&user)
		defer db.Close()

		ctx.Redirect(302, "/")
	})

	router.Run()
}

func sqlConnect() (database *gorm.DB) {
	DBMS := "mysql"
	USER := os.Getenv("MYSQL_USER")
	PASS := os.Getenv("MYSQL_PASSWORD")
	PROTOCOL := "tcp(db:3306)"
	DBNAME := os.Getenv("MYSQL_DATABASE")

	CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME + "?charset=utf8&parseTime=true&loc=Asia%2FTokyo"

	count := 0
	db, err := gorm.Open(DBMS, CONNECT)
	if err != nil {
		for {
			if err == nil {
				fmt.Println("")
				break
			}
			fmt.Print(".")
			time.Sleep(time.Second)
			count++
			if count > 180 {
				fmt.Println("")
				fmt.Println("DB connection failed")
				panic(err)
			}
			db, err = gorm.Open(DBMS, CONNECT)
		}
	}
	fmt.Println("DB connection succeeded")

	return db
}

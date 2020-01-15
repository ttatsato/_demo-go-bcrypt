package main

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"log"
)

/**
 * User{}
 * gormでORM
 */
type User struct {
	Name string
	LoginId string `gorm:"unique;not null"`
	Password string
}

func ConnectDB() (*gorm.DB, error) {
	db, err := gorm.Open("mysql", "root:@(localhost)/dekin_list?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Println("データベースへの接続に失敗しました")
		panic("データベースへの接続に失敗しました")
	}
	return db, err
}

func migrate () {
	conn, _ := ConnectDB()
	defer conn.Close()

	loginId := "これはIdです"
	demoPassWord, err := generateHash("thisIsPassword")
	if err != nil {
		panic(err)
	}
	// テーブルをdrop
	if conn.HasTable(&User{}) {
		conn.DropTable(&User{})
	}
	conn.AutoMigrate(&User{})
	if err := conn.Create(User{Name: "テスト太郎", LoginId: loginId, Password: demoPassWord}).Error; err != nil {
		// error handling...
		log.Fatal(err)
	}
}

func generateHash(str string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		return "", nil
	}
	return string(hash), err
}

func login (loginId string, password string) {
	log.Println("func login():" + "loginId = " + loginId + "  +  password =" + password)

	conn, _ := ConnectDB()
	defer conn.Close()
	var user User
	if err := conn.Where("login_id = ?", loginId).Find(&user).Error; err != nil {
		log.Println("このログインIDは存在しません。: " + loginId)
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err == nil {
		// 成功
		log.Println("ログイン成功")
	} else {
		// 失敗
		log.Println("ログイン失敗")
	}
}


func main() {
	e := echo.New()
	migrate()

	// ログインテスト
	login("間違ったID", "間違ったpassword")
	login("これはIdです", "thisIsPassword")

	e.Logger.Fatal(e.Start(":1333"))
}
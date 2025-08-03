package models

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
func ConnectDB() error{
	err := godotenv.Load(".env")
	if err != nil{
		return err
	}
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")	

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",  dbUser, dbPass, dbHost, dbPort, dbName)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Connect to database failed, %v\n", err)
	} else {
		log.Printf("Connect to database success, host: %s, port: %s, user: %s, dbname: %s\n", dbHost, dbPort, dbUser, dbName)
	}
	// 迁移数据表
	err = DB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci").AutoMigrate(&User{}, &Post{}, &Comment{})
    if err != nil {
        panic(err)
    }

	return nil
}
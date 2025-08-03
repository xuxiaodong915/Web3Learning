package token

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(user_id uint) (string, error) {
	token_lifespan, err :=strconv.Atoi(os.Getenv("TOKEN_HOUR_LIFESPAN"))
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = user_id
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(token_lifespan)).Unix()
	token :=jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func TokenValid(c *gin.Context) error {
	tokenString := ExtractToken(c)
	fmt.Println(tokenString)
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error){
		if _,ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil,fmt.Errorf("unexpected signing method:%v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")),nil
	})
	if err != nil {
		return err
	}
	return nil
}

func ExtractToken(c *gin.Context) string {
	return strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
}

func ExtractTokenID(c *gin.Context) (uint, error) {
	tokenString := ExtractToken(c)
	log.Println(tokenString)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error){
		if _,ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil,fmt.Errorf("unexpected signing method:%v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		log.Println(err)
		return 0, err
	}
	claims , ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		uid, err := strconv.ParseInt(fmt.Sprintf("%.0f",claims["user_id"]),10,32)
		if err != nil {
			return 0,err
		}
		return uint(uid),nil
	}
	return 0,nil
}
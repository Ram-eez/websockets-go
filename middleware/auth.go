package middleware

import (
	"fmt"
	"time"
	"websockets/manager"

	"github.com/golang-jwt/jwt/v5"
)

var SecretKey = []byte("secret-key")

func CreateToken(user *manager.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"userID":   user.ID,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		fmt.Println("err signing the jwt token")
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) error {

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return SecretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

func GetUserFromToken(tokenString *jwt.Token) {
	claims := tokenString.Claims.(jwt.MapClaims)

	var user User
	username := claims["username"].(string)
	userID := claims["userID"].(string)
	user.Username = username
	user.ID = userID
}

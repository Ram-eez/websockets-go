package middleware

import (
	"errors"
	"fmt"
	"time"
	"websockets/models"

	"github.com/golang-jwt/jwt/v5"
)

var SecretKey = []byte("secret-key")

func CreateToken(user *models.User) (string, error) {
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

func GetUserFromToken(tokenString string) (*models.User, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return models.SecretKey, nil
	})
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("could not get claims from the token", err)
	}

	var user models.User

	switch {
	case token.Valid:
		user.Username = claims["username"].(string)
		user.ID = claims["userID"].(string)
		user.Password = ""
	case errors.Is(err, jwt.ErrTokenMalformed):
		fmt.Println("That's not even a token")
		return nil, err
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		fmt.Println("Invalid signature")
		return nil, err
	case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
		fmt.Println("Timing is everything")
		return nil, err
	default:
		fmt.Println("Couldn't handle this token:", err)
		return nil, err
	}

	return &user, nil
}

package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/teresa-solution/api-gateway/internal/middleware"
)

func main() {
	claims := jwt.MapClaims{
		"sub":  "test-user",
		"role": "admin",
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
		"iat":  time.Now().Unix(),
	}

	token, err := middleware.GenerateTestJWT(claims)
	if err != nil {
		fmt.Printf("Failed to generate token: %v\n", err)
		return
	}

	fmt.Printf("Generated JWT: %s\n", token)
}

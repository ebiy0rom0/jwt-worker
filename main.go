package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"jwt-worker/config"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprint(os.Stderr, err)
	}
}

func run() error {
	r := echo.New()
	r.Use(middleware.Logger(), middleware.CORS(), middleware.Recover())

	r.GET("/verify", func(c echo.Context) error {
		h := c.Request().Header.Get("Authorization")
		tokenString := strings.TrimPrefix(h, "Bearer ")

		if err := decodeToken(tokenString); err != nil {
			c.String(http.StatusUnauthorized, "invalid jwt token format")
			return err
		}
		token, err := verifyToken(tokenString)
		if err != nil {
			c.String(http.StatusUnauthorized, "unverified token")
			return err
		}

		sub, err := token.Claims.GetSubject()
		if err != nil {
			return err
		}
		c.String(http.StatusOK, sub)
		return nil
	})

	r.Start(":8080")
	return nil
}

func decodeToken(tokenString string) error {
	// Split the token into its three parts
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return fmt.Errorf("invalid token format")
	}

	// Decode the header
	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return err
	}
	var header map[string]interface{}
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return err
	}

	// Decode the payload
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return err
	}
	var payload map[string]interface{}
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return err
	}
	return nil
}

func verifyToken(tokenString string) (*jwt.Token, error) {
	// Define the key to use for signing the token
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte(config.C.Supabase.Secret), nil
	}

	// Parse the token
	token, err := jwt.Parse(tokenString, keyFunc)
	if err != nil {
		return nil, err
	}

	// Verify the token
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}

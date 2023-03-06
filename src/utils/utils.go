package utils

import (
	"encoding/base64"
	"finances_manager_go/models"
	"fmt"
	"math/rand"
	"net/http"
	"time"
	"unicode"

	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

func ReturnJsonResponse(res http.ResponseWriter, httpCode int, resMessage []byte) {
	res.Header().Set("Content-type", "application/json")
	res.WriteHeader(httpCode)
	res.Write(resMessage)
}

func MethodValidation(res http.ResponseWriter, req *http.Request, expectedMethod string) {
	if expectedMethod != "POST" && expectedMethod != "GET" && expectedMethod != "DELETE" &&
		expectedMethod != "UPDATE" && expectedMethod != "PATCH" {
		msg := []byte(`{
			"success": false,
			"message":"Invalid HTTP method expected"
		}`)

		ReturnJsonResponse(res, http.StatusMethodNotAllowed, msg)
	}
	if req.Method != expectedMethod {
		msg := []byte(`{
			"success": false,
			"message":"Invalid HTTP method"
		}`)

		ReturnJsonResponse(res, http.StatusMethodNotAllowed, msg)
	}
}

func ValidatePwd(s string) bool {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	if len(s) >= 7 {
		hasMinLen = true
	}
	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

func GenID[Model models.Income | models.Expense | models.User](db *gorm.DB, model Model) (id uint) {
	randomId := rand.Uint32()
	if err := db.First(&model, randomId).Error; err != gorm.ErrRecordNotFound {
		return GenID(db, model)
	}

	return uint(randomId)
}

func GenToken(duration int, payload interface{}, privateKey string) (string, error) {
	decodedPrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return "", fmt.Errorf("could not decode: %w", err)
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)

	if err != nil {
		return "", fmt.Errorf("validate: parse key: %w", err)
	}

	nanosecondsToMinutes := 60000000000

	var tokenDuration time.Duration = time.Duration(duration) * time.Duration(nanosecondsToMinutes)

	now := time.Now().UTC()

	claims := make(jwt.MapClaims)
	claims["sub"] = payload
	claims["exp"] = now.Add(tokenDuration).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		return "", fmt.Errorf("Error: could not generate token: %w", err)
	}

	return token, nil
}

func ValidateToken(token string, publicKey string) (interface{}, error) {
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return nil, fmt.Errorf("could not decode: %w", err)
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)
	if err != nil {
		return "", fmt.Errorf("validate: parse key: %w", err)
	}

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}
		return key, nil
	})

	if err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("validate: invalid token")
	}

	return claims["sub"], nil
}

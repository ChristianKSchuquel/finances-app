package middlewares

import (
	"finances_manager_go/utils"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
)

var accessTokenPublicKey = os.Getenv("ACCESS_TOKEN_PUBLIC_KEY")

func Auth(handler func(res http.ResponseWriter, req *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.Header["access_token"] == nil {
			msg := []byte(`{
				"success": false,
				"msg": "Token header empty"
			}`)

			utils.ReturnJsonResponse(res, http.StatusUnauthorized, msg)
		}

		token, err := jwt.Parse(req.Header["access_token"][0], func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodRSA)
			if !ok {
				msg := []byte(`{
					"success": false,
					"msg": "invalid token"
				}`)
				utils.ReturnJsonResponse(res, http.StatusUnauthorized, msg)
				return nil, nil
			}
			return "", nil
		})
		if err != nil {
			msg := []byte(`
				"success": false,
				"msg": "Error while parsing token"
			`)
			utils.ReturnJsonResponse(res, http.StatusUnauthorized, msg)
			return
		}

		if token.Valid {
			handler(res, req)
		} else {
			msg := []byte(`{
				"success": false,
				"msg": "Invalid token"
			}`)

			utils.ReturnJsonResponse(res, http.StatusUnauthorized, msg)
			return
		}
	})
}

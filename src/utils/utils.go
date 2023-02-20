package utils

import "net/http"

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

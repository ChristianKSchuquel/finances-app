package handlers

import (
	"encoding/json"
	"finances_manager_go/models"
	"finances_manager_go/utils"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type APIEnv struct {
	DB *gorm.DB
}

var accessTokenPrivateKey = os.Getenv("ACCESS_TOKEN_PRIVATE_KEY")
var accessTokenPublicKey = os.Getenv("ACCESS_TOKEN_PUBLIC_KEY")
var refreshTokenPrivateKey = os.Getenv("REFRESH_TOKEN_PRIVATE_KEY")
var refreshTokenPublicKey = os.Getenv("REFRESH_TOKEN_PUBLIC_KEY")
var refreshTokenMaxAge = os.Getenv("REFRESH_TOKEN_MAXAGE")
var accessTokenMaxAge = os.Getenv("ACCESS_TOKEN_MAXAGE")

//=========================================================================
// POST: /signup
//=========================================================================

func (a *APIEnv) CreateUser(res http.ResponseWriter, req *http.Request) {
	utils.MethodValidation(res, req, "POST")

	var user models.User

	payload := req.Body

	defer req.Body.Close()

	err := json.NewDecoder(payload).Decode(&user)
	if err != nil {
		msg := []byte(`{
			"success": false,
			"msg": "Error parsing provided data"
		}`)

		utils.ReturnJsonResponse(res, http.StatusBadRequest, msg)
		return
	}

	user.Email = strings.ToLower(user.Email)

	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if re.MatchString(user.Email) != true {
		msg := []byte(`{
			"success": false,
			"msg": "Invalid email"
		}`)

		utils.ReturnJsonResponse(res, http.StatusBadRequest, msg)
		return
	}

	if err := a.DB.Where("email = ?", user.Email).First(&user).Error; err != gorm.ErrRecordNotFound {
		msg := []byte(`{
			"success": false,
			"msg": "Email in already in use"
		}`)

		utils.ReturnJsonResponse(res, http.StatusBadRequest, msg)
		return
	}

	user.ID = utils.GenID(a.DB, models.User{})

	valid := utils.ValidatePwd(user.Password)
	if valid != true {
		msg := []byte(`{
			"success": false,
			"msg": "Invalid password"
		}`)

		utils.ReturnJsonResponse(res, http.StatusBadRequest, msg)
		return
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		msg := []byte(`{
			"success": false,
			"msg": "Error while encrypting password"
		}`)

		utils.ReturnJsonResponse(res, http.StatusInternalServerError, msg)
		return
	}

	user.Password = string(hashedPwd[:])

	a.DB.Create(&user)

	msg := []byte(`{
		"success": true,
		"msg": "User created"
	}`)

	utils.ReturnJsonResponse(res, http.StatusCreated, msg)
	return
}

//=========================================================================
// POST: /login
//=========================================================================

func (a *APIEnv) Login(res http.ResponseWriter, req *http.Request) {
	utils.MethodValidation(res, req, "POST")

	var user models.User

	payload := req.Body

	defer req.Body.Close()

	err := json.NewDecoder(payload).Decode(&user)
	if err != nil {
		msg := []byte(`{
			"success": false,
			"msg": "Error parsing provided data"
		}`)

		utils.ReturnJsonResponse(res, http.StatusBadRequest, msg)
		return
	}

	var foundUser models.User

	err = a.DB.Where("email = ?", user.Email).First(&foundUser).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		msg := []byte(`{
			"success": false,
			"msg": "Error while querying database"
		}`)

		utils.ReturnJsonResponse(res, http.StatusInternalServerError, msg)
		return
	}

	if err == gorm.ErrRecordNotFound {
		msg := []byte(`{
			"success": false,
			"msg": "Account not found"
		}`)

		utils.ReturnJsonResponse(res, http.StatusBadRequest, msg)
		return
	}

	pwdCheck := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password))
	if pwdCheck != nil {
		msg := []byte(`{
			"success": false,
			"msg": "Invalid password"
		}`)

		utils.ReturnJsonResponse(res, http.StatusBadRequest, msg)
		return
	}

	accessTokenMaxAgeInt, err := strconv.ParseInt(accessTokenMaxAge, 10, 32)
	refreshTokenMaxAgeInt, err := strconv.ParseInt(refreshTokenMaxAge, 10, 32)

	_, err = utils.GenToken(int(accessTokenMaxAgeInt), user.ID, accessTokenPrivateKey)
	if err != nil {
		msg := []byte(`{
			"success": false,
			"msg": "Error while logging user"
		}`)

		utils.ReturnJsonResponse(res, http.StatusInternalServerError, msg)
		return
	}

	res.Header().Set("access_token", "%v")

	_, err = utils.GenToken(int(refreshTokenMaxAgeInt), user.ID, refreshTokenPrivateKey)
	if err != nil {
		msg := []byte(`{
			"success": false,
			"msg": "Error while logging user"
		}`)

		utils.ReturnJsonResponse(res, http.StatusInternalServerError, msg)
		return
	}

	res.Header().Set("refresh_token", "%v")

}

//=========================================================================
// POST: /add
//=========================================================================

func (a *APIEnv) AddIncome(res http.ResponseWriter, req *http.Request) {
	utils.MethodValidation(res, req, "POST")

	var income models.Income

	payload := req.Body

	defer req.Body.Close()

	err := json.NewDecoder(payload).Decode(&income)
	if err != nil {
		msg := []byte(`{
			"success": false,
			"msg": "Error parsing provided data"
		}`)

		utils.ReturnJsonResponse(res, http.StatusBadRequest, msg)
		return
	}

	income.ID = utils.GenID(a.DB, models.Income{})

	// incomeList := []models.Income{income}
	log.Println(income)

	msg := []byte(`{
		"success": true,
		"msg": "Income added"
	}`)

	utils.ReturnJsonResponse(res, http.StatusCreated, msg)
	return
}

var tpl = template.Must(template.ParseFiles("index.html"))

func IndexHandler(res http.ResponseWriter, req *http.Request) {
	tpl.Execute(res, nil)
}

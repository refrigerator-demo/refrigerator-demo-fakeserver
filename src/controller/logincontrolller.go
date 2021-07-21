package controller

import (
	"encoding/json"
	"fridge/src/auth"
	"fridge/src/model"
	"fridge/src/response"
	"io/ioutil"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (server *RestServer) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if nil != err {
		response.MakeJsonError(w, http.StatusUnprocessableEntity, err)
		return
	}

	user := model.User{}
	err = json.Unmarshal(body, &user)

	if nil != err {
		response.MakeJsonError(w, http.StatusUnprocessableEntity, err)
		return
	}

	token, err := server.SignIn(user.Email, user.Password)
	if nil != err {
		response.MakeJsonError(w, http.StatusUnprocessableEntity, err)
		return
	}

	response.MakeJson(w, http.StatusOK, token)
}

func (server *RestServer) SignIn(email, password string) (string, error) {
	var err error

	user := model.User{}

	err = server.DB.Debug().Model(model.User{}).Where("email = ?", email).Take(&user).Error // 나중에 user model 로직으로 빼야함
	if nil != err {
		return "", err
	}

	err = model.VerifyPassword(user.Password, password)
	if nil != err && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}

	return auth.CreateToken(user.ID)
}

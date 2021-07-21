package controller

import (
	"encoding/json"
	"errors"
//	"fmt"
	"fridge/src/auth"
	"fridge/src/model"
	"fridge/src/response"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (server *RestServer) CreateUser(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if nil != err {
		response.MakeJsonError(w, http.StatusUnprocessableEntity, err)
	}

	user := model.User{}
	err = json.Unmarshal(body, &user)

	if nil != err {
		response.MakeJsonError(w, http.StatusUnprocessableEntity, err)
		return
	}

	user.Initialize()
	//err = user.Validate("")

	if nil != err {
		response.MakeJsonError(w, http.StatusUnprocessableEntity, err)
		return
	}

	if err := user.DoBeforeSave(); nil != err {
		response.MakeJsonError(w, http.StatusUnprocessableEntity, err)
		
		return
	}

	userCreated, err := user.SaveUser(server.DB)

	if nil != err {

		//formattedError := formaterror.FormatError(err.Error())

		response.MakeJsonError(w, http.StatusInternalServerError, errors.New("formattedError"))
		return
	}

	//w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, userCreated.ID))
	userCreated.Password = ""
	response.MakeJson(w, http.StatusCreated, userCreated)
}

func (server *RestServer) GetUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		response.MakeJsonError(w, http.StatusBadRequest, err)
		return
	}
	user := model.User{}
	userGotten, err := user.FindUserByID(server.DB, uint64(uid))
	if err != nil {
		response.MakeJsonError(w, http.StatusBadRequest, err)
		return
	}
	response.MakeJson(w, http.StatusOK, userGotten)
}

func (server *RestServer) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	uid, err := strconv.ParseUint(vars["id"], 10, 64)
	if nil != err {
		response.MakeJsonError(w, http.StatusBadRequest, err)
		return
	}

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

	tokenID, err := auth.ExtractTokenID(r)

	if nil != err {
		response.MakeJsonError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	if tokenID != uint64(uid) {
		response.MakeJsonError(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	user.DoBeforeSave()

	updatedUser, err := user.UpdateUser(server.DB, uint64(uid))
	if err != nil {
		//formattedError := formaterror.FormatError(err.Error())
		response.MakeJsonError(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	response.MakeJson(w, http.StatusOK, updatedUser)
}

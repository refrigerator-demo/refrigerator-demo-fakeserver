package controller

import (
	"fridge/src/auth"
	"fridge/src/model"
	"fridge/src/response"
	"net/http"
	"errors"
	"io/ioutil"
	"encoding/json"
	"reflect"
	//"strings"
	"fmt"
)

func (server *RestServer) AllowOrigin(w http.ResponseWriter, r *http.Request) {
	
	keys := reflect.ValueOf(r.Header).MapKeys()
	strkeys := make([]string, len(keys))
	for i := 0; i < len(keys); i++ {
		strkeys[i] = keys[i].String()
	}
	//headres := strings.Join(strkeys, ",") + ", X-Naver-Client-Id,X-Naver-Client-Secret,X-TARGET-URL"

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	//w.Header().Set("Access-Control-Allow-Headers", headres)
	w.Header().Set("Access-Control-Max-Age", "3600");
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With");
   

	fmt.Printf("hah")
}


func (server *RestServer) CreateInventory(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if nil != err {
		response.MakeJsonError(w, http.StatusUnprocessableEntity, err)
		return
	}

	inventory := model.Inventory{}
	inventory.Initialize()

	err = json.Unmarshal(body, &inventory)
	if nil != err {
		response.MakeJsonError(w, http.StatusUnprocessableEntity, err)
		return
	}

	uid, err := auth.ExtractTokenID(r)
	if nil != err {
		response.MakeJsonError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		
		return
	}

	inventory.DoBeforeSave(uid)

	err = inventory.CreateInventory(server.DB)
	if nil != err {
		response.MakeJsonError(w, http.StatusUnprocessableEntity, err)
		return
	}

	
	response.MakeJson(w, http.StatusOK, inventory)
}

func (server *RestServer) GetInventoryByID(w http.ResponseWriter, r *http.Request) {
	
	uid, err := auth.ExtractTokenID(r)
	if nil != err {
		response.MakeJsonError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		
		return
	}

	inventory := model.Inventory{}

	inventories, err := inventory.FindInventoryById(server.DB, uid)
	if nil != err {
		response.MakeJsonError(w, http.StatusUnprocessableEntity, err)
		
		return
	}

	response.MakeJson(w, http.StatusOK, inventories)
}
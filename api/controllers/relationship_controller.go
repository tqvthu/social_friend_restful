package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tqvthu/social_friend_restful/api/models"
	_ "github.com/tqvthu/social_friend_restful/api/models"
	"github.com/tqvthu/social_friend_restful/api/responses"
	"github.com/tqvthu/social_friend_restful/api/utils/formaterror"
	//"github.com/tqvthu/social_friend_restful/api/utils/response"
	"io/ioutil"
	"net/http"
)
func (server *Server) Block(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	params := models.UpdateParams{}
	err = json.Unmarshal(body, &params)

	_, err = models.Block(server.DB, models.FriendService{models.UtilService{}}, params)

	if err != nil {

		formattedError := formaterror.FormatError(err.Error())

		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusBadRequest, models.MakeFriendResponse{Success: "true", Message: ""})
}

func (server *Server) SubscribeFriend(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	subscribeParams := models.UpdateParams{}
	err = json.Unmarshal(body, &subscribeParams)
	_, err = models.Subscribe(server.DB, models.FriendService{models.UtilService{}}, subscribeParams)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusCreated, models.MakeFriendResponse{Success: "true"})

}
func (server *Server) GetFriendList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars["email"]
	friendSlice, err := models.FindAllFriendsByEmail(server.DB, email)
	if (err != nil) {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusBadRequest, formattedError)
		return
	}
	listFriendRes := models.ListFriendResponse{}
	if len(friendSlice) == 0 {
		listFriendRes.Success = "true"
		listFriendRes.Message = "Friend Not Found"
		responses.JSON(w, http.StatusOK, listFriendRes)
	} else {
		listFriendRes.Success = "true"
		listFriendRes.Friend = friendSlice
		listFriendRes.Count = len(friendSlice)
		responses.JSON(w, http.StatusOK, listFriendRes)
	}
}

func (server *Server) GetRecipients(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	params := models.GetRecipientParams{}
	err = json.Unmarshal(body, &params)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	recipients, err := models.FindAllRecipients(server.DB, params)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusBadRequest, formattedError)
		return
	}
	//build response
	listFriendRes := models.ListRecipientResponse{}
	listFriendRes.Success = "true"
	listFriendRes.Recipient = recipients

	responses.JSON(w, http.StatusOK, listFriendRes)

}

func (server *Server) GetCommonFriends(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	type params struct {
		Friend []string  `json:"friends"`
	}
	friendsParams := params{}
	err = json.Unmarshal(body, &friendsParams)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	var commonFriendsSlice []string
	commonFriendsSlice = models.FindCommonFriends(server.DB, friendsParams.Friend)

	//build response
	listFriendRes := models.ListFriendResponse{}
	listFriendRes.Success = "true"
	listFriendRes.Friend = commonFriendsSlice
	listFriendRes.Count = len(commonFriendsSlice)
	responses.JSON(w, http.StatusOK, listFriendRes)

}

func (server *Server) MakeFriend(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	friendInfo := models.MakeFriendParams{}
	err = json.Unmarshal(body, &friendInfo)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	relationship := models.Relationship{}
	//relationship.Prepare()
	data, err := relationship.MakeFriend(server.DB, models.FriendService{models.UtilService{}}, friendInfo)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusBadRequest, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, data.ID))

	if data.ID != 0 {
		dataResponse := models.MakeFriendResponse{}
		dataResponse.Success = "true"
		dataResponse.Message = ""
		responses.JSON(w, http.StatusCreated, dataResponse)
	} else {
		dataResponse := models.MakeFriendResponse{}
		dataResponse.Success = "false"
		dataResponse.Message = "Bad Request"
		responses.JSON(w, http.StatusCreated, dataResponse)
	}
}


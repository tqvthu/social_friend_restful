package response

import (
	"fmt"
	"github.com/tqvthu/social_friend_restful/api/models"
)
type MakeFriendResponse struct{
	Success  string    `json:"success"`
	Message string `json:"message"`
}

type ListFriendResponse struct{
	Success  string    `json:"success" default:"true"`
	Friend []string `json:"friends"`
	Message string `json:"message"`
	Count int `json:"count"`
}

func FriendListResponse (currentUser models.User) ListFriendResponse  {
	data := ListFriendResponse{}
	if currentUser.ID == 0 {
		data.Success = "false"
		data.Message = "Bad Request"
		data.Count = 0
	}
	return data
}
func FormatResponse(item *models.Relationship, err error) MakeFriendResponse {
	fmt.Println(item.ID)
	data := MakeFriendResponse{}
	if item.ID != 0 {
		data.Success = "true"
		data.Message = ""
	} else {
		data.Success = "false"
		if err != nil {
			data.Message = err.Error()
		} else {
			data.Message = "No error found"
		}
	}
	return data
}

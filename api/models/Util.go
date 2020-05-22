package models
type MakeFriendResponse struct{
	Success  string    `json:"success"`
	Message string `json:"message"`
}
type ListRecipientResponse struct{
	Success  string    `json:"success" default:"true"`
	Recipient []string `json:"recipients"`
	Message string `json:"message"`
}
type ListFriendResponse struct{
	Success  string    `json:"success" default:"true"`
	Friend []string `json:"friends"`
	Message string `json:"message"`
	Count int `json:"count"`
}

type MakeFriendParams =  struct {
	Friend []string  `json:"friends"`
}
type UpdateParams =  struct {
	Requestor string  `json:"requestor"`
	Target string `json:"target"`
}
type FriendInfo = struct {
	RelationshipID uint64 `json:"ID"`
	Status uint32 `json:"status"`
	ActionEmail string `json:"action_email"`
	BlockedEmail string `json:"block_email"`
}
type GetRecipientParams = struct {
	Sender string `json:"sender"`
	Text string `json:"text"`
}
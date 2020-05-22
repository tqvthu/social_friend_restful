package models

import (
	"errors"
	_ "fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm"
	"github.com/tqvthu/social_friend_restful/api/utils"
	"html"
	"regexp"
	"strings"
	"time"
)

type FriendService struct {
	FriendDBAction UtilAction
}
type UtilService struct {}
type UtilAction interface {
	isBadRequest(db *gorm.DB, email []string) bool
	CheckRelationship(db *gorm.DB, email[]string) (FriendInfo, error)
}

type Relationship struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	UserOneEmail   string    `gorm:"size:255;not null;" json:"user_one_email"`
	UserTwoEmail   string    `gorm:"size:255;not null;" json:"user_two_email"`
	Status  uint32    `gorm:"not null" json:"status"`
	ActionUserID  uint32    `gorm:"not null" json:"action_user_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (r *Relationship) Prepare() {
	r.ID = 0
	r.UserOneEmail = html.EscapeString(strings.TrimSpace(r.UserOneEmail))
	r.UserTwoEmail = html.EscapeString(strings.TrimSpace(r.UserTwoEmail))
	r.Status = r.Status
	r.ActionUserID = r.ActionUserID
	r.CreatedAt = time.Now()
	r.UpdatedAt = time.Now()
}

func (r *Relationship) Validate() error {

	if r.UserOneEmail == "" {
		return errors.New("Required User One Email")
	}
	if r.UserTwoEmail == "" {
		return errors.New("Required User Two Email")
	}
	if r.ActionUserID < 1 {
		return errors.New("Required Action User ID")
	}
	return nil
}

func Block(db *gorm.DB, fs FriendService ,params UpdateParams) (Relationship, error) {
	emails := []string{params.Requestor, params.Target}
	var err error
	isBadRequest := fs.FriendDBAction.isBadRequest(db, emails)
	if isBadRequest {
		err = errors.New("Bad Request")
		return Relationship{}, err
	}
	connectionInfo, _ := fs.FriendDBAction.CheckRelationship(db, emails)
	if connectionInfo.Status == utils.BlOCK {
		return Relationship{}, errors.New("Has Blocked Already")
	}
	relationship := Relationship{}
	if connectionInfo.Status != utils.BlOCK && connectionInfo.Status != utils.PENDING { // update current relationship record

		db.Where("id = ?", connectionInfo.RelationshipID).First(&relationship)
		requestor :=  User{}
		db.Where("email = ?", params.Requestor).First(&requestor)

		relationship.ActionUserID = requestor.ID
		relationship.UserOneEmail = params.Requestor
		relationship.UserTwoEmail = params.Target
		relationship.Status = utils.BlOCK // block
		err := db.Debug().Save(&relationship).Error
		if err != nil {
			return relationship, err
		}
	} else { // create new record relationship and set status to block it means 3
		// find the requestor to set action user id
		user := User{}
		db.Where("email = ?", params.Requestor).First(&user)

		relationship.UserOneEmail = params.Requestor
		relationship.UserTwoEmail = params.Target
		relationship.ActionUserID = user.ID
		relationship.Status = utils.BlOCK // block
		err := db.Debug().Create(&relationship).Error
		if err != nil {
			return relationship, err
		}
	}
	return relationship, nil
}
/**
check if any bad request
 */
func (us UtilService) isBadRequest(db *gorm.DB, emails []string) bool{
	var count int
	db.Debug().Model(&User{}).Where("email =?", emails[0]).Count(&count)
	if count == 0 {
		return true
	}
	db.Debug().Model(&User{}).Where("email =?", emails[1]).Count(&count)
	if count == 0 {
		return true
	}
	return false
}

/**
checking relationship before processing friend connection
 */
func (us UtilService) CheckRelationship(db *gorm.DB, emails []string) (FriendInfo, error) {
	var err error
	relationship := Relationship{}
	db.Where("(user_one_email =? and user_two_email =?) or (user_one_email =? and user_two_email=?)", emails[0], emails[1], emails[1], emails[0]).First(&relationship)
	userActionID := relationship.ActionUserID
	actionUser := User{}
	db.Where("id = ?", userActionID).First(&actionUser)
	actionEmail := actionUser.Email
	var blockEmail string
	if (emails[0] == actionEmail) {
		blockEmail = emails[1]
	} else {
		blockEmail = emails[0]
	}
	return FriendInfo{
		RelationshipID: relationship.ID,
		Status: relationship.Status,
		ActionEmail: actionUser.Email,
		BlockedEmail: blockEmail,
	}, err


}
/**
make friend connection
 */
func (re *Relationship) MakeFriend(db *gorm.DB, fs FriendService, friendInfo MakeFriendParams) (Relationship, error) {
	isBadRequest := fs.FriendDBAction.isBadRequest(db, friendInfo.Friend)
	var err error
	if isBadRequest {
		return Relationship{}, errors.New("Bad Request")
	}
	r := Relationship{}
	connectionInfo, _ := fs.FriendDBAction.CheckRelationship(db,friendInfo.Friend)
	if connectionInfo.Status == utils.BlOCK {
		return r, errors.New("The user with email " + connectionInfo.ActionEmail + " Has Already Blocked the user " + connectionInfo.BlockedEmail)
	} else if connectionInfo.Status != utils.PENDING && connectionInfo.Status != utils.SUBSCRIBE {
		return r, errors.New("Has been connected already")
	}

	user1 := User{}
	user2 := User{}
	db.Where("email = ?", friendInfo.Friend[0]).First(&user1)
	db.Where("email = ?", friendInfo.Friend[1]).First(&user2)
	if user1.ID != 0 && user2.ID != 0 {
		r.UserOneEmail = friendInfo.Friend[0]
		r.UserTwoEmail = friendInfo.Friend[1]
		r.ActionUserID = user1.ID
		r.Status = utils.ACCEPTED
		err = db.Debug().Create(&r).Error
		if err != nil {
			return r, err
		}
	}

	return r, nil
}

/**
find all friends by email
 */
func FindAllFriendsByEmail(db *gorm.DB, email string) ([]string, error) {
	user := User{}
	db.Where("email =?", email).First(&user)
	var friendSlice []string
	if user.ID == 0 {
		return friendSlice, errors.New("Bad Request")
	}

	friends := []Relationship{}
	db.Where("user_one_email =? or user_two_email =? and status in (?)", email, email, []uint32{utils.ACCEPTED,utils.SUBSCRIBE}).Find(&friends)

	for _, elem := range friends {
		if !utils.Contains(friendSlice, elem.UserOneEmail) && elem.UserOneEmail != email {
			friendSlice = append(friendSlice, elem.UserOneEmail)
		}
		if !utils.Contains(friendSlice, elem.UserTwoEmail) && elem.UserTwoEmail != email {
			friendSlice = append(friendSlice, elem.UserTwoEmail)
		}
	}
	return friendSlice, nil
}

/**
	Find All Recipeints
 */
func FindAllRecipients(db *gorm.DB, params GetRecipientParams) ([]string, error) {
	var count int
	var err error
	var str []string
	db.Debug().Model(&User{}).Where("email =?", params.Sender).Count(&count)
	if count == 0 {
		return str , errors.New("Bad Request")
	}
	friends := []Relationship{}

	tx := db.Where("(user_one_email =? or user_two_email =?)", params.Sender, params.Sender )

	r, _ := regexp.Compile("([a-zA-Z0-9+._-]+@[a-zA-Z0-9._-]+\\.[a-zA-Z0-9_-]+)")
	mentionEmail := r.FindString(params.Text)
	tx = tx.Where("status in (?)", []uint32{1,4})
	tx.Find(&friends)
	var friendSlice []string
	if mentionEmail != "" {
		friendSlice = append(friendSlice, mentionEmail)
	}
	for _, elem := range friends {
		if !utils.Contains(friendSlice, elem.UserOneEmail) && elem.UserOneEmail != params.Sender {
			friendSlice = append(friendSlice, elem.UserOneEmail)
		}
		if !utils.Contains(friendSlice, elem.UserTwoEmail) && elem.UserTwoEmail != params.Sender {
			friendSlice = append(friendSlice, elem.UserTwoEmail)
		}
	}
	return friendSlice, err
}

/**
subscribe
 */
func Subscribe(db *gorm.DB, fs FriendService, params UpdateParams) (Relationship, error) {
	emails := []string{params.Requestor, params.Target}
	isBadRequest := fs.FriendDBAction.isBadRequest(db, emails)
	var err error
	r := Relationship{}
	if isBadRequest {
		return r, errors.New("Bad Request")
	}

	connectionInfo, _ := fs.FriendDBAction.CheckRelationship(db, emails)

	if connectionInfo.RelationshipID != 0 { // has connected already
		if  connectionInfo.Status == utils.BlOCK && connectionInfo.ActionEmail == params.Target {
			return r, errors.New(params.Requestor + "has been blocked already")
		} else{
			return r, errors.New("Has connected already")
		}
	}

	requestor := User{}
	db.Where("email = ?", params.Requestor).First(&requestor)

	r.UserOneEmail = params.Requestor
	r.UserTwoEmail = params.Target
	r.Status = utils.SUBSCRIBE // subscribe
	r.ActionUserID = requestor.ID

	err = db.Debug().Create(&r).Error
	if err != nil {
		return r, err
	}
	return r, nil
}

/**
find common friends
 */
func FindCommonFriends(db *gorm.DB, email []string) []string {
	friends := [] Relationship{}
	db.Where("user_one_email in(?) or user_two_email in(?)", email, email).Find(&friends)
	var commonFriendsSlice []string
	for _, elem := range friends {
		if !utils.Contains(commonFriendsSlice, elem.UserOneEmail) && elem.UserOneEmail != email[0] && elem.UserOneEmail != email[1] {
			commonFriendsSlice = append(commonFriendsSlice, elem.UserOneEmail)
		}
		if !utils.Contains(commonFriendsSlice, elem.UserTwoEmail) && elem.UserTwoEmail != email[1] && elem.UserOneEmail != email[0]{
			commonFriendsSlice = append(commonFriendsSlice, elem.UserTwoEmail)
		}
	}

	return commonFriendsSlice
}

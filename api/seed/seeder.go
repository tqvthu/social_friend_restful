package seed

import (
	"github.com/jinzhu/gorm"
	"github.com/tqvthu/social_friend_restful/api/models"
	"log"
)

var users = []models.User{
	models.User{
		Nickname: "Steven victor",
		Email:    "steven@gmail.com",
		Password: "password",
	},
	models.User{
		Nickname: "Martin Luther",
		Email:    "luther@gmail.com",
		Password: "password",
	},
	models.User{
		Nickname: "Jennifer Lopez",
		Email:    "jenifer@gmail.com",
		Password: "password",
	},
	models.User{
		Nickname: "Tome Cruise",
		Email:    "tom@gmail.com",
		Password: "password",
	},
	models.User{
		Nickname: "David Beckham",
		Email:    "david@gmail.com",
		Password: "password",
	},
	models.User{
		Nickname: "Messi",
		Email:    "messi@gmail.com",
		Password: "password",
	},
}

func Load(db *gorm.DB) {
	var count int
	db.Debug().Model(&models.User{}).Count(&count)
	if count != 0 {
		return
	}
	err := db.Debug().DropTableIfExists(&models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.Relationship{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.User{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}
	err = db.Debug().Model(&models.Relationship{}).AddForeignKey("action_user_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	for i, _ := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
	}
}
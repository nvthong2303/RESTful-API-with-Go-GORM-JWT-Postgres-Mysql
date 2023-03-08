package seed

import (
	"github.com/jinzhu/gorm"
	"github.com/nvthong2303/echo_crud_01/api/models"
	"log"
)

var users = []models.User{
	models.User{
		Nickname: "Nguyen van thong",
		Email:    "nvthong2303@gmail.com",
		Password: "123456",
	},
	models.User{
		Nickname: "Nguyen van thong test",
		Email:    "nvthong2704@gmail.com",
		Password: "123456",
	},
}

var posts = []models.Post{
	models.Post{
		Title:   "Title 1",
		Content: "Content 1",
	},
	models.Post{
		Title:   "Title 2",
		Content: "Content 2",
	},
}

func Load(db *gorm.DB) {
	err := db.Debug().DropTableIfExists(&models.Post{}, &models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.User{}, &models.Post{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	err = db.Debug().Model(&models.Post{}).AddForeignKey("author_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error : %v", err)
	}

	for i, _ := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed user table: %v", err)
		}
		posts[i].AuthorID = users[i].ID

		err = db.Debug().Model(&models.Post{}).Create(&posts[i]).Error
		if err != nil {
			log.Fatalf("cannot seed post table: %v", err)
		}
	}
}

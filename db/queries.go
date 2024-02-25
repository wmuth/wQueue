package db

import (
	"errors"
	"fmt"
	"os"

	"wQueue/encrypt"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupConn() (*gorm.DB, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, errors.New("error loading .env")
	}

	envVars := [...]string{"DB_HOST", "DB_USER", "DB_PASS", "DB_NAME", "DB_PORT", "DB_ENCR"}
	for _, s := range envVars {
		_, exists := os.LookupEnv(s)
		if !exists {
			return nil, fmt.Errorf("error loading environment variable %s", s)
		}
	}

	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_ENCR"))

	return gorm.Open(postgres.Open(connStr), &gorm.Config{})
}

func Login(c *gorm.DB, uName, pass string) (bool, User, error) {
	var u User
	if err := c.Where("username = $1", uName).First(&u).Error; err != nil {
		return false, u, err
	}

	corr, err := encrypt.CheckPass(u.Password, pass)
	if err != nil {
		return false, u, err
	}

	if corr {
		return true, u, nil
	}
	return false, u, nil
}

func GetAllQueues(c *gorm.DB) ([]Queue, error) {
	var q []Queue
	if err := c.Find(&q).Error; err != nil {
		return q, err
	}
	return q, nil
}

func GetQueueByTitle(c *gorm.DB, title string) (Queue, error) {
	var q Queue
	if err := c.Where("title = $1", title).First(&q).Error; err != nil {
		return q, err
	}
	return q, nil
}

func GetUserByID(c *gorm.DB, id int) (User, error) {
	var u User
	if err := c.Where("id = $1", id).First(&u).Error; err != nil {
		return u, err
	}
	return u, nil
}

func UserIsAdminOfQueue(c *gorm.DB, u *User, q *Queue) error {
	var a Admin
	if err := c.Where("userid = $1 AND queueid = $2", u.Id, q.Id).Limit(1).Find(&a).Error; err != nil {
		return err
	}
	if a.Userid == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func AdminCloseQueue(c *gorm.DB, q *Queue) error {
	q.Open = false
	if err := c.Save(q).Error; err != nil {
		return err
	}
	return nil
}

func AdminOpenQueue(c *gorm.DB, q *Queue) error {
	q.Open = true
	if err := c.Save(q).Error; err != nil {
		return err
	}
	return nil
}

func AddToQueue(c *gorm.DB, u *User, q *Queue, loc, com string) error {
	qi := Queueitem{
		Userid:   u.Id,
		Queueid:  q.Id,
		Location: loc,
		Active:   true,
		Comment:  com,
	}

	if err := c.Create(&qi).Error; err != nil {
		return err
	}
	return nil
}

func GetUserAccesses(c *gorm.DB, u *User) ([]Access, error) {
	var a []Access
	if err := c.Where("userid = $1", u.Id).Find(&a).Error; err != nil {
		return a, err
	}
	return a, nil
}

func GetQueueitem(c *gorm.DB, u *User, q *Queue) (Queueitem, error) {
	var qi Queueitem
	if err := c.Where("userid = $1 AND queueid = $2", u.Id, q.Id).First(&qi).Error; err != nil {
		return qi, err
	}
	return qi, nil
}

func GetQueueitemsInQueue(c *gorm.DB, q *Queue) ([]Queueitem, error) {
	var qs []Queueitem
	if err := c.Where("queueid = $1", q.Id).Find(&qs).Error; err != nil {
		return qs, err
	}
	return qs, nil
}

func CountQueueitemsInQueue(c *gorm.DB, q *Queue) (int64, error) {
	var count int64
	if err := c.Model(&Queueitem{}).Where("queueid = $1", q.Id).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func RemoveFromQueue(c *gorm.DB, qi *Queueitem) error {
	if err := c.Delete(&qi).Error; err != nil {
		return err
	}
	return nil
}

func SetQueueMessage(c *gorm.DB, q *Queue, m string) error {
	q.Message = m
	if err := c.Updates(&q).Error; err != nil {
		return err
	}
	return nil
}

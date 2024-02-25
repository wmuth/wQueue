package db

import (
	"time"
)

type Queue struct {
	Id      int   `gorm:"primaryKey"`
	InQueue int64 `gorm:"-"`
	Message string
	Open    bool
	Title   string
}

type User struct {
	Id       int `gorm:"primaryKey"`
	Password string
	Username string
}

type Message struct {
	Id     int `gorm:"primaryKey"`
	Fromid int
	Text   string
	Time   time.Time
	Toid   int
}

type Queueitem struct {
	Id       int `gorm:"primaryKey"`
	Active   bool
	Comment  string
	Location string
	Queueid  int
	Userid   int
}

type Admin struct {
	Queueid int `gorm:"primaryKey"`
	Userid  int `gorm:"primaryKey"`
}

type Access struct {
	Queueid int `gorm:"primaryKey"`
	Userid  int `gorm:"primaryKey"`
}

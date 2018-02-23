package main

import (
	"fmt"
	"time"

	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"github.com/golang/glog"
	_ "github.com/lib/pq"
)

const (
	DB_HOST     = "localhost"
	DB_PORT     = 5433
	DB_USER     = "postgres"
	DB_PASSWORD = "root"
	DB_NAME     = "students"
)

type User_info struct {
	Uid       int64 `xorm:"pk not null autoincr"`
	Name      string
	Dept      string
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
}

type Library_info struct {
	Entry_id   int64 `xorm:"pk not null autoincr"`
	Book_id    int64
	Student_id int64
}

type Book_info struct {
	Book_id     int64 `xorm:"pk not null autoincr"`
	Author_name string
	Version     int `xorm:"version"`
}

//Structure required for joining tables

type User_library struct {
	User_info    `xorm:"extends"`
	Library_info `xorm:"extends"`
	Book_info    `xorm:"extends"`
}

func (User_library) TableName() string {
	return "user_info"
}

func connect_database() *xorm.Engine {
	var err error

	dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)

	en, err := xorm.NewEngine("postgres", dbinfo)
	if err != nil {
		glog.Fatalln("engine creation failed", err)
	}

	err = en.Ping()
	if err != nil {
		glog.Fatalln("Connect to database failed", err)
	} else {
		en.SetMapper(core.SameMapper{})
	}

	glog.Infoln("successfully connected to postgres database.")
	return en
}

func sync_tables(en *xorm.Engine) {
	err := en.Sync(new(User_info))
	err = en.Sync(new(Library_info))
	err = en.Sync(new(Book_info))

	if err != nil {
		glog.Fatalln("creation error", err)
		return
	}

	glog.Infoln("Successfully synced")
}

func insert_data(en *xorm.Engine) {
	u := new(User_info)
	u.Uid = 1104029
	u.Name = "Maruf"
	u.Dept = "CSE"
	affected, _ := en.Insert(u)
	//affected, err := en.Insert(&User_info{Uid: 11040501, Name: "Diptal", Dept: "ME"})
	glog.Infoln("insertion single data. affected:", affected)
}

func insert_multiple_data(en *xorm.Engine) {
	users := make([]User_info, 3)

	users[0].Name = "Shohel"
	users[0].Dept = "Archi"
	users[0].Uid = 1104023

	users[1].Name = "Hasan"
	users[1].Dept = "CSE"
	users[1].Uid = 1104009

	users[2].Name = "Uttam"
	users[2].Dept = "ETE"
	users[2].Uid = 1104089

	affected, _ := en.Insert(&users)
	glog.Infoln("insertion multiple data. affected:", affected)
}

func main() {
	en := connect_database()
	sync_tables(en)
	insert_data(en)
	insert_multiple_data(en)
}

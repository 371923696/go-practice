package main

import (
	"fmt"
	"log"
	"strconv"
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
		glog.Fatalln("engine creation failed", err.Error())
	}

	err = en.Ping()
	if err != nil {
		glog.Fatalln("Connect to database failed", err.Error())
	} else {
		//en.SetTableMapper(core.SameMapper{})
		//en.SetColumnMapper(core.GonicMapper{})
		en.SetMapper(core.GonicMapper{})
	}

	glog.Infoln("successfully connected to postgres database.")
	return en
}

func sync_tables(en *xorm.Engine) {
	err := en.Sync(new(User_info))
	err = en.Sync(new(Library_info))
	err = en.Sync(new(Book_info))

	if err != nil {
		glog.Fatalln("creation error", err.Error())
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

func query_single_data(en *xorm.Engine) {
	user := User_info{Uid: 1104029, Dept: "CSE"}
	has, err := en.Get(&user)
	if err != nil {
		glog.Fatalln("query failed.", err.Error())
	}

	glog.Infoln(has)
	glog.Infoln(user)
}

func query_multiple_data(en *xorm.Engine) {
	var users []User_info
	err := en.Find(&users)
	if err != nil {
		glog.Fatalln("find user info failed.", err.Error())
	}

	for index, user := range users {
		glog.Infoln(index, user)
	}
}

func query_conditional_data(en *xorm.Engine) {
	var err error

	var tenusers []User_info
	//err = en.Cols("Uid", "Name").Where("Dept = ?", "CSE").Limit(10).Find(&tenusers)
	err = en.Cols("Uid", "Name").Where("Dept = ?", "CSE").Limit(10).Find(&tenusers)
	if err != nil {
		glog.Fatalln("query condition failed.", err.Error())
	}

	for index, user := range tenusers {
		glog.Infoln(index, user.Uid, user.Name)
	}

}

func query_exec_sql(en *xorm.Engine) {
	sql := "select * from user_info"
	results, err := en.QueryString(sql)
	if err != nil {
		glog.Fatalln("query exec sql failed.", err.Error())
	}

	for index, result := range results {
		glog.Infoln(index, result)
	}
}

func join_User_Library(en *xorm.Engine) {
	var users []User_library
	err := en.Join("INNER", "library_info", "library_info.student_id = user_info.uid").Join("INNER", "book_info", "book_info.book_id = library_info.book_id").Find(&users)
	if err != nil {
		glog.Fatalln(err)
	}
	//log.Println(users)

	for _, user := range users {
		glog.Infoln(user)
	}
}

func update_data(en *xorm.Engine) {
	user := new(User_info)
	user.Name = "zhen.zhang"
	affected, err := en.Id(1104029).Update(user)
	if err != nil {
		glog.Fatalln("update failed.", err.Error())
	}
	glog.Infoln("affected row: ", affected)

}

func delete_data(en *xorm.Engine) {
	affected, err := en.Id(1104029).Delete(&User_info{})
	if err != nil {
		glog.Fatalln("delete failed.", err.Error())
	}

	glog.Infoln("delete affected:", affected)
}

func sql_command(en *xorm.Engine) {
	//en.Query for select, en.Exec for insert, update or delete
	sql := "select * from user_info"
	results, _ := en.Query(sql)

	for _, result := range results {
		//convert results []map[string][]byte to users []UserInfo
		var user User_info
		user.Uid, _ = strconv.ParseInt(string(result["uid"]), 10, 64) //10 base, 64bit
		user.Name = string(result["name"])
		user.Dept = string(result["dept"])

		layout := "2006-01-02T15:04:05Z"
		user.UpdatedAt, _ = time.Parse(layout, string(result["updated_at"]))
		user.CreatedAt, _ = time.Parse(layout, string(result["created_at"]))

		log.Println(user)
	}

	//Exec for update
	sql = "update user_info set name=? where uid=?"
	res, _ := en.Exec(sql, "Mynul Hasan", 1104009)
	log.Println("affected", res)
}

func main() {
	en := connect_database()
	//sync_tables(en)
	//insert_data(en)
	//insert_multiple_data(en)
	//query_single_data(en)
	//query_multiple_data(en)
	//query_conditional_data(en)
	//query_exec_sql(en)
	join_User_Library(en)
	//update_data(en)
	//delete_data(en)
	//sql_command(en)

}

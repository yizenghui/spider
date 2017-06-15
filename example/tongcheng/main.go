package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/yizenghui/spider/tongcheng"
)

// Job 职位模型
type Job struct {
	gorm.Model
	PublishAt    int `sql:"index" default:"0"`
	Title        string
	Position     string
	Category     string
	PositionName string
	CategoryName string
	Location     string
	Salery       string
	PayType      string
	Description  string
	Number       string
	Education    string
	WorkYears    string
	FromURL      string `sql:"index"`
	Company      string
	CompanyURL   string
	Linkman      string
	Telephone    string
	Email        string
	Address      string
	Lng          string
	Lat          string
	Welfare      string
}

var db *gorm.DB

func main() {

	area := "gz"
	_, e := fmt.Scanln(&area)
	if nil != e {
		panic("初始化地区失败")
	}

	var err error
	db, err = gorm.Open("sqlite3", area+"_job.db")
	// db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=spider sslmode=disable password=123456")

	if err != nil {
		panic("连接数据库失败")
	}

	// 自动迁移模式
	db.AutoMigrate(&Job{})
	defer db.Close()

	getList(area)

}

func getList(area string) {

	pcListURL := "http://" + area + ".58.com/job/"
	ticker := time.NewTicker(time.Minute * 2)
	for _ = range ticker.C {
		fmt.Printf("ticked at %v spider %v \n", time.Now(), pcListURL)
		go spiderJobList(pcListURL)
	}
}

func spiderJobList(pcListURL string) {
	rows1, err := tongcheng.GetJobList(pcListURL)
	if err == nil {
		// fmt.Println(rows1)
		for k, v := range rows1 {
			fmt.Printf("%v -> %v  %v\n", k, v.Title, v.InfoURL)

			time.Sleep(1 * time.Second)
			syncJob(v.InfoURL)
		}
	}
}

func getInfo(infoURL string) (tongcheng.Job, error) {

	info, err := tongcheng.GetJobInfoData(infoURL)
	if err != nil {
		return info, err
	}
	return info, nil

}

// 同步职位
func syncJob(link string) {

	// info := getInfo(link)
	// fmt.Println(info)
	// 读取
	var job Job
	db.Where(Job{FromURL: link}).FirstOrCreate(&job)
	// fmt.Println(job)

	updateTime := job.UpdatedAt.Unix()
	createTime := job.CreatedAt.Unix()
	todayTime := getUpdateTime()

	// 如果更新时间小于今天或者是刚刚创建的数据，采集详细页并更新
	if createTime == updateTime || updateTime < todayTime {

		// fmt.Println(job)
		info, err := getInfo(link)
		if err == nil {

			// fmt.Println(info)
			fmt.Printf("%s -> %s  %s\n", link, info.Title, info.Company)
			job.PublishAt = 0
			job.Title = info.Title
			job.Position = info.Position
			job.PositionName = info.PositionName
			job.Category = info.Category
			job.CategoryName = info.CategoryName
			job.Location = info.Location
			job.Salery = info.Salery
			job.PayType = info.PayType
			job.Description = info.Description
			job.Number = info.Number
			job.Education = info.Education
			job.WorkYears = info.WorkYears
			job.FromURL = info.FromURL
			job.Company = info.Company
			job.CompanyURL = info.CompanyURL
			job.Linkman = info.Linkman
			job.Telephone = info.Telephone
			job.Email = info.Email
			job.Address = info.Address
			job.Lng = info.Lng
			job.Lat = info.Lat
			job.Welfare = strings.Join(info.Tags, ",")
			db.Save(&job)
		}
		// fmt.Println(job)
	}

}

// 获取更新时间界限 (如果更新时间小于界限，就去更新职位)
func getUpdateTime() int64 {

	// nTime := time.Now()
	// yesTime := nTime.AddDate(0, 0, -1).Unix()

	timeStr := time.Now().Format("2006-01-02")
	// fmt.Println("timeStr:", timeStr)
	today, _ := time.Parse("2006-01-02", timeStr)
	todayTime := today.Unix() - 8*3600
	// fmt.Println("timeNumber:", todayTime)
	return todayTime
}

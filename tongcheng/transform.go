package tongcheng

import (
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/lunny/html2md"        // html 2 markdown
	"github.com/xuebing1110/location" // location name 2 location code
	"github.com/yizenghui/gps"
	"github.com/yizenghui/spider/code"
	"github.com/yizenghui/spider/conf"
)

// Job 采集58职位数据结构
type Job struct {
	Title        string
	Position     string
	PositionName string
	Category     string
	CategoryName string
	Location     string
	Salery       string
	PayType      string
	Description  string
	Number       string
	Education    string
	WorkYears    string
	FromURL      string
	Company      string
	CompanyURL   string
	Linkman      string
	Telephone    string
	Email        string
	Address      string
	Lng          string
	Lat          string
	Tags         []string
}

// JobData 本地保存数据结构
type JobData struct {
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

// PubJob 提交转换的数据结构
type PubJob struct {
	Title       string // 职位标题
	Position    string // 原职位分类
	Company     string // 公司名
	Category    int    // 分类
	Area        int    // 地区
	MinPay      int    // 最小月薪
	MaxPay      int    // 最大月薪
	Education   int    // 学历
	Experience  int    // 工作经验
	Welfare     []int  // 标签
	Description string
	SourceFrom  string // string默认长度为255, 使用这种tag重设。
	CompanyURL  string // string默认长度为255, 使用这种tag重设。
	Linkman     string
	Telephone   string
	Email       string
	Address     string
	Lng         float64
	Lat         float64
}

// TransformJob 数据转换
func TransformJob(job JobData) PubJob {
	var pj PubJob
	pj.Title = job.Title
	pj.Position = job.Position
	pj.Company = job.Company
	pj.Category = TransformCategory(job.CategoryName + job.Category)
	pj.Area = TransformArea(job.Location)
	pj.MinPay, pj.MaxPay = TransformSalary(job.Salery)
	pj.Education = TransformEducation(job.Education)
	pj.Experience = TransformExperience(job.WorkYears)
	pj.Welfare = TransformWelfare(job.Welfare)
	pj.SourceFrom = job.FromURL
	pj.CompanyURL = job.CompanyURL
	pj.Linkman = job.Linkman
	pj.Telephone = job.Telephone
	pj.Email = job.Email
	pj.Address = job.Address
	lng, ex := strconv.ParseFloat(job.Lng, 64)
	lat, ey := strconv.ParseFloat(job.Lat, 64)
	if ey == nil && ex == nil && lng != 0 && lat != 0 {
		// 百度转火星
		latG1, lngG1 := gps.BD09ToGCJ02(lat, lng)
		// 火星转国际
		pj.Lat, pj.Lng = gps.GCJ02ToWGS84(latG1, lngG1)
	}

	pj.Description = html2md.Convert(job.Description)
	return pj
}

// TransformArea 转化城市 使用了 github.com/xuebing1110/location 的 GetAdcode 满足目前需求
func TransformArea(text string) int {
	var code int
	level := []string{"市", "区", ""}
	cities := strings.Split(text, "-")
	for l, city := range cities {
		val := location.GetAdcode(city + level[l])
		value, err := strconv.Atoi(val)
		if err == nil && value > code {
			code = value
		}
		// fmt.Println(city+level[l], code)
	}
	return code
}

// TransformEducation 转化学历要求
func TransformEducation(text string) int {
	education := code.Mate(text, EducationMate, conf.Education)
	return education
}

// TransformExperience 转化工作经验
func TransformExperience(text string) int {
	experience := code.Mate(text, ExperienceMate, conf.Experience)
	return experience
}

// TransformSalary 转化待遇
func TransformSalary(text string) (min, max int) {

	// Contains Compare

	if text == "面议" {
		return 0, 0
	}

	if text == string("1000") {
		return 0, 1000
	}

	if text == string("1000-2000") {
		return 1000, 2000
	}

	if text == string("2000-3000") {
		return 2000, 3000
	}

	if text == string("3000-5000") {
		return 3000, 5000
	}

	if text == string("5000-8000") {
		return 5000, 8000
	}

	if text == string("8000-12000") {
		return 8000, 12000
	}

	if text == string("12000-20000") {
		return 12000, 20000
	}

	if text == string("20000-25000") {
		return 20000, 25000
	}

	if text == string("25000") {
		return 25000, 0
	}

	return 0, 0
}

// TransformSalaryContains 转化待遇(包含匹配)
func TransformSalaryContains(text string) (min, max int) {

	// Contains Compare

	if b := strings.Contains(text, string("面议")); b == true {
		return 0, 0
	}

	if b := strings.Contains(text, string("1000")); b == true {
		return 0, 1000
	}

	if b := strings.Contains(text, string("1000-2000")); b == true {
		return 1000, 2000
	}
	if b := strings.Contains(text, string("2000-3000")); b == true {
		return 2000, 3000
	}
	if b := strings.Contains(text, string("3000-5000")); b == true {
		return 3000, 5000
	}
	if b := strings.Contains(text, string("5000-8000")); b == true {
		return 5000, 8000
	}
	if b := strings.Contains(text, string("8000-12000")); b == true {
		return 8000, 12000
	}
	if b := strings.Contains(text, string("12000-20000")); b == true {
		return 12000, 20000
	}
	if b := strings.Contains(text, string("20000-25000")); b == true {
		return 20000, 25000
	}
	if b := strings.Contains(text, string("25000")); b == true {
		return 25000, 0
	}

	return 0, 0
}

// TransformWelfare 转化标签
func TransformWelfare(text string) []int {
	welfares := strings.Split(text, ",")
	var tags []int
	if welfares != nil {
		for _, welfare := range welfares {
			val := code.Mate(welfare, WelfareMate, conf.Welfare)
			if val != 0 {
				tags = append(tags, val)
			}
		}
	}
	return tags
}

// TransformCategory 获取线上分类ID 目前由 job.CategoryName+job.Category 组合获取
func TransformCategory(text string) int {
	category := code.Mate(text, CategoryMate, conf.Category)
	return category
}

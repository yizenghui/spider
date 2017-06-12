package tongcheng

import (
	"strconv"
	"strings"

	"github.com/xuebing1110/location"
	"github.com/yizenghui/spider/code"
	"github.com/yizenghui/spider/conf"
)

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

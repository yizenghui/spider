package tongcheng

import (
	"bytes"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/yizenghui/spider/code"
)

// JobRows 58PC职位列表记录集
type JobRows []JobItem

// JobItem 58PC职位列参数
type JobItem struct {
	Title      string
	Company    string
	CompanyURL string
	Location   string
	Date       string
	InfoURL    string
}

// GetJobInfoData 获取职位详细数据
func GetJobInfoData(mobileInfoURL string) (SourceJob, error) {

	var job SourceJob

	job.FromURL = mobileInfoURL

	client := &http.Client{}

	reqest, err := http.NewRequest("GET", mobileInfoURL, nil)

	if err != nil {
		//	 handle error
		return job, err

	}
	reqest.Header.Add("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B143 Safari/601.1")
	//	time.Sleep(10 * time.Second)
	// fmt.Println(mobile_info_url)

	resp, err := client.Do(reqest)

	if err != nil {
		//	 handle error
		return job, err

	}
	g, err := goquery.NewDocumentFromReader(resp.Body)
	resp.Body.Close()
	//
	// g, e := goquery.NewDocument(mobile_info_url)

	if err != nil {
		// fmt.Println(e)
		return job, err
	}
	//
	// 标题
	job.Title = strings.TrimSpace(g.Find("#d_title").Text())

	// 企业
	job.Company = strings.TrimSpace(g.Find("#openCom").Find(".c_tit").Find("a").Eq(0).Text())
	// 企业
	job.CompanyURL, _ = g.Find("#openCom").Find(".c_tit").Find("a").Eq(0).Attr("href")

	// 职位
	job.Position = strings.TrimSpace(g.Find(".job_con ul li").Eq(0).Find(".attrValue").Text())

	/* 通过正则从链接地址中匹配出分类标识start */
	// http://m.58.com/gz/meirongzhuli/
	var positionNameURL, _ = g.Find(".job_con ul li").Eq(0).Find(".attrValue").Find("a").Eq(0).Attr("href")

	// fmt.Println(positionNameURL)
	// 职位分类标识
	job.PositionName = code.FindString(`http://m.58.com/(?P<location>\w+)/(?P<category>\w+)/`, positionNameURL, "category")
	/* 通过正则从链接地址中匹配出分类标识end */

	// 分类
	job.Category = strings.TrimSpace(g.Find(".job_con ul li").Eq(0).Find(".attrValue").Text())

	// 工作地点
	job.Location = strings.TrimSpace(g.Find(".job_con ul li").Eq(2).Find(".dizhiValue").Text())

	// 工作地址
	job.Address = strings.TrimSpace(g.Find(".company_con ul li").Eq(3).Find(".dizhiValue").Text())

	// 待遇
	job.Salery = strings.TrimSpace(g.Find(".pay").Eq(0).Text())

	// pay_yy
	job.PayType = strings.TrimSpace(g.Find(".pay_yy").Eq(0).Text())

	// 要求
	// demand := g.Find(".job_con ul li").Eq(1).Find(".attrValue").Text()
	// fmt.Println(strings.TrimSpace(demand))

	// dis_con
	// 要求
	job.Description, _ = g.Find(".position_dis").Eq(0).Find(".dis_con p").Html()

	job.Description = strings.TrimSpace(job.Description)

	html, _ := g.Html()
	//	fmt.Printf("html %s\n", html)

	// 联系人
	job.Linkman = code.FindString(`{"I":"5333","V":"(?P<linkman>[^"]+)"}`, html, "linkman")

	// 电话
	job.Telephone, _ = g.Find(".phoneWrap").Eq(0).Find(".phone").Attr("phoneno")
	job.Telephone = strings.TrimSpace(job.Telephone)

	//listname:'tech'},locallist

	// 分类
	categoryMap := code.SelectString(`,name:'(?P<category>[^']+)',listname:'(?P<category_name>[^']+)'},locallist`, html)

	job.Category = categoryMap["category"]
	// 分类标识
	job.CategoryName = categoryMap["category_name"]

	// 人数
	job.Number = code.FindString(`{"I":"5353","V":"(?P<number>[^"]+)"}`, html, "number")

	//学历要求
	job.Education = code.FindString(`,学历要求(?P<education>[^,]+),`, html, "education")

	//工作经验
	job.WorkYears = code.FindString(`,工作经验(?P<workYears>[^,]+),`, html, "workYears")

	//TODO 是否可接收应届生
	freshGraduate := code.FindString(`,可接受(?P<freshGraduate>[^,]+),学历,`, html, "freshGraduate")
	if freshGraduate == "" {

	}

	//企业性质
	job.CompanyType = code.FindString(`<li><span class="attrName">性质：</span><span class="attrValue">(?P<companyType>[^<]+)</span></li>`, html, "companyType")

	// 公司规模
	job.CompanySize = code.FindString(`,{"I":"5755","V":"(?P<companySize>[^"]+)"},{"I":"comp_id",`, html, "companySize")

	// 公司所属行业
	job.CompanyIndustry = strings.TrimSpace(code.FindString(`{"I":"camp_indus","V":"(?P<companyIndustry>[^"]+)"}`, html, "companyIndustry"))

	// 公司详细描述
	job.CompanyDescription = strings.TrimSpace(g.Find(".company_con ul p").Eq(0).Text())

	job.Lat = code.FindString(`{"I":"6691","V":"(?P<lat>[^,]+)"},`, html, "lat")

	job.Lng = code.FindString(`{"I":"6692","V":"(?P<lng>[^,]+)"},`, html, "lng")

	email := code.FindString(`{"I":"5360","V":"(?P<email>[^"]+)"}`, html, "lng")
	checkEmail, _ := regexp.MatchString(`(?P<first>\w+).58.com`, email)
	// 不要58的邮箱
	if !checkEmail {
		job.Email = email
	}
	// 福利标签
	g.Find(".fulivalue").Find("span").Each(func(i int, content *goquery.Selection) {
		text := content.Text()
		// fmt.Println(text)
		job.Tags = append(job.Tags, text)
	})
	return job, nil

}

// GetJobList 获取职位列表
func GetJobList(pcListURL string) (JobRows, error) {
	var rows JobRows
	var item JobItem
	g, e := goquery.NewDocument(pcListURL)
	if e != nil {
		return rows, e
	}

	// 下列内容于 2017年4月4日 20:50:24 抓取
	g.Find("#infolist dl").Each(func(i int, content *goquery.Selection) {

		// 职位详细页链接地址
		item.InfoURL, _ = content.Find(".t").Attr("href")
		// 验证URL是否我们需要的58PC job info url
		checkLinkIsJobInfo, _ := regexp.MatchString(`http://(?P<site>\w+).58.com/(?P<cate>\w+)/(?P<id>\d+)x.shtml`, item.InfoURL)
		if checkLinkIsJobInfo {

			item.InfoURL = GetMobileJobInfoLink(item.InfoURL)
			// 标题
			item.Title = strings.TrimSpace(content.Find(".t").Text())
			// 地点
			item.Location = strings.TrimSpace(content.Find(".w96").Text())
			// 企业名
			item.Company = strings.TrimSpace(content.Find(".fl").Text())
			// 企业详细页链接地址
			item.CompanyURL, _ = content.Find(".fl").Attr("href")

			// 职位更新时间
			item.Date = strings.TrimSpace(content.Find(".w68").Text())

			// fmt.Println(item)
			rows = append(rows, item)
		}
	})

	return rows, nil
}

//GetMobileJobInfoLink 通过pc详细页地址 获取手机端详细页地址
func GetMobileJobInfoLink(pcInfoLink string) string {
	linkMap := code.SelectString(`http://(?P<site>\w+).58.com/(?P<cate>\w+)/(?P<id>\d+)x.shtml`, pcInfoLink)
	var buffer bytes.Buffer
	buffer.WriteString("http://m.58.com/")
	buffer.WriteString(linkMap["site"])
	buffer.WriteString("/")
	buffer.WriteString(linkMap["cate"])
	buffer.WriteString("/")
	buffer.WriteString(linkMap["id"])
	buffer.WriteString("x.shtml")
	return buffer.String()
}

package tongcheng

import (
	"bytes"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/yizenghui/spider/code"
)

// GetJobInfoData 获取职位详细数据
func GetJobInfoData(mobileInfoURL string) (Job, error) {

	var job Job

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
	positionNameExp := myRegexp{regexp.MustCompile(`http://m.58.com/(?P<location>\w+)/(?P<category>\w+)/`)}

	positionNameExpMap := positionNameExp.FindStringSubmatchMap(positionNameURL)

	// 职位分类标识
	job.PositionName, _ = positionNameExpMap["category"]
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

	// 工作经验经验不限

	numberExp := myRegexp{regexp.MustCompile(`{"I":"5353","V":"(?P<number>[^"]+)"}`)}

	numberMap := numberExp.FindStringSubmatchMap(html)
	job.Number = numberMap["number"]

	//,学历要求学历不限,
	educationExp := myRegexp{regexp.MustCompile(`,学历要求(?P<education>[^,]+),`)}

	educationMap := educationExp.FindStringSubmatchMap(html)
	job.Education = educationMap["education"]

	//,工作经验1-2年,

	workYearsExp := myRegexp{regexp.MustCompile(`,工作经验(?P<workYears>[^,]+),`)}

	workYearsMap := workYearsExp.FindStringSubmatchMap(html)
	job.WorkYears = workYearsMap["workYears"]

	latExp := myRegexp{regexp.MustCompile(`{"I":"6691","V":"(?P<lat>[^,]+)"},`)}
	latMap := latExp.FindStringSubmatchMap(html)
	job.Lat = latMap["lat"]

	lngExp := myRegexp{regexp.MustCompile(`{"I":"6692","V":"(?P<lng>[^,]+)"},`)}
	lngMap := lngExp.FindStringSubmatchMap(html)
	job.Lng = lngMap["lng"]

	emailExp := myRegexp{regexp.MustCompile(`{"I":"5360","V":"(?P<email>[^"]+)"}`)}

	mmap := emailExp.FindStringSubmatchMap(html)
	email := mmap["email"]
	// fmt.Println(mmap)
	// fmt.Println(email)

	checkEmail, _ := regexp.MatchString(`(?P<first>\w+).58.com`, email)

	// 不要58的邮箱
	if !checkEmail {
		job.Email = email
	}

	// fmt.Println(myExp.FindStringSubmatch(html))

	//

	// 福利标签
	g.Find(".fulivalue").Find("span").Each(func(i int, content *goquery.Selection) {
		text := content.Text()
		// fmt.Println(text)
		job.Tags = append(job.Tags, text)
	})
	return job, nil

}

func GetJobList(pc_list_url string) (Rows, error) {
	var rows Rows
	var item Item
	g, e := goquery.NewDocument(pc_list_url)
	if e != nil {
		return rows, e
	}

	// 下列内容于 2017年4月4日 20:50:24 抓取
	g.Find("#infolist dl").Each(func(i int, content *goquery.Selection) {

		// 职位详细页链接地址
		item.InfoUrl, _ = content.Find(".t").Attr("href")
		// 验证URL是否我们需要的58PC job info url
		checkLinkIsJobInfo, _ := regexp.MatchString(`http://(?P<site>\w+).58.com/(?P<cate>\w+)/(?P<id>\d+)x.shtml`, item.InfoUrl)
		if checkLinkIsJobInfo {

			item.InfoUrl = GetMobileJobInfoLink(item.InfoUrl)
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

// 通过pc详细页地址 获取手机端详细页地址
func GetMobileJobInfoLink(pc_info_link string) string {
	linkExp := myRegexp{regexp.MustCompile(`http://(?P<site>\w+).58.com/(?P<cate>\w+)/(?P<id>\d+)x.shtml`)}
	linkMap := linkExp.FindStringSubmatchMap(pc_info_link)
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

type Rows []Item

type Item struct {
	Title      string
	Company    string
	CompanyURL string
	Location   string
	Date       string
	InfoUrl    string
}

type myRegexp struct {
	*regexp.Regexp
}

func (r *myRegexp) FindStringSubmatchMap(s string) map[string]string {
	captures := make(map[string]string)

	match := r.FindStringSubmatch(s)
	if match == nil {
		return captures
	}

	for i, name := range r.SubexpNames() {
		//
		if i == 0 {
			continue
		}
		captures[name] = match[i]

	}
	return captures
}

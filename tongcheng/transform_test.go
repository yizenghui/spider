// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tongcheng

import (
	"fmt"
	//	"log"
	"testing"
)

func Test_TransformJob(t *testing.T) {
	var job Job
	job.Title = "服装设计"
	job.Position = "服装设计师"
	job.Category = "服装/纺织/食品"
	job.PositionName = "fuzhuangsheji"
	job.CategoryName = "xiaofeipin"
	job.Location = "广州-海珠-沥滘"
	job.Salery = "面议"
	job.PayType = ""
	job.Description = "本公司主要为服装网店 海珠区南洲路附近<br/>岗位职责：<br/>1 前期核算用料 配比<br/>1 跟进生产进度<br/>2 把控货品质量<br/>3 严格合理利用面辅料 节约成本<br/>4 做事灵活 及时处理突发情况 保证货期 质量<br/>任职资格：<br/>‍‍1二年以上服饰行业工厂或公司跟单经验，能够独立制作服装工艺单<br/>2 熟悉各类面料（毛织/梭织）、服装工艺和生产、熟悉服装加工流程，熟练掌  握服装跟单技能（验货、裁片数量、催货等）‍‍<br/>工作时间：<br/>【上班时间：月休4天，09:00-18:00，】<br/>    特殊情况 弹性工作"
	job.Number = "2"
	job.Education = "学历不限"
	job.WorkYears = "经验不限"
	job.FromURL = "http://m.58.com/gz/xiaofeipin/27916549939388x.shtml"
	job.Company = "广州声色服装有限公司"
	job.CompanyURL = "http://qy.m.58.com/m_detail/18568009933574/"
	job.Linkman = "李先生"
	job.Telephone = "15915892736"
	job.Email = "2859258223@qq.com"
	job.Address = "广州市海珠区南洲路后滘景业工业区"
	job.Lng = "113.324762"
	job.Lat = "23.073643"
	job.Welfare = "五险一金, 社会保险, 活动经费, 发展空间大"
	// var pubJob PubJob
	postJob := TransformJob(job)
	fmt.Println(postJob)
}

func Test_TransformArea(t *testing.T) {
	w := TransformArea("广州-番禺-番禺广场")
	fmt.Println(w)
	w2 := TransformArea("广州-越秀-越秀周边")
	fmt.Println(w2)
}

func Test_TransformCategory(t *testing.T) {
	w := TransformCategory("zpshengchankaifa普工/技工")
	fmt.Println(w)
	w2 := TransformCategory("zptaobao淘宝职位")
	fmt.Println(w2)
}
func Test_TransformExperience(t *testing.T) {
	w := TransformExperience("1-2年")
	fmt.Println(w)
	w2 := TransformExperience("经验不限")
	fmt.Println(w2)
}
func Test_TransformEducation(t *testing.T) {
	w := TransformEducation("学历不限")
	fmt.Println(w)
	w2 := TransformEducation("技校")
	fmt.Println(w2)
	w3 := TransformEducation("中专")
	fmt.Println(w3)
}

func Test_TransformWelfare(t *testing.T) {
	w := TransformWelfare("五险一金,包吃,包住,年底双薪,饭补,绩效奖金,员工旅游,美味下午茶")
	fmt.Println(w)
}
func Test_TransformSalary(t *testing.T) {
	var n, x int
	n, x = TransformSalary("面议")
	fmt.Println(n, x)
	n, x = TransformSalary("5000-8000")
	fmt.Println(n, x)
}

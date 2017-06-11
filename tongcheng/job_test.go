// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tongcheng

import (
	"fmt"
	//	"log"
	"testing"
)

func Test_GetJobInfo(t *testing.T) {
	//	que := yande.NewQuerier("", "1", "10")
	// http://m.58.com/bj/yewu/8822526542722x.shtml
	// http://m.58.com/gz/kefu/29097137457359x.shtml
	job1, err := GetJobInfoData("http://m.58.com/gz/meirongjianshen/29506574767430x.shtml")
	if err != nil {
		panic("spider job data error")

	}
	fmt.Println(job1)
	fmt.Println(job1.Title)
	fmt.Println(job1.Company)
	fmt.Println(job1.CompanyURL)
	// fmt.Println(job1.Tags)
	// job2 := job.GetInfoData("http://m.58.com/gz/renli/27714016037839x.shtml")
	// fmt.Println(job2.Title)
	// job3 := job.GetInfoData("http://m.58.com/gz/zpshengchankaifa/28764443069867x.shtml")
	// fmt.Println(job3.Title)
}

func Test_GetLisRows(t *testing.T) {
	rows1, err := GetJobList("http://gz.58.com/job/")
	if err != nil {
		panic("spider job data error")
	}
	// fmt.Println(rows1)
	for k, v := range rows1 {
		fmt.Printf("%v -> %v\n", k, v.Title)
	}
}

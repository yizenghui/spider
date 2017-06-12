// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tongcheng

import (
	"fmt"
	//	"log"
	"testing"
)

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

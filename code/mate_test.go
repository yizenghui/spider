package code

import (
	"fmt"
	"testing"
)

func Test_MateEducation(t *testing.T) {

	var Education = map[int]string{
		1: "高中J",
		2: "本科",
	}

	var EducationMate = map[string]string{
		//
		"高中": "高中J",
		"本科": "本科",
	}

	edu := Mate("本科", EducationMate, Education)
	edu2 := Mate("高中", EducationMate, Education)
	if edu != 2 || edu2 != 1 {
		fmt.Println(edu, edu2)
	}
}

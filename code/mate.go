package code

// Mate 数据匹配
func Mate(val string, pro map[string]string, tar map[int]string) int {
	var text string
	var edu int
	for k, v := range pro {
		if val == k {
			text = v
			break
		}
	}
	for key, value := range tar {
		if value == text {
			edu = key
			break
		}
	}
	return edu
}

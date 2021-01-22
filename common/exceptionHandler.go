package common

// Must 检测是否出现异常
func Must(err error) {
	if err != nil {
		panic(err)
	}
}

// Must2 检测是否出现异常，同时返回第一个值
func Must2(val interface{}, err error) interface{} {
	Must(err)
	return val
}

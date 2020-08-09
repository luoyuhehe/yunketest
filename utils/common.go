package utils

import "time"

func GetCurrentTimeStr() string {
	return time.Now().Format("2006/01/02 15:04:05")
}

func GetCurrentTime() time.Time {
	return time.Now()
}

func arraySearch(obj interface{}, arr ...interface{}) int {
	for i, v := range arr {
		if obj == v {
			return i
		}
	}
	return -1
}

func InInt32Array(obj int32, arr ...int32) bool {
	newArr := make([]interface{}, 0, len(arr))
	for _, v := range arr {
		newArr = append(newArr, v)
	}
	return arraySearch(obj, newArr...) != -1
}

func InStringArray(obj string, arr ...string) bool {
	newArr := make([]interface{}, 0, len(arr))
	for _, v := range arr {
		newArr = append(newArr, v)
	}
	return arraySearch(obj, newArr...) != -1
}

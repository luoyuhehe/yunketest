package thinkgo

import (
	"fmt"
	"strings"
)

// 错误码字典
var errMsgDict = make(map[int]string)

// 注册错误码
func RegisterErrMsgDict(dict map[int]string) {
	for code, errMsg := range dict {
		if errMsgDict[code] != "" {
			panic(fmt.Errorf("错误码初始化错误,重复定义的code:%d", code))
			return
		}

		if code == 0 || code < -10000000 || code > 99999999 {
			panic(fmt.Errorf("错误码初始化错误,不符合规范的code:%d", code))
		}

		errMsgDict[code] = errMsg
	}
}

func GetErrMsg(code int) string {
	return errMsgDict[code]
}

// BusinessError 业务错误类型
func NewBusinessError(errCode int, errMsg ...string) *BusinessError {
	return &BusinessError{
		ErrorCode: errCode,
		ErrorMsg:  errMsg,
	}
}

type BusinessError struct {
	ErrorCode int
	ErrorMsg  []string
}

func (e *BusinessError) Error() string {
	if len(e.ErrorMsg) != 0 {
		return strings.Join(e.ErrorMsg, "；")
	}
	return GetErrMsg(e.ErrorCode)
}

func GetErrorCode(err error) int {
	switch v := err.(type) {
	case *BusinessError:
		return v.ErrorCode
	default:
		return defaultErrorCode
	}
}

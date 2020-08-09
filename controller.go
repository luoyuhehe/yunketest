package thinkgo

import (
	"gitee.com/luoyusnnu/thinkgo/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	defaultErrorCode   = -1
	defaultErrorMsg    = "未知错误"
	successCode        = 0
	defaultSuccessMsg  = "操作成功"
	autoLogAllResp     = "all"
	autoLogErrorResp   = "error"
	autoLogSuccessResp = "success"
)

type BaseController struct {
}

func NewBaseController() *BaseController {
	return &BaseController{}
}

// isAutoLog 是否自动打印响应日志
func (b *BaseController) isAutoLogResp(code int) bool {
	if !utils.InStringArray(AppConfig.AutoLogResp, autoLogAllResp, autoLogErrorResp, autoLogSuccessResp) {
		return false
	}

	if AppConfig.AutoLogResp == autoLogAllResp {
		return true
	}

	if code == 0 && AppConfig.AutoLogResp == autoLogSuccessResp {
		return true
	}

	if code != 0 && AppConfig.AutoLogResp == autoLogErrorResp {
		return true
	}

	return false
}

// Json 响应json格式结果
func (b *BaseController) Json(ctx *gin.Context, code int, data interface{}, msg string) {
	if data == nil {
		data = map[string]interface{}{}
	}

	obj := gin.H{
		"code": code,
		"data": data,
		"msg":  msg,
	}

	if b.isAutoLogResp(code) {
		if code == 0 {
			b.Info(obj)
		} else {
			b.Error(obj)
		}
	}

	// 开始时间
	ctx.JSON(http.StatusOK, obj)
}

// SuccessResponse 响应成功结果
func (b *BaseController) SuccessResponse(ctx *gin.Context, data interface{}, msg ...string) {
	message := defaultSuccessMsg
	if len(msg) != 0 {
		message = strings.Join(msg, "；")
	}

	b.Json(ctx, successCode, data, message)
}

// ErrorResponse 返回失败响应信息
func (b *BaseController) ErrorResponse(ctx *gin.Context, code int, errMsg ...string) {
	message := b.getDictErrMsg(code)
	if len(errMsg) != 0 {
		message = strings.Join(errMsg, "；")
	}

	b.Json(ctx, code, map[string]interface{}{}, message)
}

// WrapErrorResponse 将错误包装后返回
func (b *BaseController) WrapErrorResponse(ctx *gin.Context, err error) {
	errCode := GetErrorCode(err)
	if b.isAutoLogResp(errCode) {
		b.Error(err)
	}
	b.Json(ctx, errCode, map[string]interface{}{}, err.Error())
}

// ErrorResponse 返回失败响应信息
func (b *BaseController) ErrorResponseWithData(ctx *gin.Context, code int, errMsg string, data interface{}) {
	b.Json(ctx, code, data, errMsg)
}

// getDictErrMsg 解析错误信息
func (b *BaseController) getDictErrMsg(code int) string {
	errMsg, ok := errMsgDict[code]
	if ok {
		return errMsg
	}

	return defaultErrorMsg
}

func (b *BaseController) Infof(format string, args ...interface{}) {
	GetLogger().Infof(format, args)
}

func (b *BaseController) Errorf(format string, args ...interface{}) {
	GetLogger().Errorf(format, args)
}

func (b *BaseController) Debugf(format string, args ...interface{}) {
	GetLogger().Debugf(format, args)
}

func (b *BaseController) Panicf(format string, args ...interface{}) {
	GetLogger().Panicf(format, args)
}

func (b *BaseController) Error(args ...interface{}) {
	GetLogger().Error(args)
}

func (b *BaseController) Info(args ...interface{}) {
	GetLogger().Info(args)
}

func (b *BaseController) Debug(args ...interface{}) {
	GetLogger().Error(args)
}

func (b *BaseController) Panic(args ...interface{}) {
	GetLogger().Panic(args)
}

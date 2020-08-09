package validate

import (
	"fmt"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	validator_engine "github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"strings"

	"github.com/pkg/errors"
	"reflect"
)

type Validator struct {
	translator ut.Translator
	engine     *validator_engine.Validate
}

//var _ binding.StructValidator = &Validator{}
var defaultValidator = NewDefaultValidator()

// GetDefaultValidator 获取默认的数据校验引擎
func GetDefaultValidator() *Validator {
	return defaultValidator
}

// init 初始化语言转换工具
func NewDefaultValidator() *Validator {
	validator := &Validator{
		engine: validator_engine.New(),
	}
	zhTranslator := zh.New()
	uni := ut.New(zhTranslator, zhTranslator)
	var ok bool
	if validator.translator, ok = uni.GetTranslator("zh"); !ok {
		panic("get translator not found")
	}

	if err := zh_translations.RegisterDefaultTranslations(validator.engine, validator.translator); err != nil {
		panic(errors.New("注册Translations失败"))
	}
	validator.engine.SetTagName("validate")
	validator.engine.RegisterTagNameFunc(func(fld reflect.StructField) string {
		return fld.Tag.Get("comment")
	})

	// 注册自定义验证规则失败
	if err := validator.engine.RegisterValidation("is-awesome", validateMyVal); err != nil {
		panic(errors.WithMessage(err, "注册自定义验证规则失败"))
	}

	// 注册自定义校验规则
	if err := validator.registerValidation(); err != nil {
		panic(err)
	}

	if err := validator.translateOverride(); err != nil {
		panic(err)
	}

	return validator
}

// ValidateStruct 验证结构体
func (v *Validator) ValidateStruct(obj interface{}) error {
	if kindOfData(obj) == reflect.Struct {
		if err := v.engine.Struct(obj); err != nil {
			return v.translateAll(err)
		}
	}
	return nil
}

// Engine 获取校验引擎
func (v *Validator) Engine() interface{} {
	return v.engine
}

func kindOfData(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}

// validateMyVal implements validator.Func
func validateMyVal(fl validator_engine.FieldLevel) bool {
	return fl.Field().String() == "awesome"
}

// registerValidation 注册自定义规则
func (v *Validator) registerValidation() error {
	// 注册自定义验证规则失败
	if err := v.engine.RegisterValidation("is-awesome", validateMyVal); err != nil {
		panic(errors.WithMessage(err, "注册自定义验证规则失败"))
	}
	return nil
}

// translateAll 翻译所有验证错误
func (v *Validator) translateAll(err error) error {
	// 解析所有验证错误信息
	if errs, ok := err.(validator_engine.ValidationErrors); ok {
		errMsg := ""
		for _, e := range errs {
			// can translate each error one at a time.
			errMsg = fmt.Sprintf("%s;%s", errMsg, fmt.Sprintf("%s", e.Translate(v.translator)))
		}
		errMsg = strings.Trim(errMsg, ";")
		return errors.New(errMsg)
	}

	return err
}

// translateOverride 错误信息重新定义
func (v *Validator) translateOverride() error {
	if err := v.engine.RegisterTranslation("required", v.translator, func(ut ut.Translator) error {
		return ut.Add("required", "{0}参数不能为空", true)
	}, func(ut ut.Translator, fe validator_engine.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	}); err != nil {
		return err
	}
	return nil
}

package validation

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"sync"

	"go-server/internal/errcode"
	"go-server/internal/response"

	"github.com/gin-gonic/gin/binding"
	zh "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhtranslations "github.com/go-playground/validator/v10/translations/zh"
)

var (
	translator ut.Translator
	initOnce   sync.Once
)

func Init() error {
	var initErr error
	initOnce.Do(func() {
		engine, ok := binding.Validator.Engine().(*validator.Validate)
		if !ok {
			return
		}

		uni := ut.New(zh.New(), zh.New())
		var found bool
		translator, found = uni.GetTranslator("zh")
		if !found {
			initErr = errors.New("zh translator not found")
			return
		}

		engine.RegisterTagNameFunc(func(field reflect.StructField) string {
			name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			switch name {
			case "", "-":
				return field.Name
			default:
				return name
			}
		})

		initErr = zhtranslations.RegisterDefaultTranslations(engine, translator)
	})
	return initErr
}

func BindError(err error) error {
	if err == nil {
		return nil
	}

	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) && translator != nil {
		messages := make([]string, 0, len(validationErrs))
		for _, item := range validationErrs {
			messages = append(messages, item.Translate(translator))
		}
		return response.NewAppError(errcode.ErrInvalidParam.Status(), "INVALID_ARGUMENT", strings.Join(messages, "; "))
	}

	var syntaxErr *json.SyntaxError
	if errors.As(err, &syntaxErr) {
		return response.NewAppError(errcode.ErrInvalidParam.Status(), "INVALID_ARGUMENT", "请求体 JSON 格式错误")
	}

	return errcode.ErrInvalidParam.AsError()
}

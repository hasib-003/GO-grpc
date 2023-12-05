package entity

import (
	"log"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var Validate *validator.Validate
var Trans ut.Translator

func init() {
	translator := en.New()
	uni := ut.New(translator, translator)
	var found bool
	Trans, found = uni.GetTranslator("en")
	if !found {
		log.Fatal("translator not found")
	}

	Validate = validator.New()

	if err := en_translations.RegisterDefaultTranslations(Validate, Trans); err != nil {
		log.Fatal(err)
	}

	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

}

func validate(s interface{}) []FieldError {
	err := Validate.Struct(s)

	if err == nil {
		return nil
	}
	errs := err.(validator.ValidationErrors)
	errcount := len(errs)
	validationErrors := make([]FieldError, errcount)
	for i, e := range errs {
		validationErrors[i] = FieldError{
			Field: e.Field(),
			Error: e.Translate(Trans),
		}
		//	fmt.Println(errors.Unwrap(e).Error())

	}
	return validationErrors
}

package libraries

import (
	"database/sql"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en2 "github.com/go-playground/validator/v10/translations/en"
	"github.com/taufiqkba/go_auth/config"
	"reflect"
	"strings"
)

type Validation struct {
	conn *sql.DB
}

func NewValidation() *Validation {
	conn, err := config.DBConn()
	if err != nil {
		panic(err)
	}

	return &Validation{conn}
}

func (v *Validation) Init() (*validator.Validate, ut.Translator) {
	//translator error
	translator := en.New()
	uni := ut.New(translator, translator)
	trans, _ := uni.GetTranslator("en")

	//	validation using go playground
	validate := validator.New()

	//register default translation english
	err := en2.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		panic(err)
	}

	//change default label error
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		labelName := field.Tag.Get("label")
		return labelName
	})

	validate.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} can't empty", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	//register is_unique validation
	validate.RegisterValidation("is_unique", func(fl validator.FieldLevel) bool {
		params := fl.Param()
		splitParams := strings.Split(params, "-")

		tableName := splitParams[0]
		fieldName := splitParams[1]
		fieldValue := fl.Field().String()

		//	check to database
		return v.checkIsUnique(tableName, fieldName, fieldValue)
	})

	//error message is_unique
	validate.RegisterTranslation("is_unique", trans, func(ut ut.Translator) error {
		return ut.Add("is_unique", "{0} duplicate data", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("is_unique", fe.Field())
		return t
	})
	return validate, trans
}

func (v *Validation) Struct(s interface{}) interface{} {
	validate, trans := v.Init()
	vErrors := make(map[string]interface{})
	err := validate.Struct(s)

	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			vErrors[e.StructField()] = e.Translate(trans)
		}
	}

	if len(vErrors) > 0 {
		return vErrors
	}
	return nil
}

func (v *Validation) checkIsUnique(tableName, fieldName, fieldValue string) bool {
	row, _ := v.conn.Query("select "+fieldName+" from "+tableName+" where "+fieldName+" = ?", fieldValue)
	// example query = select email from users where email = user3@gmail.com

	defer row.Close()

	var result string
	for row.Next() {
		row.Scan(&result)
	}

	return result != fieldValue
}

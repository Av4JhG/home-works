package hw09structvalidator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	errTypeNotStruct = fmt.Errorf("объект не является структурой")
	errLen           = fmt.Errorf("ошибка превышения длины строки")
	errRegexp        = fmt.Errorf("строка не соответствует регулярному выражению")
	errContains      = fmt.Errorf("число не входит в указанное множество")
	errMin           = fmt.Errorf("значение меньше минимально допустимого")
	errMax           = fmt.Errorf("значение больше максимально допустимого")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var resultErr string
	for _, item := range v {
		resultErr += fmt.Sprintf("field: %s | err: %s\n", item.Field, item.Err)
	}

	return resultErr
}

func Validate(v interface{}) error {
	var errors ValidationErrors
	var err error

	structType := reflect.TypeOf(v)
	structValue := reflect.ValueOf(v)

	// проверим что нам пришла структура.
	if structType.Kind().String() != "struct" {
		return errTypeNotStruct
	}

	// перебираем поля типа.
	for i := 0; i < structType.NumField(); i++ {
		fieldType := structType.Field(i)
		fieldValue := structValue.Field(i)
		tagValue := fieldType.Tag.Get("validate")

		if tagValue == "" {
			continue
		}

		// валидируем, набираем слайс ошибок.
		switch fieldType.Type.String() {
		case "string":
			errors, err = validateString(fieldType.Name, fieldValue.String(), tagValue, errors)
			if err != nil {
				return err
			}

		case "[]string":
			for _, item := range fieldValue.Interface().([]string) {
				errors, err = validateString(fieldType.Name, item, tagValue, errors)
				if err != nil {
					return err
				}
			}

		case "int":
			errors, err = validateInt(fieldType.Name, fieldValue.Interface().(int), tagValue, errors)
			if err != nil {
				return err
			}

		case "[]int":
			for _, item := range fieldValue.Interface().([]int) {
				errors, err = validateInt(fieldType.Name, item, tagValue, errors)
				if err != nil {
					return err
				}
			}

		default:
			continue
		}
	}

	return errors
}

// валидация строк.
func validateString(
	fieldName string,
	fieldValue string,
	rulesStr string,
	errors ValidationErrors,
) (ValidationErrors, error) {
	rules := strings.Split(rulesStr, "|")

	for _, value := range rules {
		rulesSlice := strings.Split(value, ":")

		switch rulesSlice[0] {
		case "len":
			intValue, err := strconv.Atoi(rulesSlice[1])
			if err != nil {
				return nil, err
			}

			if len(fieldValue) != intValue {
				errors = append(errors, ValidationError{
					Field: fieldName,
					Err:   errLen,
				})
			}
		case "regexp":
			matchString, err := regexp.MatchString(rulesSlice[1], fieldValue)
			if err != nil {
				return nil, err
			}

			if !matchString {
				errors = append(errors, ValidationError{
					Field: fieldName,
					Err:   errRegexp,
				})
			}
		case "in":
			for _, item := range strings.Split(rulesSlice[1], ",") {
				if !strings.Contains(fieldValue, item) {
					errors = append(errors, ValidationError{
						Field: fieldName,
						Err:   errContains,
					})
				}
			}
		default:
			continue
		}
	}

	return errors, nil
}

// валидация целых чисел.
func validateInt(fieldName string, fieldValue int, rulesStr string, errors ValidationErrors) (ValidationErrors, error) {
	rules := strings.Split(rulesStr, "|")

	for _, value := range rules {
		rulesSlice := strings.Split(value, ":")
		switch rulesSlice[0] {
		case "min":
			intValue, err := strconv.Atoi(rulesSlice[1])
			if err != nil {
				return nil, err
			}

			if fieldValue < intValue {
				errors = append(errors, ValidationError{
					Field: fieldName,
					Err:   errMin,
				})
			}
		case "max":
			intValue, err := strconv.Atoi(rulesSlice[1])
			if err != nil {
				return nil, err
			}

			if fieldValue > intValue {
				errors = append(errors, ValidationError{
					Field: fieldName,
					Err:   errMax,
				})
			}
		case "in":
			contains := false
			for _, item := range strings.Split(rulesSlice[1], ",") {
				if strings.Contains(strconv.Itoa(fieldValue), item) {
					contains = true
				}
			}
			if !contains {
				errors = append(errors, ValidationError{
					Field: fieldName,
					Err:   errContains,
				})
			}
		default:
			continue
		}
	}

	return errors, nil
}

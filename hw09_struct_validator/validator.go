package hw09_struct_validator //nolint:golint,stylecheck

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var resultString string
	for _, error := range v {
		resultString += fmt.Sprintf("error in field: %v, %v\n", error.Field, error.Err)
	}
	return resultString
}

func Validate(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("expected struct, but received %T", v)
	}

	var validationErrors ValidationErrors
	for i := 0; i < val.NumField(); i++ {
		if !val.Field(i).CanInterface() {
			continue
		}
		tag := string(val.Type().Field(i).Tag)
		for _, iterTag := range strings.Split(tag, " ") {
			if strings.Split(iterTag, ":")[0] == "validate" {
				normalizedTag := strings.Trim(iterTag[9:], "\"")
				value := val.Field(i)
				valueName := val.Type().Field(i).Name
				err := validateField(normalizedTag, value)
				if err != nil {
					er := ValidationError{
						valueName,
						err,
					}
					validationErrors = append(validationErrors, er)
				}
			}
		}
	}
	fmt.Printf("%v\n", validationErrors)
	if validationErrors != nil {
		return validationErrors
	}
	return nil
}

func validateField(validationTag string, val reflect.Value) error {
	subTags := strings.Split(validationTag, "|")
	rules := map[string]string{}
	for _, tag := range subTags {
		validationRule := strings.Split(tag, ":")[0]
		validationValue := strings.Split(tag, ":")[1]
		rules[validationRule] = validationValue
	}

	switch val.Kind() {
	case reflect.Int:
		var value []int64
		value = append(value, val.Int())
		err := parseIntField(rules, value)
		if err != nil {
			return err
		}
		return nil
	case reflect.String:
		var value []string
		value = append(value, val.String())
		err := parseStringField(rules, value)
		if err != nil {
			return err
		}
		return nil
	case reflect.SliceOf(val.Type()).Kind():
		switch t := val.Interface().(type) {
		case []string:
			var value []string
			value = append(value, t...)
			if err := parseStringField(rules, value); err != nil {
				return err
			}
			return nil
		case []int:
			var value []int64
			for _, v := range t {
				value = append(value, int64(v))
			}
			if err := parseIntField(rules, value); err != nil {
				return err
			}
		}
	}
	return nil
}

func parseStringField(rules map[string]string, value []string) error {
	for _, v := range value {
		for ruleName, ruleValue := range rules {
			if ruleName == "in" {
				if !strings.Contains(ruleValue, v) {
					return fmt.Errorf("error: %v not in [%v] set", v, ruleValue)
				}
			}
			if ruleName == "len" {
				length, err := strconv.Atoi(ruleValue)
				if err != nil {
					return fmt.Errorf("cast error")
				}
				if len(v) != length {
					return fmt.Errorf("error: length of \"%v\" not equal %v", v, ruleValue)
				}
			}
			if ruleName == "regexp" {
				ruleValue = strings.ReplaceAll(ruleValue, "\\\\", `\`)
				validString := regexp.MustCompile(ruleValue)
				if !validString.MatchString(v) {
					return fmt.Errorf("error: %v does not match %v regex expression", v, ruleValue)
				}
			}
		}
	}
	return nil
}

func parseIntField(rules map[string]string, value []int64) error {
	for _, v := range value {
		for ruleName, ruleValue := range rules {
			if ruleName == "in" {
				if !strings.Contains(ruleValue, strconv.FormatInt(v, 10)) {
					return fmt.Errorf("error: %v not in [%v] set", v, ruleValue)
				}
			}
			if ruleName == "min" {
				min, _ := strconv.ParseInt(ruleValue, 10, 64)
				if v < min {
					return fmt.Errorf("error: %v less than %v", v, ruleValue)
				}
			}
			if ruleName == "max" {
				max, _ := strconv.ParseInt(ruleValue, 10, 64)
				if v > max {
					return fmt.Errorf("error: %v greater than %v", v, ruleValue)
				}
			}
		}
	}
	return nil
}

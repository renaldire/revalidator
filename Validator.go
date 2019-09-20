package Validator

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Rules struct {
	Value interface{}
	Rule  string
}

const (
	INVALID_MIN              = "%s must be not less than %d"
	INVALID_MIN_LENGTH       = "%s must be at least %d characters"

	INVALID_MAX              = "%s must be not more than %d"
	INVALID_MAX_LENGTH       = "%s must be not more than %d characters"

	INVALID_REQUIRED         = "%s is required."
	INVALID_EMAIL            = "%s must be in e-mail format."
	INVALID_NUMERIC          = "%s must be a numeric"
	INVALID_DATE             = "%s must be in yyyy-mm-dd format."
	INVALID_DATETIME_FORMAT  = "%s must be in yyyy-MM-dd hh:mm:ss format."
	INVALID_TIME_FORMAT      = "%s must be in hh:mm:ss format."
	DEFAULT_DATE_FORMAT      = "2006-01-02"
	DEFAULT_DATE_TIME_FORMAT = "2006-01-02T15:04:05"
	DEFAULT_TIME_FORMAT      = "15:04:05"
)

func min_length(key, value string, min int) (error) {
	if (len(value) < min) {
		return fmt.Errorf(INVALID_MIN_LENGTH, key, min)
	}
	return nil
}
func min(key string, value int, min int) (error) {
	if (value < min) {
		return fmt.Errorf(INVALID_MIN, key, min)
	}
	return nil
}
func max_length(key, value string, min int) (error) {
	if (len(value) > min) {
		return fmt.Errorf(INVALID_MAX_LENGTH, key, min)
	}
	return nil
}
func max(key string, value, max int) (error) {
	if (value > max) {
		return fmt.Errorf(INVALID_MAX, key, max)
	}
	return nil
}
func required(key, value string) (error) {
	if value == "" {
		return fmt.Errorf(INVALID_REQUIRED, key)
	}
	return nil
}
func numeric(key, value string) (error) {
	_, err := strconv.Atoi(value)

	if (err != nil) {
		return fmt.Errorf(INVALID_NUMERIC, key)
	}

	return nil
}
func email(key, value string) (error) {
	rxEmail := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if len(value) > 254 || !rxEmail.MatchString(value) {
		return fmt.Errorf(INVALID_EMAIL, key)
	}
	return nil
}
func getInt(rule, specific string) (int, error) {
	data := strings.Split(rule, specific)
	value := strings.Split(data[1], "|")[0]
	return strconv.Atoi(value)
}
func getString(rule, specific string)(string){
	data := strings.Split(rule, specific)
	value := strings.Split(data[1], "|")[0]
	return value
}
func getByLayout(rule,specific string,layout string)(time.Time,error){
	data := strings.Split(rule, specific)
	value := strings.Split(data[1], "|")[0]
	return time.Parse(layout, value)
}
func date(key, value string) (error) {
	_, err := time.Parse(DEFAULT_DATE_FORMAT, value)
	if err != nil {
		return fmt.Errorf(INVALID_DATE, key)
	}
	return nil
}
func datetime(key, value string) (error) {
	_, err := time.Parse(DEFAULT_DATE_TIME_FORMAT, value)
	if err != nil {
		return fmt.Errorf(INVALID_DATETIME_FORMAT, key)
	}
	return nil
}
func timeFormat(key, value string) (error) {
	_, err := time.Parse(DEFAULT_TIME_FORMAT, value)
	if err != nil {
		return fmt.Errorf(INVALID_TIME_FORMAT, key)
	}
	return nil
}

func Validate(validator map[string]Rules) (errors []error) {

	for k, v := range (validator) {
		if strings.Contains(v.Rule, "allowempty") {
			if (v.Value == "" || v.Value == nil) {
				continue
			}
		}
		if strings.Contains(v.Rule,"required_if"){
			otherField := getString(v.Rule, "required_if:")
			data:=strings.Split(otherField, ",")

			testValidator:=validator[data[0]]
			if testValidator.Value==data[1]{
				err := required(k, fmt.Sprintf("%v", v.Value))
				if err != nil {
					errors = append(errors, err)
				}
			}else{
				continue
			}
		}
		if strings.Contains(v.Rule, "required") {
			err := required(k, fmt.Sprintf("%v", v.Value))
			if err != nil {
				errors = append(errors, err)
			}
		}
		if strings.Contains(v.Rule, "numeric") {
			err := numeric(k, fmt.Sprintf("%v", v.Value))
			if err != nil {
				errors = append(errors, err)
			}
		}
		if strings.Contains(v.Rule, "min:") {
			minValue, err := getInt(v.Rule, "min:")
			if err != nil {
				panic("min validator rule value must be numeric")
				os.Exit(1)
			}

			switch v.Value.(type) {
			case int:
				num, _ := v.Value.(int)
				err = min(k, num, minValue)
			case string:
				err = min_length(k, fmt.Sprintf("%v", v.Value), minValue)
			}

			if err != nil {
				errors = append(errors, err)
			}
		}
		if strings.Contains(v.Rule, "max:") {
			maxValue, err := getInt(v.Rule, "max:")
			if err != nil {
				panic("min validator rule value must be numeric")
				os.Exit(1)
			}

			switch v.Value.(type) {
			case int:
				num, _ := v.Value.(int)
				err = max(k, num, maxValue)
			case string:
				err = max_length(k, fmt.Sprintf("%v", v.Value), maxValue)
			}

			if err != nil {
				errors = append(errors, err)
			}
		}
		if strings.Contains(v.Rule, "email") {
			err := email(k, fmt.Sprintf("%v", v.Value))
			if err != nil {
				errors = append(errors, err)
			}
		}

		if strings.Contains(v.Rule, "datetime") {
			err := datetime(k, fmt.Sprintf("%v", v.Value))
			if err != nil {
				errors = append(errors, err)
			}
		} else if strings.Contains(v.Rule, "date") {
			err := date(k, fmt.Sprintf("%v", v.Value))
			if err != nil {
				errors = append(errors, err)
			}
		} else if strings.Contains(v.Rule, "time") {
			err := timeFormat(k, fmt.Sprintf("%v", v.Value))
			if err != nil {
				errors = append(errors, err)
			}
		}
	}
	return
}

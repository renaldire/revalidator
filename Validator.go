package Validator

import (
	"fmt"
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
	LOG = "[500] Internal Server Error. %s attribute must be a numeric int.\n"
	INVALID_MIN            = "%s must be not less than %d"
	INVALID_MIN_CHARACTERS = "%s must be at least %d characters"

	INVALID_MAX            = "%s must be not more than %d"
	INVALID_MAX_CHARACTERS = "%s must be not more than %d characters"

	INVALID_REQUIRED         = "%s is required."
	INVALID_EMAIL            = "%s must be in e-mail format."
	INVALID_IN               = "%s is invalid."
	INVALID_NUMERIC          = "%s must be a numeric"
	INVALID_DATE             = "%s must be in yyyy-mm-dd format."
	INVALID_DATETIME_FORMAT  = "%s must be in yyyy-MM-dd hh:mm:ss format."
	INVALID_TIME_FORMAT      = "%s must be in hh:mm:ss format."
	DEFAULT_DATE_FORMAT      = "2006-01-02"
	DEFAULT_DATE_TIME_FORMAT = "2006-01-02T15:04:05"
	DEFAULT_TIME_FORMAT      = "15:04:05"
)

func date(validator Rules,key string,errors *[]error)(bool){
	if !strings.Contains(validator.Rule, "date") {
		return false
	}
	_, err := time.Parse(DEFAULT_DATE_FORMAT, fmt.Sprintf("%s",validator.Value))
	if err != nil {
		*errors=append(*errors,fmt.Errorf(INVALID_DATE, key))
	}
	return true
}
func datetime(validator Rules,key string,errors *[]error)(bool){
	if !strings.Contains(validator.Rule, "datetime") {
		return false
	}
	_, err := time.Parse(DEFAULT_DATE_TIME_FORMAT, fmt.Sprintf("%s",validator.Value))
	if err != nil {
		*errors=append(*errors,fmt.Errorf(INVALID_DATETIME_FORMAT, key))
	}
	return true
}
func timeFormat(validator Rules,key string,errors *[]error)(bool) {
	if !strings.Contains(validator.Rule, "|time") && !strings.HasPrefix(validator.Rule,"time") {
		return false
	}
	_, err := time.Parse(DEFAULT_TIME_FORMAT, fmt.Sprintf("%s",validator.Value))
	if err != nil {
		*errors=append(*errors,fmt.Errorf(INVALID_TIME_FORMAT, key))
	}
	return true
}
func getInt(rule, specific string) (int, error) {
	data := strings.Split(rule, specific)
	value := strings.Split(data[1], "|")[0]
	return strconv.Atoi(value)
}
func getString(rule, specific string) (string) {
	data := strings.Split(rule, specific)
	value := strings.Split(data[1], "|")[0]
	return value
}
func numeric(validator Rules,key string,errors *[]error){
	if !strings.Contains(validator.Rule, "numeric") {
		return
	}

	_, err := strconv.Atoi(fmt.Sprintf("%s",validator.Value))
	if (err != nil) {
		*errors=append(*errors,fmt.Errorf(INVALID_NUMERIC, key))
	}
	regex := regexp.MustCompile("[0-9]*")
	if !regex.MatchString(fmt.Sprintf("%s",validator.Value)) {
		*errors=append(*errors,fmt.Errorf(INVALID_EMAIL, key))
	}
}
func allowEmpty(validator Rules)(bool){
	if !strings.Contains(validator.Rule, "allowempty") {
		return false
	}

	if (validator.Value == "" || validator.Value == nil) {
		return true
	}
	return false
}
func requiredIf(validator Rules,validators map[string]Rules,key string,errors *[]error)(bool){
	if !strings.Contains(validator.Rule, "required_if") {
		return false
	}
	otherField := getString(validator.Rule, "required_if:")

	data := strings.Split(otherField, ",")
	targetKey:=data[0]
	targetValue:=data[1]

	targetValidator := validators[targetKey]
	if targetValidator.Value == targetValue {
		required(validator,key,errors)
	}
	return true
}
func min(validator Rules,key string,errors *[]error){
	if !strings.Contains(validator.Rule, "min:") {
		return
	}

	minValue, err := getInt(validator.Rule, "min:")
	if err != nil {
		fmt.Printf(LOG,"min")
		return
	}

	switch validator.Value.(type) {
	case int:
		if validator.Value.(int)<minValue{
			*errors=append(*errors,fmt.Errorf(INVALID_MIN,key,minValue))
		}
	case string:
		if len(fmt.Sprintf("%s",validator.Value))<minValue{
			*errors=append(*errors,fmt.Errorf(INVALID_MIN_CHARACTERS,key,minValue))
		}
	}
}
func max(validator Rules,key string,errors *[]error){
	if !strings.Contains(validator.Rule, "max:") {
		return
	}
	maxValue, err := getInt(validator.Rule, "max:")
	if err != nil {
		fmt.Printf(LOG,"max")
		return
	}

	switch validator.Value.(type) {
	case int:
		if validator.Value.(int)>maxValue{
			*errors=append(*errors,fmt.Errorf(INVALID_MAX,key,maxValue))
		}
	case string:
		if len(fmt.Sprintf("%s",validator.Value))>maxValue{
			*errors=append(*errors,fmt.Errorf(INVALID_MAX_CHARACTERS,key,maxValue))
		}
	}
}
func email(validator Rules,key string,errors *[]error){
	if !strings.Contains(validator.Rule, "email") {
		return
	}

	regex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if len(fmt.Sprintf("%s",validator.Value)) > 254 || !regex.MatchString(fmt.Sprintf("%s",validator.Value)) {
		*errors=append(*errors,fmt.Errorf(INVALID_EMAIL, key))
	}
	return
}
func in(validator Rules,key string,errors *[]error){
	if !strings.Contains(validator.Rule, "|in:") && !strings.HasPrefix(validator.Rule,"in:"){
		return
	}

	data:=strings.Split(getString(validator.Rule, "in:"),",")
	for _, value := range (data) {
		if value == validator.Value {
			return
		}
	}
	*errors=append(*errors,fmt.Errorf(INVALID_IN,key))
}
func required(validator Rules,key string,errors *[]error){
	if !strings.Contains(validator.Rule, "required") {
		return
	}

	if validator.Value == "" || validator.Value == nil {
		*errors=append(*errors,fmt.Errorf(INVALID_REQUIRED, key))
	}
}
func calendar(validator Rules,key string,errors *[]error){
	if datetime(validator,key,errors)==false {
		if date(validator, key, errors)==false{
			timeFormat(validator,key,errors)
		}
	}
}
func skipValidateOtherField(validator Rules,validators map[string]Rules,key string,errors *[]error)(bool){
	if allowEmpty(validator)==true { return true }
	if requiredIf(validator,validators,key,errors)==true { return true }
	return false
}

func Validate(validators map[string]Rules) (errors []error) {
	for key, validator := range (validators) {
		if skipValidateOtherField(validator,validators,key,&errors)==true {continue}

		required(validator,key,&errors)
		numeric(validator,key,&errors)
		min(validator,key,&errors)
		max(validator,key,&errors)
		email(validator,key,&errors)
		in(validator,key,&errors)
		calendar(validator,key,&errors)
	}
	return
}

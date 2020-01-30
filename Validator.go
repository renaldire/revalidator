package Validator

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func unique(validator Rules, key string, errors *[]error) {
	if !strings.Contains(validator.Rule, "unique:") {
		return
	}

	data := strings.Split(getString(validator.Rule, "unique:"), ",")

	if len(data) != 2 {
		return
	}

	table := data[0]
	column := data[1]

	statement := fmt.Sprintf("SELECT EXISTS(SELECT %s FROM %s WHERE %s=$1)", column, table, column)
	stmt, err := getDb().Prepare(statement)

	if err == sql.ErrNoRows {
		*errors = append(*errors, err)
		return
	} else if err!=nil{
		panic(err.Error())
	}

	row := stmt.QueryRow(validator.Value)
	var result bool
	row.Scan(&result)

	if result == true {
		err = fmt.Errorf(invalidUnique, key)
		*errors = append(*errors, err)
		return
	}
}
func date(validator Rules, key string, errors *[]error) (bool) {
	if !strings.Contains(validator.Rule, "date") {
		return false
	}
	_, err := time.Parse(defaultDateFormat, fmt.Sprintf("%s", validator.Value))
	if err != nil {
		*errors = append(*errors, fmt.Errorf(invalidDate, key))
	}
	return true
}
func datetime(validator Rules, key string, errors *[]error) (bool) {
	if !strings.Contains(validator.Rule, "datetime") {
		return false
	}
	_, err := time.Parse(defaultDateTimeFormat, fmt.Sprintf("%s", validator.Value))
	if err != nil {
		*errors = append(*errors, fmt.Errorf(invalidDateTime, key))
	}
	return true
}
func timeFormat(validator Rules, key string, errors *[]error) (bool) {
	if !strings.Contains(validator.Rule, "|time") && !strings.HasPrefix(validator.Rule, "time") {
		return false
	}
	_, err := time.Parse(defaultTimeFormat, fmt.Sprintf("%s", validator.Value))
	if err != nil {
		*errors = append(*errors, fmt.Errorf(invalidTime, key))
	}
	return true
}
func getInt(rule, attribute string) (int, error) {
	data := strings.Split(rule, attribute)
	value := strings.Split(data[1], "|")[0]
	return strconv.Atoi(value)
}
func getString(rule, attribute string) (string) {
	data := strings.Split(rule, attribute)
	value := strings.Split(data[1], "|")[0]
	return value
}
func numeric(validator Rules, key string, errors *[]error) {
	if !strings.Contains(validator.Rule, "numeric") {
		return
	}

	_, err := strconv.Atoi(fmt.Sprintf("%s", validator.Value))
	if (err != nil) {
		*errors = append(*errors, fmt.Errorf(invalidNumeric, key))
	}
	regex := regexp.MustCompile("[0-9]*")
	if !regex.MatchString(fmt.Sprintf("%s", validator.Value)) {
		*errors = append(*errors, fmt.Errorf(invalidEmail, key))
	}
}
func allowEmpty(validator Rules) (bool) {
	if !strings.Contains(validator.Rule, "allowempty") {
		return false
	}

	if (validator.Value == "" || validator.Value == nil) {
		return true
	}
	return false
}
func requiredIf(validator Rules, validators map[string]Rules, key string, errors *[]error) (bool) {
	if !strings.Contains(validator.Rule, "required_if") {
		return false
	}
	otherField := getString(validator.Rule, "required_if:")

	data := strings.Split(otherField, ",")
	targetKey := data[0]
	targetValue := data[1]

	targetValidator := validators[targetKey]
	if targetValidator.Value == targetValue {
		required(validator, key, errors)
	}
	return true
}
func min(validator Rules, key string, errors *[]error) {
	if !strings.Contains(validator.Rule, "min:") {
		return
	}

	minValue, err := getInt(validator.Rule, "min:")
	if err != nil {
		fmt.Printf(log, "min")
		return
	}

	switch validator.Value.(type) {
	case int:
		if validator.Value.(int) < minValue {
			*errors = append(*errors, fmt.Errorf(invalidMin, key, minValue))
		}
	case string:
		if len(fmt.Sprintf("%s", validator.Value)) < minValue {
			*errors = append(*errors, fmt.Errorf(invalidMinCharacters, key, minValue))
		}
	}
}
func max(validator Rules, key string, errors *[]error) {
	if !strings.Contains(validator.Rule, "max:") {
		return
	}
	maxValue, err := getInt(validator.Rule, "max:")
	if err != nil {
		fmt.Printf(log, "max")
		return
	}

	switch validator.Value.(type) {
	case int:
		if validator.Value.(int) > maxValue {
			*errors = append(*errors, fmt.Errorf(invalidMax, key, maxValue))
		}
	case string:
		if len(fmt.Sprintf("%s", validator.Value)) > maxValue {
			*errors = append(*errors, fmt.Errorf(invalidMaxCharacters, key, maxValue))
		}
	}
}
func email(validator Rules, key string, errors *[]error) {
	if !strings.Contains(validator.Rule, "email") {
		return
	}

	regex := regexp.MustCompile(regexEmail)
	if len(fmt.Sprintf("%s", validator.Value)) > 254 || !regex.MatchString(fmt.Sprintf("%s", validator.Value)) {
		*errors = append(*errors, fmt.Errorf(invalidEmail, key))
	}
	return
}
func in(validator Rules, key string, errors *[]error) {
	if !strings.Contains(validator.Rule, "|in:") && !strings.HasPrefix(validator.Rule, "in:") {
		return
	}

	data := strings.Split(getString(validator.Rule, "in:"), ",")
	for _, value := range (data) {
		if value == validator.Value {
			return
		}
	}
	*errors = append(*errors, fmt.Errorf(invalidIn, key))
}
func required(validator Rules, key string, errors *[]error) {
	if !strings.Contains(validator.Rule, "required") {
		return
	}

	if validator.Value == "" || validator.Value == nil {
		*errors = append(*errors, fmt.Errorf(invalidRequired, key))
	}
}
func calendar(validator Rules, key string, errors *[]error) {
	if datetime(validator, key, errors) == false {
		if date(validator, key, errors) == false {
			timeFormat(validator, key, errors)
		}
	}
}
func allowedEmptyOrDependsOtherField(validator Rules, validators map[string]Rules, key string, errors *[]error) (bool) {
	if allowEmpty(validator) == true {
		return true
	}
	if requiredIf(validator, validators, key, errors) == true {
		return true
	}
	return false
}
func startsWith(validator Rules, key string, errors *[]error){
	if !strings.Contains(validator.Rule, "starts_with:") {
		return
	}

	targetPrefix:=getString(validator.Rule,"starts_with:")
	value:=fmt.Sprintf("%s",validator.Value)

	if len(value)<len(targetPrefix){
		*errors = append(*errors, fmt.Errorf(invalidStartsWith, key,targetPrefix))
		return
	}

	if value[0:len(targetPrefix)-1]!=targetPrefix{
		*errors = append(*errors, fmt.Errorf(invalidStartsWith, key,targetPrefix))
		return
	}
}
func endsWith(validator Rules, key string, errors *[]error){
	if !strings.Contains(validator.Rule, "ends_with:") {
		return
	}

	targetPrefix:=getString(validator.Rule,"ends_with:")
	value:=fmt.Sprintf("%s",validator.Value)

	if len(value)<len(targetPrefix){
		*errors = append(*errors, fmt.Errorf(invalidEndsWith, key,targetPrefix))
		return
	}

	if value[len(value)-len(targetPrefix):]!=targetPrefix{
		*errors = append(*errors, fmt.Errorf(invalidEndsWith, key,targetPrefix))
		return
	}
}
func Validate(validators map[string]Rules) (errors []error) {
	for key, validator := range (validators) {
		if allowedEmptyOrDependsOtherField(validator, validators, key, &errors) == true {
			continue
		}

		required(validator, key, &errors)
		numeric(validator, key, &errors)
		min(validator, key, &errors)
		max(validator, key, &errors)
		email(validator, key, &errors)
		in(validator, key, &errors)
		unique(validator, key, &errors)
		calendar(validator, key, &errors)
		startsWith(validator,key,&errors)
		endsWith(validator,key,&errors)
	}
	return
}

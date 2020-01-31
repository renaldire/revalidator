package Validator

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func getRule(rules, targetRule string) (string,bool) {
	//return !strings.Contains(rules, "|"+targetRule) && !strings.HasPrefix(rules, targetRule)
	for _,v:=range strings.Split(rules,"|") {
		if strings.HasPrefix(v, targetRule){
			return v,true
		}
	}
	return "",false
}
func unique(validator Rules, key string, errors *[]error) {
	uniqueRule,ok:=getRule(validator.Rule, "unique")
	if !ok {
		return
	}

	data := strings.Split(getStringRuleValue(uniqueRule), ",")

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
	} else if err != nil {
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
func date(validator Rules, key string, errors *[]error) bool {
	_,ok:=getRule(validator.Rule, "date")
	if !ok {
		return false
	}
	_, err := time.Parse(defaultDateFormat, fmt.Sprintf("%s", validator.Value))
	if err != nil {
		*errors = append(*errors, fmt.Errorf(invalidDate, key))
	}
	return true
}
func datetime(validator Rules, key string, errors *[]error) bool {
	_,ok:=getRule(validator.Rule, "datetime")
	if !ok {
		return false
	}
	_, err := time.Parse(defaultDateTimeFormat, fmt.Sprintf("%s", validator.Value))
	if err != nil {
		*errors = append(*errors, fmt.Errorf(invalidDateTime, key))
	}
	return true
}
func timeFormat(validator Rules, key string, errors *[]error) bool {
	_,ok:=getRule(validator.Rule, "time")
	if !ok {
		return false
	}
	_, err := time.Parse(defaultTimeFormat, fmt.Sprintf("%s", validator.Value))
	if err != nil {
		*errors = append(*errors, fmt.Errorf(invalidTime, key))
	}
	return true
}
func getIntRuleValue(rule string) (int, error) {
	data := strings.Split(rule, ":")
	value := data[1]
	return strconv.Atoi(value)
}
func getStringRuleValue(rule string) string {
	data := strings.Split(rule,":")
	value := data[1]
	return value
}
func checkRegex(validator Rules, key, regexPattern string) bool {
	regex := regexp.MustCompile(regexPattern)
	if !regex.MatchString(fmt.Sprintf("%s", validator.Value)) {
		return false
	}
	return true
}
func numeric(validator Rules, key string, errors *[]error) {
	_,ok:=getRule(validator.Rule, "numeric")
	if !ok {
		return
	}
	if checkRegex(validator, key, RegexNotNumeric) == true {
		*errors = append(*errors, fmt.Errorf(invalidNumeric, key))
	}
}
func allowEmpty(validator Rules) bool {
	_,ok:=getRule(validator.Rule, "allowempty")
	if !ok {
		return false
	}

	if validator.Value == "" || validator.Value == nil {
		return true
	}
	return false
}
func requiredIf(validator Rules, validators map[string]Rules, key string, errors *[]error) bool {
	requiredIfRule,ok:=getRule(validator.Rule, "required_if:")
	if !ok {
		return false
	}

	otherField := getStringRuleValue(requiredIfRule)

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
	minRule,ok:=getRule(validator.Rule, "min:")
	if !ok {
		return
	}

	minValue, err := getIntRuleValue(minRule)
	if err != nil {
		panic("min attribute value must be numeric")
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
	maxRule,ok:=getRule(validator.Rule, "max:")
	if !ok {
		return
	}
	maxValue, err := getIntRuleValue(maxRule)
	if err != nil {
		panic("max attribute value must be numeric")
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
	_,ok:=getRule(validator.Rule, "email")
	if !ok {
		return
	}

	value := fmt.Sprintf("%s", validator.Value)

	if len(value) > 254 || checkRegex(validator, key, RegexEmail) == false {
		*errors = append(*errors, fmt.Errorf(invalidEmail, key))
	}
	return
}
func in(validator Rules, key string, errors *[]error) {
	inRule,ok:=getRule(validator.Rule, "in:")
	if !ok {
		return
	}

	data := strings.Split(getStringRuleValue(inRule), ",")
	for _, value := range data {
		if value == validator.Value {
			return
		}
	}
	*errors = append(*errors, fmt.Errorf(invalidIn, key))
}
func required(validator Rules, key string, errors *[]error) {
	_,ok:=getRule(validator.Rule, "required")
	if !ok {
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
func allowedEmptyOrDependsOtherField(validator Rules, validators map[string]Rules, key string, errors *[]error) bool {
	if allowEmpty(validator) == true {
		return true
	}
	if requiredIf(validator, validators, key, errors) == true {
		return true
	}
	return false
}
func startsWith(validator Rules, key string, errors *[]error) {
	startsWithRule,ok:=getRule(validator.Rule, "starts_with")
	if !ok {
		return
	}

	targetPrefix := getStringRuleValue(startsWithRule)
	value := fmt.Sprintf("%s", validator.Value)

	if len(value) < len(targetPrefix) {
		*errors = append(*errors, fmt.Errorf(invalidStartsWith, key, targetPrefix))
		return
	}

	if value[0:len(targetPrefix)-1] != targetPrefix {
		*errors = append(*errors, fmt.Errorf(invalidStartsWith, key, targetPrefix))
	}
}
func endsWith(validator Rules, key string, errors *[]error) {
	endsWithRule,ok:=getRule(validator.Rule, "ends_with")
	if !ok {
		return
	}

	targetPrefix := getStringRuleValue(endsWithRule)
	value := fmt.Sprintf("%s", validator.Value)

	if len(value) < len(targetPrefix) {
		*errors = append(*errors, fmt.Errorf(invalidEndsWith, key, targetPrefix))
		return
	}

	if value[len(value)-len(targetPrefix):] != targetPrefix {
		*errors = append(*errors, fmt.Errorf(invalidEndsWith, key, targetPrefix))
	}
}
func regex(validator Rules, key string, errors *[]error) {
	regexRule,ok:=getRule(validator.Rule, "regex:")
	if !ok {
		return
	}

	targetRegex := getStringRuleValue(regexRule)

	if checkRegex(validator, key, targetRegex) == false {
		*errors = append(*errors, fmt.Errorf(invalidRegex, key))
	}
}

func Validate(validators map[string]Rules) (errors []error) {
	for key, validator := range validators {
		validator.Value = strings.TrimSpace(fmt.Sprintf("%s", validator.Value))

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
		startsWith(validator, key, &errors)
		endsWith(validator, key, &errors)
		regex(validator, key, &errors)
	}
	return
}

func ValidateStruct(param interface{}) (errors []error) {
	structParam := reflect.TypeOf(param)
	valueParam := reflect.ValueOf(param)

	isParamTypeOfStruct := structParam.Kind() == reflect.Struct

	if isParamTypeOfStruct != true {
		panic("Parameter of ValidateStruct must be a struct type")
		return
	}

	validators := make(map[string]Rules)

	// loop through each struct attribute
	for i := 0; i < structParam.NumField(); i++ {
		currField := structParam.Field(i)
		currValue := valueParam.Field(i)

		key, isUsingJsonTag := currField.Tag.Lookup("json")
		if !isUsingJsonTag {
			key = structParam.Field(i).Name
		}

		value := currValue.Interface()
		rule, isUsingRuleTag := currField.Tag.Lookup("rule")

		if !isUsingRuleTag {
			continue
		}

		validators[key] = Rules{value, rule}
	}

	errors = Validate(validators)
	return
}

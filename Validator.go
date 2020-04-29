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

func getRule(rules, targetRule string) (string, bool) {
	//return !strings.Contains(rules, "|"+targetRule) && !strings.HasPrefix(rules, targetRule)
	for _, v := range strings.Split(rules, "|") {
		if strings.HasPrefix(v, targetRule) {
			return v, true
		}
	}
	return "", false
}
func getIntRuleValue(rule string) (int, error) {
	data := strings.Split(rule, ":")
	value := data[1]
	return strconv.Atoi(value)
}
func getStringRuleValue(rule string) string {
	data := strings.Split(rule, ":")
	value := data[1]
	return value
}
func setError(message string, rule Rule) error {
	if rule.CustomMessage != "" {
		return fmt.Errorf(rule.CustomMessage)
	}
	return fmt.Errorf(message)
}
func checkRegex(rule Rule, key, regexPattern string) bool {
	regex := regexp.MustCompile(regexPattern)
	if !regex.MatchString(fmt.Sprintf("%s", rule.Value)) {
		return false
	}
	return true
}

func unique(rule Rule, key string, errors *[]error) {
	uniqueRule, ok := getRule(rule.Rule, "unique")
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

	row := stmt.QueryRow(rule.Value)
	var result bool
	row.Scan(&result)

	if result == true {
		err = fmt.Errorf(invalidUnique, key)
		*errors = append(*errors, err)
		return
	}
}
func date(rule Rule, key string, errors *[]error) bool {
	_, ok := getRule(rule.Rule, "date")
	if !ok {
		return false
	}
	_, err := time.Parse(defaultDateFormat, fmt.Sprintf("%s", rule.Value))
	if err != nil {
		errMessage := fmt.Sprintf(invalidDate, key)
		*errors = append(*errors, setError(errMessage, rule))
	}
	return true
}
func datetime(rule Rule, key string, errors *[]error) bool {
	_, ok := getRule(rule.Rule, "datetime")
	if !ok {
		return false
	}
	_, err := time.Parse(defaultDateTimeFormat, fmt.Sprintf("%s", rule.Value))
	if err != nil {
		errMessage := fmt.Sprintf(invalidDateTime, key)
		*errors = append(*errors, setError(errMessage, rule))
	}
	return true
}
func timeFormat(rule Rule, key string, errors *[]error) bool {
	_, ok := getRule(rule.Rule, "time")
	if !ok {
		return false
	}
	_, err := time.Parse(defaultTimeFormat, fmt.Sprintf("%s", rule.Value))
	if err != nil {
		errMessage := fmt.Sprintf(invalidTime, key)
		*errors = append(*errors, setError(errMessage, rule))
	}
	return true
}
func numeric(rule Rule, key string, errors *[]error) {
	_, ok := getRule(rule.Rule, "numeric")
	if !ok {
		return
	}
	if checkRegex(rule, key, RegexNotNumeric) == true {
		errMessage := fmt.Sprintf(invalidNumeric, key)
		*errors = append(*errors, setError(errMessage, rule))
	}
}
func allowEmpty(rule Rule) bool {
	_, ok := getRule(rule.Rule, "allowempty")
	if !ok {
		return false
	}

	if rule.Value == "" || rule.Value == nil {
		return true
	}
	return false
}
func requiredIf(rule Rule, validators map[string]Rule, key string, errors *[]error) bool {
	requiredIfRule, ok := getRule(rule.Rule, "required_if:")
	if !ok {
		return false
	}

	otherField := getStringRuleValue(requiredIfRule)

	data := strings.Split(otherField, ",")
	targetKey := data[0]
	targetValue := data[1]

	targetValidator := validators[targetKey]

	if targetValidator.Value == targetValue {
		required(rule, key, errors)
	}
	return true
}
func min(rule Rule, key string, errors *[]error) {
	minRule, ok := getRule(rule.Rule, "min:")
	if !ok {
		return
	}

	minValue, err := getIntRuleValue(minRule)
	if err != nil {
		panic("min attribute value must be numeric")
		return
	}

	switch rule.Value.(type) {
	case int:
		if rule.Value.(int) < minValue {
			errMessage := fmt.Sprintf(invalidMin, key, minValue)
			*errors = append(*errors, setError(errMessage, rule))
		}
	case string:
		if len(fmt.Sprintf("%s", rule.Value)) < minValue {
			errMessage := fmt.Sprintf(invalidMinCharacters, key, minValue)
			*errors = append(*errors, setError(errMessage, rule))
		}
	}
}
func max(rule Rule, key string, errors *[]error) {
	maxRule, ok := getRule(rule.Rule, "max:")
	if !ok {
		return
	}
	maxValue, err := getIntRuleValue(maxRule)
	if err != nil {
		panic("max attribute value must be numeric")
		return
	}

	switch rule.Value.(type) {
	case int:
		if rule.Value.(int) > maxValue {
			errMessage := fmt.Sprintf(invalidMax, key, maxValue)
			*errors = append(*errors, setError(errMessage, rule))
		}
	case string:
		if len(fmt.Sprintf("%s", rule.Value)) > maxValue {
			errMessage := fmt.Sprintf(invalidMaxCharacters, key, maxValue)
			*errors = append(*errors, setError(errMessage, rule))
		}
	}
}
func email(rule Rule, key string, errors *[]error) {
	_, ok := getRule(rule.Rule, "email")
	if !ok {
		return
	}

	value := fmt.Sprintf("%s", rule.Value)

	if len(value) > 254 || checkRegex(rule, key, RegexEmail) == false {
		errMessage := fmt.Sprintf(invalidEmail, key)
		*errors = append(*errors, setError(errMessage, rule))
	}
	return
}
func in(rule Rule, key string, errors *[]error) {
	inRule, ok := getRule(rule.Rule, "in:")
	if !ok {
		return
	}

	data := strings.Split(getStringRuleValue(inRule), ",")
	for _, value := range data {
		if value == rule.Value {
			return
		}
	}
	errMessage := fmt.Sprintf(invalidIn, key)
	*errors = append(*errors, setError(errMessage, rule))
}
func required(rule Rule, key string, errors *[]error) {
	_, ok := getRule(rule.Rule, "required")
	if !ok {
		return
	}

	if rule.Value == "" || rule.Value == nil {
		errMessage := fmt.Sprintf(invalidRequired, key)
		*errors = append(*errors, setError(errMessage, rule))
	}
}
func calendar(rule Rule, key string, errors *[]error) {
	if datetime(rule, key, errors) == false {
		if date(rule, key, errors) == false {
			timeFormat(rule, key, errors)
		}
	}
}
func allowedEmptyOrDependsOtherField(rule Rule, validators map[string]Rule, key string, errors *[]error) bool {
	if allowEmpty(rule) == true {
		return true
	}
	if requiredIf(rule, validators, key, errors) == true {
		return true
	}
	return false
}
func startsWith(rule Rule, key string, errors *[]error) {
	startsWithRule, ok := getRule(rule.Rule, "starts_with")
	if !ok {
		return
	}

	targetPrefix := getStringRuleValue(startsWithRule)
	value := fmt.Sprintf("%s", rule.Value)

	if len(value) < len(targetPrefix) {
		errMessage := fmt.Sprintf(invalidStartsWith, key, targetPrefix)
		*errors = append(*errors, setError(errMessage, rule))
		return
	}

	if value[0:len(targetPrefix)-1] != targetPrefix {
		errMessage := fmt.Sprintf(invalidStartsWith, key, targetPrefix)
		*errors = append(*errors, setError(errMessage, rule))
	}
}
func endsWith(rule Rule, key string, errors *[]error) {
	endsWithRule, ok := getRule(rule.Rule, "ends_with")
	if !ok {
		return
	}

	targetPrefix := getStringRuleValue(endsWithRule)
	value := fmt.Sprintf("%s", rule.Value)

	if len(value) < len(targetPrefix) {
		errMessage := fmt.Sprintf(invalidEndsWith, key, targetPrefix)
		*errors = append(*errors, setError(errMessage, rule))
		return
	}

	if value[len(value)-len(targetPrefix):] != targetPrefix {
		errMessage := fmt.Sprintf(invalidEndsWith, key, targetPrefix)
		*errors = append(*errors, setError(errMessage, rule))
	}
}
func regex(rule Rule, key string, errors *[]error) {
	regexRule, ok := getRule(rule.Rule, "regex:")
	if !ok {
		return
	}

	targetRegex := getStringRuleValue(regexRule)

	if checkRegex(rule, key, targetRegex) == false {
		errMessage := fmt.Sprintf(invalidRegex, key)
		*errors = append(*errors, setError(errMessage, rule))
	}
}

func Validate(validators map[string]Rule) (errors []error) {
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

	validators := make(map[string]Rule)

	// loop through each struct attribute
	for i := 0; i < structParam.NumField(); i++ {
		currField := structParam.Field(i)
		currValue := valueParam.Field(i)

		key, ok := currField.Tag.Lookup("json")
		if !ok {
			key = structParam.Field(i).Name
		}

		value := currValue.Interface()
		rule, ok := currField.Tag.Lookup("rule")
		if !ok {
			continue
		}

		customMessage, ok := currField.Tag.Lookup("message")
		if !ok {
			customMessage = ""
		}

		validators[key] = Rule{Value: value, Rule: rule, CustomMessage: customMessage}
	}

	errors = Validate(validators)
	return
}

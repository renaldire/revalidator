package main

import (
	"fmt"
	Validator "revalidator"
	"time"
)

type User struct {
	Name        string `rule:"required|min:5|max:10|regex:^[A-Za-z]+$"`
	FullName    string `rule:"required"`
	Age         string `rule:"required|numeric|min:18"`
	Email       string `rule:"allowempty|email"`
	Birthday    string `rule:"required|date"`
	Gender      string `rule:"in:male,female"`
	PaymentType string `rule:"required|min:2|max:8|in:cc,debit,cash"`
	CCNumber    string `rule:"required_if:PaymentType,cc"`
	Address     string `rule:"ends_with:Street"`
	Site        string `rule:"allowempty|starts_with:http://"`
}

func main() {
	name := "john96"
	fullName := "john doe96"
	age := "103120381381907401701294112312313123131435435345346654"
	email := "asd"
	birthday := "1997-03-30asd"
	gender := "male"
	payment_type := "cc"
	cc_number := ""
	address := "Nut Farm Street"
	site := "renaldi.xyz"
	//username:="john"

	// Set Validator Environment
	// This only required once, only if you have to check unique value in database
	Validator.ConnectionString = "Your Connection String here"
	Validator.DbDriver = "postgres" // postgres || mysql
	// If ConnectionString is invalid, then it will shows panic

	start := time.Now()

	// Validator Usages
	validator := map[string]Validator.Rules{
		// fieldName : {value, rule},

		"full name": {fullName, "required|regex:"+Validator.RegexAlphabetWithSpace},
		"name":      {name, "required|min:5|max:10|regex:"+Validator.RegexAlphabet},
		//"username":{username,"required|unique:users,username"}, // unique format: unique:table,column
		"age":                {age, "required|numeric|min:18"},
		"gender":             {gender, "in:male,female"}, //in format: in:value1,value2,...valueN
		"address":            {address, "ends_with:Street"},
		"email":              {email, "allowempty|email"},
		"site":               {site, "allowempty|starts_with:http://"},
		"birthday":           {birthday, "required|date"},
		"payment type":       {payment_type, "required|min:2|max:8|in:cc,debit,cash"},
		"credit card number": {cc_number, "required_if:payment type,cc"}, //required_if format: required_if:desiredField,desiredValue
	}
	// Get All Errors based on defined rule
	errs := Validator.Validate(validator)

	execution_time := time.Since(start)
	fmt.Println("\nRun Time for Validate :", execution_time)

	if errs != nil {
		fmt.Printf("Found %d numbers of errors\n",len(errs))
		fmt.Println(errs)
	}

	// Validate certain struct
	// rule is defined in the struct tag

	start = time.Now()
	user := User{
		Name:        name,
		FullName:    fullName,
		Age:         age,
		Email:       email,
		Birthday:    birthday,
		Gender:      gender,
		PaymentType: payment_type,
		CCNumber:    cc_number,
		Address:     address,
		Site:        site,
	}

	errs = Validator.ValidateStruct(user)

	execution_time = time.Since(start)
	fmt.Println("\nRun Time for ValidateStruct :", execution_time)

	if errs != nil {
		fmt.Printf("Found %d numbers of errors\n",len(errs))
		fmt.Println(errs)
	}
	return
}

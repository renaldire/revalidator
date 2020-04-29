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
	Email       string `rule:"allowempty|email" message:"Emailnya salah"`
	Birthday    string `rule:"required|datetime"`
	Gender      string `rule:"in:male,female"`
	PaymentType string `rule:"required|min:2|max:8|in:cc,debit,cash"`
	CCNumber    string `rule:"required_if:PaymentType,cc"`
	Address     string `rule:"ends_with:Street"`
	Site        string `rule:"allowempty|starts_with:http://"`
}

func main() {
	name := "john96"
	fullName := "john doe96"
	age := "1"
	email := "asd"
	birthday := "1997-03-30"
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
	validator := map[string]Validator.Rule{
		// fieldName : {value, rule},

		"full name": {Value: fullName, Rule: "required|regex:" + Validator.RegexAlphabetWithSpace},
		"name":      {Value: name, Rule: "required|min:5|max:10|regex:" + Validator.RegexAlphabet},
		//"username":{username,"required|unique:users,username"}, // unique format: unique:table,column
		"age":                {Value: age, Rule: "required|numeric|min:18", CustomMessage: "Belum Cukup Umur"},
		"gender":             {Value: gender, Rule: "in:male,female"}, //in format: in:value1,value2,...valueN
		"address":            {Value: address, Rule: "ends_with:Street"},
		"email":              {Value: email, Rule: "allowempty|email", CustomMessage: "Emailnya Salah"},
		"site":               {Value: site, Rule: "allowempty|starts_with:http://"},
		"birthday":           {Value: birthday, Rule: "required|datetime"},
		"payment type":       {Value: payment_type, Rule: "required|min:2|max:8|in:cc,debit,cash"},
		"credit card number": {Value: cc_number, Rule: "required_if:payment type,cc"}, //required_if format: required_if:desiredField,desiredValue
	}
	// Get All Errors based on defined rule
	errs := Validator.Validate(validator)

	execution_time := time.Since(start)
	fmt.Println("\nRun Time for Validate :", execution_time)

	if errs != nil {
		fmt.Printf("Found %d numbers of errors\n", len(errs))
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
		fmt.Printf("Found %d numbers of errors\n", len(errs))
		fmt.Println(errs)
	}
	return
}

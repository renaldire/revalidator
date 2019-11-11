package main

import (
	"fmt"
	Validator "revalidator"
	"time"
)

func main(){
	name:=""
	age:=10
	email:="asd"
	birthday:="1997-03-30asd"
	gender:=""
	payment_type:="cc"
	cc_number:=""
	username:="john"

	// Set Validator Environtment
	// This only required once, only if you have to check unique value in database
	Validator.ConnectionString="Your Connection String here"
	Validator.DbDriver="postgres" // postgres || mysql

	start:=time.Now()

	// Validator Usages
	validator:=map[string]Validator.Rules{
		"name":{name,"required|min:5|max:10"},
		"username":{username,"required|unique:users,username"}, // unique format: unique:table,column
		"age":{age,"required|numeric|min:18"},
		"gender":{gender,"in:male,female"}, //in format: in:value1,value2,...valueN
		"email":{email,"allowempty|email"},
		"birthday":{birthday,"required|date"},
		"payment type":{payment_type,"required|min:2|max:8|in:cc,debit,cash"},
		"credit card number":{cc_number,"required_if:payment type,cc"}, //required_if format: required_if:desiredField,desiredValue
	}
	// Get All Errors based on defined rule
	errs := Validator.Validate(validator)

	end:=time.Now()
	execution_time:=end.Sub(start)
	fmt.Println("Run Time :",execution_time)
	if errs !=nil {
		fmt.Println(errs)
		return
	}
	fmt.Println("Success")
	return
}



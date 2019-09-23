package main

import (
	"fmt"
	Validator "revalidator"
)

func main(){
	name:="john doe"
	age:=19
	email:="john@mail.com"
	birthday:="1997-03-30"
	gender:="male"
	payment_type:="test"
	cc_number:=""
	validator:=map[string]Validator.Rules{
		"name":{name,"required|min:5|max:10"},
		"age":{age,"required|numeric|min:18"},
		"gender":{gender,"required|in:male,female"},
		"email":{email,"allowempty|email"},
		"birthday":{birthday,"required|date"},
		"payment type":{payment_type,"required|in:cc,debit,cash"},
		"credit card number":{cc_number,"required_if:payment type,cc"},
	}
	err:= Validator.Validate(validator)

	if err!=nil {
		fmt.Println(err[0].Error())
		return
	}
	fmt.Println("Success")
	return
}



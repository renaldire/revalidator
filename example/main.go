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
	payment_type:="cc"
	cc_number:=""
	validator:=map[string]Validator.Rules{
		"name":{name,"required|min:5|max:10"},
		"age":{age,"required|numeric|min:18"},
		"email":{email,"allowempty|email"},
		"birthday":{birthday,"required|date"},
		"payment type":{payment_type,"required"},
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



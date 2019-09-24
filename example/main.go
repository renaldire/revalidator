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
	payment_type:=""
	cc_number:=""

	start:=time.Now()
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
	end:=time.Now()

	execution_time:=end.Sub(start)
	fmt.Println("Run Time :",execution_time)
	if err!=nil {
		fmt.Println(err[0].Error())
		return
	}
	fmt.Println("Success")
	return
}



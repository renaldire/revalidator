package main

import (
	"Validator"
	"fmt"
)

func main(){
	name:="john"
	age:="asdada"
	email:="renaldi"
	birthday:="1997-03-30adssad"
	validator:=map[string]Validator.Rules{
		"name":{name,"required|min:5|max:10"},
		"age":{age,"required|numeric|min:18"},
		"email":{email,"allowempty|email"},
		"birthday":{birthday,"required|date"},
	}

	err:= Validator.Validate(validator)
	if err!=nil {
		fmt.Println(err[0].Error())
		return
	}
	fmt.Println("Success")
	return
}



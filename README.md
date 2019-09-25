# revalidator

Validator for golang.

## Instalation

    go get "github.com/renaldire/revalidator"

## Usage

    validator:=map[string]Validator.Rules{  
	   "<Field Name>":{<Value>,"<Rule 1>|<Rule 2>|...|<Rule n>"} 
	}
	
	errs:=revalidator.Validate(validator)
	if errs!=nil{
		return errs[0] // retrieve the first error
	}

## Example

    
    import Validator "revalidator"
    
    name:=""  
    age:=10  
    email:="asd"  
    birthday:="1997-03-30asd"  
    gender:=""  
    payment_type:=""  
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
    errs:= Validator.Validate(validator) 
     
    if errs!=nil {  
       fmt.Println(errs)  
       return  
    }

Result:

    [name is required. 
    name must be at least 5 characters 
    age must be not less than 18 
    gender is required. 
    gender is invalid. 
    email must be in e-mail format. 
    birthday must be in yyyy-mm-dd format. 
    payment type is required. 
    payment type is invalid.]

## Rules

 - allowempty
 - required_if
 - required
 - min
 - max
 - numeric
 - email
 - date
 - time
 - datetime
 - in

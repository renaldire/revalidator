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
	
## What's new
    Jan 30, 2020
    ============
    - Starts With Rule
    - Ends With Rule
    - Regex Rule

    Nov 11,2019
    ===========
    - Unique rule 
	
	To use this rule, you have to set Validator environtment
	
	Validator.ConnectionString="Your Connection String here"
    Validator.DbDriver="postgres" // postgres || mysql
	
    This only required once, only if you have to check unique value in database

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

|     Rule    |                                                              Description                                                              |
|:-----------:|:-------------------------------------------------------------------------------------------------------------------------------------:|
| allowempty  | allow field to be empty, and once the field is not empty, allow other rules to be apply                                               |
| required    | field must be filled                                                                                                                  |
| required_if | field is required but depends on other field                                                                                          |
| min         | minimum value for numeric or minimum length for a string                                                                              |
| max         | maximum value for numeric or maximum length for a string                                                                              |
| numeric     | a string can only contains numeric characters                                                                                         |
| email       | a string must be a valid email format                                                                                                 |
| date        | yyyy-mm-dd                                                                                                                            |
| datetime    | yyyy-mm-dd hh:mm:ss                                                                                                                   |
| time        | hh:mm:ss                                                                                                                              |
| in          | field only valid if the value is one of the defined value in rule.   Example: Gender can only be male or female. Rule: in:male,female |
| unique      | field must be unique in a database. Example: Validate if the username is already taken or not.          
| starts_with      | field must be starts with a certain string  |
| ends_with      | field must be ends with a certain string  |
| regex      | field must be matches with a certain regex format  |

See example package for more detail.
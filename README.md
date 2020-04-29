# revalidator

Validator for golang.

## Instalation

    go get "github.com/renaldire/revalidator"

## Usage

    validator:=map[string]Validator.Rules{  
	   "<Field Name>":{<Value>,"<Rule 1>|<Rule 2>|...|<Rule n>",<Optional Custom Error Message>} 
	}
	
	errs:=revalidator.Validate(validator)
	if errs!=nil{
		return errs[0] // retrieve the first error
	}
	
## What's new
    Apr 29, 2020
    ============
    - Custom Error Message
    "full name": {Value: "", Rule: "required", CustomMessage:"Nama harus diisi"}
    
    type User struct {
        Name        string `rule:"required"`
        FullName    string `rule:"required" message:"Nama Lengkap Wajib Diisi"`
    }
    
    Feb 03, 2020
    ============
    - Bug Fixes
    - Validate using Struct
    
    type User struct {
    	Name        string `rule:"required|min:5|max:10|regex:^[A-Za-z]+$"`
    	FullName    string `rule:"required"`
    }
    ...
    errs = Validator.ValidateStruct(user)
    ...
    // More Example already provided in example package
    
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
       "name":{Value: name, Rule: "required|min:5|max:10"},  
       "age":{Value: age, Rule: "required|numeric|min:18"},  
       "gender":{Value: gender, Rule: "required|in:male,female"},  
       "email":{Value: email, Rule: "allowempty|email", CustomMessage:"Emailnya Salah"},  
       "birthday":{Value: birthday, Rule: "required|date"},  
       "payment type":{Value: payment_type, Rule: "required|in:cc,debit,cash"},  
       "credit card number":{Value: cc_number, Rule: "required_if:payment type,cc"},  
    }  
    errs:= Validator.Validate(validator) 
     
    if errs!=nil {  
       fmt.Println(errs)  
       return  
    }

    Also see package example for more advanced example
Result:

    [name is required. 
    name must be at least 5 characters 
    age must be not less than 18 
    gender is required. 
    gender is invalid. 
    Emailnya Salah. 
    birthday must be in yyyy-mm-dd format. 
    payment type is required. 
    payment type is invalid.]

## Rules

Notes:
Every value will be trimmed before being validate.

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
package Validator

const (
	regexEmail = "^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"

	log                  = "[500] Internal Server Error. %s attribute must be a numeric int.\n"
	invalidMin           = "%s must be not less than %d"
	invalidMinCharacters = "%s must be at least %d characters"

	invalidMax           = "%s must be not more than %d"
	invalidMaxCharacters = "%s must be not more than %d characters"

	invalidRequired       = "%s is required."
	invalidEmail          = "%s must be in e-mail format."
	invalidIn             = "%s is invalid."
	invalidNumeric        = "%s must be a numeric"
	invalidUnique         = "%s is already exists."
	invalidDate           = "%s must be in yyyy-mm-dd format."
	invalidDateTime       = "%s must be in yyyy-MM-dd hh:mm:ss format."
	invalidTime           = "%s must be in hh:mm:ss format."
	defaultDateFormat     = "2006-01-02"
	defaultDateTimeFormat = "2006-01-02T15:04:05"
	defaultTimeFormat     = "15:04:05"
)

package Validator

const (
	RegexEmail                 = "^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
	RegexNotNumeric            = "\\D"
	RegexAlphabet              = "^[A-Za-z]+$"
	RegexAlphabetWithSpace     = "^[A-Za-z ]+$"
	RegexAlphanumeric          = "^[A-Za-z0-9]+$"
	RegexAlphanumericWithSpace = "^[A-Za-z0-9 ]+$"

	invalidMin           = "%s must be not less than %d"
	invalidMinCharacters = "%s must be at least %d characters"

	invalidMax           = "%s must be not more than %d"
	invalidMaxCharacters = "%s must be not more than %d characters"

	invalidStartsWith = "%s must be starts with '%s'"
	invalidEndsWith   = "%s must be ends with '%s'"
	invalidRegex      = "invalid %s format"

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

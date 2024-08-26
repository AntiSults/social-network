package security

import "regexp"

// IsValidEmail checking whether entered e-mail valid
func IsValidEmail(email string) bool {

	// Define the regular expression for email validation
	rePattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	//more complicated re pattern
	//rePattern := "^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"

	regexpObj := regexp.MustCompile(rePattern)

	return regexpObj.MatchString(email)
}

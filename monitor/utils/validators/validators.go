package validators

import (
	"errors"
	"net/url"
)

var validationErrorMap = map[string]error{
	"InvalidLongURL": errors.New("Error, the provided long url is not valid: %v \n The URL must have a Scheme, Host, and/or a Path"),
}

// IsValidURL checks if the URL has a Scheme, Host, and/or a Path
func IsValidURL(str string) (bool, error) {
	u, err := url.Parse(str)
	if err == nil && u.Scheme != "" && u.Host != "" {
		return true, nil
	}
	return false, validationErrorMap["InvalidLongURL"]
}

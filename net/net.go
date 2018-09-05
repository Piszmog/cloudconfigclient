package net

import (
	"strings"
)

// CreateUrl creates a full URL from the specified base URL and the array of URI variables.
//
// URI variables are separated by '/'.
func CreateUrl(baseUrl string, uriVariables ...string) string {
	url := strings.TrimRight(baseUrl, "/")
	for _, uriVariable := range uriVariables {
		url = url + "/" + uriVariable
	}
	return url
}

// JoinProfiles joins the array of profiles with a comma.
func JoinProfiles(profiles []string) string {
	return strings.Join(profiles, ",")
}

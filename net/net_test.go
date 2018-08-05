package net

import "testing"

func TestCreateUrl(t *testing.T) {
	baseUrl := "http://localhost:8080"
	url := CreateUrl(baseUrl, "variable1", "variable2", "variable3")
	if url != "http://localhost:8080/variable1/variable2/variable3" {
		t.Errorf("constructed url does not match excepted url")
	}
}

func TestCreateUrlWhenBaseUrlEndsWithSlash(t *testing.T) {
	baseUrl := "http://localhost:8080/"
	url := CreateUrl(baseUrl, "variable1", "variable2", "variable3")
	if url != "http://localhost:8080/variable1/variable2/variable3" {
		t.Errorf("constructed url does not match excepted url")
	}
}

func TestJoinProfiles(t *testing.T) {
	profiles := []string{"profile1", "profile2"}
	profilesString := JoinProfiles(profiles)
	if profilesString != "profile1,profile2" {
		t.Errorf("profiles were not append")
	}
}

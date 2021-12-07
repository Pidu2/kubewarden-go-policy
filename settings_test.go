package main

import (
	"testing"
)

func TestParsingSettingsWithAllValuesProvidedFromValidationReq(t *testing.T) {
	request := `
	{
		"request": "doesn't matter here",
		"settings": {
			"annotation": [ "foo", "bar" ]
		}
	}
	`
	rawRequest := []byte(request)

	settings, err := NewSettingsFromValidationReq(rawRequest)
	if err != nil {
		t.Errorf("Unexpected error %+v", err)
	}

	expected := []string{"foo", "bar"}
	for _, exp := range expected {
		if !settings.AllowNodePortAnnotations.Contains(exp) {
			t.Errorf("Missing value %s", exp)
		}
	}
}

func TestParsingSettingsWithNoValueProvided(t *testing.T) {
	request := `
	{
		"request": "doesn't matter here",
		"settings": {
		}
	}
	`
	rawRequest := []byte(request)

	settings, err := NewSettingsFromValidationReq(rawRequest)
	if err != nil {
		t.Errorf("Unexpected error %+v", err)
	}

	if settings.AllowNodePortAnnotations.Cardinality() != 0 {
		t.Errorf("Expecpted AllowNodePortAnnotations to be empty")
	}
}

func TestSettingsAreValid(t *testing.T) {
	request := `
	{
    "request": "doesnt matter",
    "settings": {
      "annotation": ["foo", "bar"]
    }
	}
	`
	rawRequest := []byte(request)

	settings, err := NewSettingsFromValidateSettingsPayload(rawRequest)
	if err != nil {
		t.Errorf("Unexpected error %+v", err)
	}

	valid, err := settings.Valid()
	if !valid {
		t.Errorf("Settings are reported as not valid")
	}
	if err != nil {
		t.Errorf("Unexpected error %+v", err)
	}
}

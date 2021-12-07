package main

import (
	mapset "github.com/deckarep/golang-set"
	"github.com/kubewarden/gjson"
	kubewarden "github.com/kubewarden/policy-sdk-go"

	"fmt"
)

type Settings struct {
	AllowNodePortAnnotations mapset.Set `json:"annotation"`
}

// Builds a new Settings instance starting from a validation
// request payload:
// {
//    "request": ...,
//    "settings": {
//       "denied_names": [...]
//    }
// }
func NewSettingsFromValidationReq(payload []byte) (Settings, error) {
	return newSettings(
		payload,
		"settings.annotation")
}

// Builds a new Settings instance starting from a Settings
// payload:
// {
//    "denied_names": ...
// }
func NewSettingsFromValidateSettingsPayload(payload []byte) (Settings, error) {
	return newSettings(
		payload,
		"annotation")
}

func newSettings(payload []byte, paths ...string) (Settings, error) {
	if len(paths) != 1 {
		return Settings{}, fmt.Errorf("wrong number of json paths")
	}

	data := gjson.GetManyBytes(payload, paths...)

	allowNodePortAnnotation := mapset.NewThreadUnsafeSet()
	data[0].ForEach(func(_, entry gjson.Result) bool {
		allowNodePortAnnotation.Add(entry.String())
		return true
	})

	return Settings{
		AllowNodePortAnnotations: allowNodePortAnnotation,
	}, nil
}

// No special check has to be done
func (s *Settings) Valid() (bool, error) {
	return true, nil
}

func validateSettings(payload []byte) ([]byte, error) {
	settings, err := NewSettingsFromValidateSettingsPayload(payload)
	if err != nil {
		return kubewarden.RejectSettings(kubewarden.Message(err.Error()))
	}

	valid, err := settings.Valid()
	if err != nil {
		return kubewarden.RejectSettings(kubewarden.Message(fmt.Sprintf("Provided settings are not valid: %v", err)))
	}
	if valid {
		return kubewarden.AcceptSettings()
	}

	return kubewarden.RejectSettings(kubewarden.Message("Provided settings are not valid"))
}

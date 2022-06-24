package main

import (
	"testing"

	"github.com/mailru/easyjson"
)

func TestParsingSettingsWithNoValueProvided(t *testing.T) {
	rawSettings := []byte(`{}`)
	settings := &Settings{}
	if err := easyjson.Unmarshal(rawSettings, settings); err != nil {
		t.Errorf("Unexpected error %+v", err)
	}

	if len(settings.ExcludedNamespaces) != 0 {
		t.Errorf("Expected ExcludedNamespaces to be empty")
	}

	valid, err := settings.Valid()
	if !valid {
		t.Errorf("Settings are reported as not valid")
	}
	if err != nil {
		t.Errorf("Unexpected error %+v", err)
	}
}

func TestIsNameDenied(t *testing.T) {
	// settings := Settings{
	// 	ExcludedNamespaces: []string{"bob"},
	// }

	// // if !settings.IsNamespaceIstioDisabled("bob") {
	// // 	t.Errorf("name should be denied")
	// // }

	// // if settings.IsNameDenied("alice") {
	// // 	t.Errorf("name should not be denied")
	// // }
}

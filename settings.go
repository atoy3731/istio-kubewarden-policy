package main

import (
	kubewarden "github.com/kubewarden/policy-sdk-go"
	kubewarden_protocol "github.com/kubewarden/policy-sdk-go/protocol"
	"github.com/mailru/easyjson"

	"fmt"
)

// The Settings class is defined inside of the `types.go` file

// No special checks have to be done
func (s *Settings) Valid() (bool, error) {
	return true, nil
}

func (s *Settings) IsNamespaceIstioDisabled(name string, annotations map[string]string) bool {
	for _, excludedNamespace := range s.ExcludedNamespaces {
		if excludedNamespace == name {
			return false
		}
	}

	for k, v := range annotations {
		if k == "istio-injection" && v == "true" {
			return false
		}
	}

	return true
}

func (s *Settings) IsPodIstioDisabled(labels map[string]string, annotations map[string]string) bool {
	for labelKey, labelValue := range s.ExcludedPodLabels {
		if labels[labelKey] == labelValue {
			return false
		}
	}

	for k, v := range annotations {
		if k == "sidecar.istio.io/inject" && v == "true" {
			return true
		}
	}

	return false
}

func NewSettingsFromValidationReq(validationReq *kubewarden_protocol.ValidationRequest) (Settings, error) {
	settings := Settings{}
	err := easyjson.Unmarshal(validationReq.Settings, &settings)
	return settings, err
}

func validateSettings(payload []byte) ([]byte, error) {
	logger.Info("validating settings")

	settings := Settings{}
	err := easyjson.Unmarshal(payload, &settings)
	if err != nil {
		return kubewarden.RejectSettings(kubewarden.Message(fmt.Sprintf("Provided settings are not valid: %v", err)))
	}

	valid, err := settings.Valid()
	if err != nil {
		return kubewarden.RejectSettings(kubewarden.Message(fmt.Sprintf("Provided settings are not valid: %v", err)))
	}
	if valid {
		return kubewarden.AcceptSettings()
	}

	logger.Warn("rejecting settings")
	return kubewarden.RejectSettings(kubewarden.Message("Provided settings are not valid"))
}

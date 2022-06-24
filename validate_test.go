package main

import (
	"testing"

	kubewarden_protocol "github.com/kubewarden/policy-sdk-go/protocol"
	kubewarden_testing "github.com/kubewarden/policy-sdk-go/testing"
	"github.com/mailru/easyjson"
)

// func TestEmptySettingsLeadsToApproval(t *testing.T) {
// 	settings := Settings{}
// 	pod := corev1.Pod{
// 		Metadata: metav1.ObjectMeta{
// 			Name:      "test-pod",
// 			Namespace: "default",
// 		},
// 	}

// 	payload, err := kubewarden_testing.BuildValidationRequest(&pod, &settings)
// 	if err != nil {
// 		t.Errorf("Unexpected error: %+v", err)
// 	}

// 	responsePayload, err := validate(payload)
// 	if err != nil {
// 		t.Errorf("Unexpected error: %+v", err)
// 	}

// 	var response kubewarden_protocol.ValidationResponse
// 	if err := easyjson.Unmarshal(responsePayload, &response); err != nil {
// 		t.Errorf("Unexpected error: %+v", err)
// 	}

// 	if response.Accepted != true {
// 		t.Errorf("Unexpected rejection: msg %s - code %d", *response.Message, *response.Code)
// 	}
// }

// func TestApproval(t *testing.T) {
// 	settings := Settings{
// 		DeniedNames: []string{"foo", "bar"},
// 	}
// 	pod := corev1.Pod{
// 		Metadata: metav1.ObjectMeta{
// 			Name:      "test-pod",
// 			Namespace: "default",
// 		},
// 	}

// 	payload, err := kubewarden_testing.BuildValidationRequest(&pod, &settings)
// 	if err != nil {
// 		t.Errorf("Unexpected error: %+v", err)
// 	}

// 	responsePayload, err := validate(payload)
// 	if err != nil {
// 		t.Errorf("Unexpected error: %+v", err)
// 	}

// 	var response kubewarden_protocol.ValidationResponse
// 	if err := easyjson.Unmarshal(responsePayload, &response); err != nil {
// 		t.Errorf("Unexpected error: %+v", err)
// 	}

// 	if response.Accepted != true {
// 		t.Error("Unexpected rejection")
// 	}
// }

// func TestApproveFixture(t *testing.T) {
// 	settings := Settings{
// 		DeniedNames: []string{},
// 	}

// 	payload, err := kubewarden_testing.BuildValidationRequestFromFixture(
// 		"test_data/pod.json",
// 		&settings)
// 	if err != nil {
// 		t.Errorf("Unexpected error: %+v", err)
// 	}

// 	responsePayload, err := validate(payload)
// 	if err != nil {
// 		t.Errorf("Unexpected error: %+v", err)
// 	}

// 	var response kubewarden_protocol.ValidationResponse
// 	if err := easyjson.Unmarshal(responsePayload, &response); err != nil {
// 		t.Errorf("Unexpected error: %+v", err)
// 	}

// 	if response.Accepted != true {
// 		t.Error("Unexpected rejection")
// 	}
// }

func TestAcceptanceBecauseNamespaceIsExcluded(t *testing.T) {
	settings := Settings{
		ExcludedNamespaces: []string{"foo", "bar"},
	}

	payload, err := kubewarden_testing.BuildValidationRequestFromFixture("./test_data/namespace-enabled.json", &settings)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	responsePayload, err := validate(payload)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	var response kubewarden_protocol.ValidationResponse
	if err := easyjson.Unmarshal(responsePayload, &response); err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if response.Accepted == false {
		t.Error("Unexpected denial")
	}
}

func TestRejectionBecauseNamespaceIsNotInjected(t *testing.T) {
	settings := Settings{
		ExcludedNamespaces: []string{"unmatched"},
	}

	payload, err := kubewarden_testing.BuildValidationRequestFromFixture("./test_data/namespace-disabled.json", &settings)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	responsePayload, err := validate(payload)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	var response kubewarden_protocol.ValidationResponse
	if err := easyjson.Unmarshal(responsePayload, &response); err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if response.Accepted == true {
		t.Error("Unexpected acceptance")
	}
}

func TestRejectionBecausePodIsIstioDiabled(t *testing.T) {
	settings := Settings{
		ExcludedNamespaces: []string{"unmatched"},
		ExcludedPodLabels: map[string]string{
			"istioException": "true",
		},
	}

	payload, err := kubewarden_testing.BuildValidationRequestFromFixture("./test_data/pod-disabled.json", &settings)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	responsePayload, err := validate(payload)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	var response kubewarden_protocol.ValidationResponse
	if err := easyjson.Unmarshal(responsePayload, &response); err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if response.Accepted == true {
		t.Error("Unexpected acceptance")
	}
}

func TestAcceptanceBecausePodIsExcluded(t *testing.T) {
	settings := Settings{
		ExcludedNamespaces: []string{"unmatched"},
		ExcludedPodLabels: map[string]string{
			"istioException": "enabled",
		},
	}

	payload, err := kubewarden_testing.BuildValidationRequestFromFixture("./test_data/pod-disabled.json", &settings)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	responsePayload, err := validate(payload)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	var response kubewarden_protocol.ValidationResponse
	if err := easyjson.Unmarshal(responsePayload, &response); err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if response.Accepted != true {
		t.Error("Unexpected denial")
	}
}

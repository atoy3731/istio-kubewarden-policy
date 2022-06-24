package main

import (
	"fmt"

	onelog "github.com/francoispqt/onelog"
	corev1 "github.com/kubewarden/k8s-objects/api/core/v1"
	kubewarden "github.com/kubewarden/policy-sdk-go"
	kubewarden_protocol "github.com/kubewarden/policy-sdk-go/protocol"
	"github.com/mailru/easyjson"
)

func checkNamespace(settings Settings, validationRequest kubewarden_protocol.ValidationRequest) ([]byte, error) {
	// Access the **raw** JSON that describes the object
	namespaceJSON := validationRequest.Request.Object

	// Try to create a Pod instance using the RAW JSON we got from the
	// ValidationRequest.
	namespace := &corev1.Namespace{}
	if err := easyjson.Unmarshal([]byte(namespaceJSON), namespace); err != nil {
		return kubewarden.RejectRequest(
			kubewarden.Message(
				fmt.Sprintf("Cannot decode Namespace object: %s", err.Error())),
			kubewarden.Code(400))
	}

	if settings.IsNamespaceIstioDisabled(namespace.Metadata.Name, namespace.Metadata.Labels) {
		logger.InfoWithFields("rejecting namespace object", func(e onelog.Entry) {
			e.String("name", namespace.Metadata.Name)
		})

		return kubewarden.RejectRequest(
			kubewarden.Message(
				fmt.Sprintf("The '%s' namespace is not Istio enabled", namespace.Metadata.Name)),
			kubewarden.NoCode)
	}

	return kubewarden.AcceptRequest()
}

func checkPod(settings Settings, validationRequest kubewarden_protocol.ValidationRequest) ([]byte, error) {
	// Access the **raw** JSON that describes the object
	podJSON := validationRequest.Request.Object

	// Try to create a Pod instance using the RAW JSON we got from the
	// ValidationRequest.
	pod := &corev1.Pod{}
	if err := easyjson.Unmarshal([]byte(podJSON), pod); err != nil {
		return kubewarden.RejectRequest(
			kubewarden.Message(
				fmt.Sprintf("Cannot decode Pod object: %s", err.Error())),
			kubewarden.Code(400))
	}

	if settings.IsPodIstioDisabled(pod.Metadata.Labels, pod.Metadata.Annotations) {
		logger.InfoWithFields("rejecting pod object", func(e onelog.Entry) {
			e.String("name", pod.Metadata.Name)
		})

		return kubewarden.RejectRequest(
			kubewarden.Message(
				fmt.Sprintf("The '%s' pod is not Istio enabled", pod.Metadata.Name)),
			kubewarden.NoCode)
	}

	return kubewarden.AcceptRequest()
}

func validate(payload []byte) ([]byte, error) {
	// Create a ValidationRequest instance from the incoming payload
	validationRequest := kubewarden_protocol.ValidationRequest{}
	err := easyjson.Unmarshal(payload, &validationRequest)
	if err != nil {
		return kubewarden.RejectRequest(
			kubewarden.Message(err.Error()),
			kubewarden.Code(400))
	}

	// Create a Settings instance from the ValidationRequest object
	settings, err := NewSettingsFromValidationReq(&validationRequest)
	if err != nil {
		return kubewarden.RejectRequest(
			kubewarden.Message(err.Error()),
			kubewarden.Code(400))
	}

	switch validationRequest.Request.RequestKind.Kind {
	case "Namespace":
		return checkNamespace(settings, validationRequest)
	case "Pod":
		return checkPod(settings, validationRequest)
	default:
		return kubewarden.AcceptRequest()
	}
}

package handlers

import (
	"encoding/json"
	"fmt"
)

type RequestBodyValidator struct{}

func NewRequestBodyValidator() RequestBodyValidator {
	return RequestBodyValidator{}
}

func (r RequestBodyValidator) ValidateRequestBody(allowedFields map[string]any, requestBody []byte) error {
	if len(requestBody) == 0 {
		return fmt.Errorf("request body is empty")
	}
	var bodyMap map[string]interface{}
	if err := json.Unmarshal(requestBody, &bodyMap); err != nil {
		return fmt.Errorf("failed to parse request body: %w", err)
	}

	if len(bodyMap) == 0 {
		return fmt.Errorf("request body is empty")
	}

	for key := range bodyMap {
		if _, ok := allowedFields[key]; !ok {
			return fmt.Errorf("invalid field '%s' in request body", key)
		}
	}

	return nil
}

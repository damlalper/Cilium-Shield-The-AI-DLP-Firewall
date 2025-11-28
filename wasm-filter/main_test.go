
package main

import (
	"testing"
)

func TestIsValidLuhn(t *testing.T) {
	testCases := []struct {
		name     string
		number   string
		expected bool
	}{
		{"valid", "49927398716", true},
		{"valid with spaces", "4992 7398 716", true},
		{"valid with dashes", "4992-7398-716", true},
		{"invalid", "49927398717", false},
		{"invalid short", "123", false},
		{"empty", "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isValidLuhn(tc.number); got != tc.expected {
				t.Errorf("isValidLuhn(%q) = %v; want %v", tc.number, got, tc.expected)
			}
		})
	}
}

func TestRedactSensitiveData(t *testing.T) {
	testCases := []struct {
		name           string
		body           string
		expectedBody   string
		expectRedacted bool
	}{
		{
			name:           "no sensitive data",
			body:           `{"message": "hello world"}`,
			expectedBody:   `{"message": "hello world"}`,
			expectRedacted: false,
		},
		{
			name:           "valid credit card",
			body:           `{"card": "49927398716", "user": "test"}`,
			expectedBody:   `{"card": "[REDACTED_CREDIT_CARD]", "user": "test"}`,
			expectRedacted: true,
		},
		{
			name:           "invalid credit card",
			body:           `{"card": "49927398717", "user": "test"}`,
			expectedBody:   `{"card": "49927398717", "user": "test"}`,
			expectRedacted: false,
		},
		{
			name:           "openai api key",
			body:           `{"api_key": "sk-proj-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"}`,
			expectedBody:   `{"api_key": "[REDACTED_API_KEY]"}`,
			expectRedacted: true,
		},
		{
			name:           "google api key",
			body:           `{"api_key": "AIzaSyxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"}`,
			expectedBody:   `{"api_key": "[REDACTED_API_KEY]"}`,
			expectRedacted: true,
		},
		{
			name:           "email address",
			body:           `{"email": "test@example.com"}`,
			expectedBody:   `{"email": "[REDACTED_EMAIL]"}`,
			expectRedacted: true,
		},
		{
			name:           "multiple redactions",
			body:           `{"email": "test@example.com", "card": "49927398716", "key": "sk-proj-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"}`,
			expectedBody:   `{"email": "[REDACTED_EMAIL]", "card": "[REDACTED_CREDIT_CARD]", "key": "[REDACTED_API_KEY]"}`,
			expectRedacted: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			redactedBody, redacted := redactSensitiveData(tc.body)
			if redacted != tc.expectRedacted {
				t.Errorf("Expected redacted to be %v, but got %v", tc.expectRedacted, redacted)
			}
			if redactedBody != tc.expectedBody {
				t.Errorf("Expected body to be %q, but got %q", tc.expectedBody, redactedBody)
			}
		})
	}
}

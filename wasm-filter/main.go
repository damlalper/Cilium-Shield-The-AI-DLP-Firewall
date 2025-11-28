
package main

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

// main is the entry point for the Wasm module.
// It sets the VM context, which is required for the Wasm module to be executed by the host.
func main() {
	proxywasm.SetVMContext(&vmContext{})
}

type vmContext struct {
	// Embed the default VM context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultVMContext
}

// NewPluginContext creates a new plugin context.
// This is called once when the plugin is loaded.
func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &pluginContext{}
}

type pluginContext struct {
	// Embed the default plugin context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultPluginContext
}

// NewHttpContext creates a new HTTP context.
// This is called for each request that the filter processes.
func (*pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	return &httpContext{}
}

type httpContext struct {
	// Embed the default HTTP context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultHttpContext
}

// OnHttpRequestBody is called when the HTTP request body is received.
// This is where we will inspect the request body for sensitive data.
func (ctx *httpContext) OnHttpRequestBody(bodySize int, endOfStream bool) types.Action {
	// We need the full request body to perform the redaction.
	// If the body is streamed, we need to wait until the full body is received.
	if !endOfStream {
		return types.ActionPause
	}

	// Get the entire request body.
	// Note on memory management: GetHttpRequestBody allocates memory in the Wasm module
	// to hold the request body. This memory is managed by the Wasm runtime and will be
	// garbage collected. However, for very large request bodies, this could lead to
	// high memory consumption. In a production environment, we might want to
	// implement a streaming approach to process the body in chunks.
	originalBody, err := proxywasm.GetHttpRequestBody(0, bodySize)
	if err != nil {
		proxywasm.LogErrorf("failed to get request body: %v", err)
		return types.ActionContinue
	}

	// Redact sensitive data from the request body.
	redactedBody, redacted := redactSensitiveData(string(originalBody))

	if redacted {
		// If data was redacted, replace the request body with the redacted version.
		// Note on memory management: ReplaceHttpRequestBody will copy the redacted body
		// into the host's memory. The memory allocated for the original body in the
		// Wasm module will be garbage collected.
		err := proxywasm.ReplaceHttpRequestBody([]byte(redactedBody))
		if err != nil {
			proxywasm.LogErrorf("failed to replace request body: %v", err)
			return types.ActionContinue
		}
		// Add a header to indicate that the request has been redacted.
		err = proxywasm.AddHttpRequestHeader("X-Cilium-Shield-Status", "REDACTED")
		if err != nil {
			proxywasm.LogErrorf("failed to add request header: %v", err)
			return types.ActionContinue
		}
	}

	return types.ActionContinue
}

// redactSensitiveData redacts sensitive data from the given string.
func redactSensitiveData(body string) (string, bool) {
	redacted := false

	// Redact credit card numbers using a regex and a Luhn algorithm check.
	// The regex is designed to be efficient and avoid backtracking.
	ccRegex := regexp.MustCompile(`\b(?:\d[ -]*?){13,16}\b`)
	if ccRegex.MatchString(body) {
		body = ccRegex.ReplaceAllStringFunc(body, func(s string) string {
			if isValidLuhn(s) {
				redacted = true
				return "[REDACTED_CREDIT_CARD]"
			}
			return s
		})
	}

	// Redact API keys. This regex covers common API key patterns.
	apiKeyRegex := regexp.MustCompile(`(?i)(sk-proj-|AIzaSy|xoxb-|ghp_)[0-9a-zA-Z]{20,}`)
	if apiKeyRegex.MatchString(body) {
		redacted = true
		body = apiKeyRegex.ReplaceAllString(body, "[REDACTED_API_KEY]")
	}

	// Redact email addresses.
	emailRegex := regexp.MustCompile(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`)
	if emailRegex.MatchString(body) {
		redacted = true
		body = emailRegex.ReplaceAllString(body, "[REDACTED_EMAIL]")
	}

	return body, redacted
}

// isValidLuhn implements the Luhn algorithm for credit card number validation.
// This is a more robust check than a simple regex.
func isValidLuhn(s string) bool {
	s = strings.Map(func(r rune) rune {
		if unicode.IsDigit(r) {
			return r
		}
		return -1
	}, s)

	if len(s) < 13 {
		return false
	}

	sum := 0
	double := false
	for i := len(s) - 1; i >= 0; i-- {
		digit := int(s[i] - '0')
		if double {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		double = !double
	}
	return sum%10 == 0
}

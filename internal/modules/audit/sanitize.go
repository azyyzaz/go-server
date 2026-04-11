package audit

import (
	"encoding/json"
	"strings"
)

const maxRequestBodyLength = 2048

var sensitiveKeys = map[string]struct{}{
	"password":      {},
	"access_token":  {},
	"refresh_token": {},
	"captcha_code":  {},
}

func SummarizeRequestBody(body []byte) string {
	body = bytesTrimSpace(body)
	if len(body) == 0 {
		return ""
	}

	var payload any
	if err := json.Unmarshal(body, &payload); err == nil {
		redactPayload(payload)
		if cleaned, err := json.Marshal(payload); err == nil {
			return truncate(string(cleaned), maxRequestBodyLength)
		}
	}
	return truncate(string(body), maxRequestBodyLength)
}

func ExtractUsername(body []byte) string {
	body = bytesTrimSpace(body)
	if len(body) == 0 {
		return ""
	}

	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return ""
	}
	if username, ok := payload["username"].(string); ok {
		return strings.TrimSpace(username)
	}
	return ""
}

func redactPayload(v any) {
	switch value := v.(type) {
	case map[string]any:
		for key, item := range value {
			if _, ok := sensitiveKeys[strings.ToLower(key)]; ok {
				value[key] = "[REDACTED]"
				continue
			}
			redactPayload(item)
		}
	case []any:
		for _, item := range value {
			redactPayload(item)
		}
	}
}

func bytesTrimSpace(body []byte) []byte {
	return []byte(strings.TrimSpace(string(body)))
}

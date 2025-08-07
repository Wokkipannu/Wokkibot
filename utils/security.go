package utils

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

func ValidateURL(inputURL string, allowedSchemes map[string]bool, shellDangerousChars []string) error {
	if inputURL == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	if !allowedSchemes[parsedURL.Scheme] {
		return fmt.Errorf("unsupported URL scheme: %s (only http/https allowed)", parsedURL.Scheme)
	}

	if parsedURL.Hostname() == "" {
		return fmt.Errorf("URL must have a valid hostname")
	}

	for _, dangerousChar := range shellDangerousChars {
		if strings.Contains(inputURL, dangerousChar) {
			return fmt.Errorf("URL contains potentially dangerous character: %s", dangerousChar)
		}
	}

	if strings.Contains(inputURL, "$(") || strings.Contains(inputURL, ")`") {
		return fmt.Errorf("URL contains command injection patterns")
	}

	hostname := strings.ToLower(parsedURL.Hostname())
	if hostname == "localhost" || hostname == "127.0.0.1" ||
		strings.HasPrefix(hostname, "192.168.") ||
		strings.HasPrefix(hostname, "10.") ||
		strings.HasPrefix(hostname, "172.16.") {
		return fmt.Errorf("access to local/internal networks is not allowed")
	}

	return nil
}

func ValidateTimeParameter(timeParam string, timeParamRegex *regexp.Regexp, timeDangerousChars []string) error {
	if timeParam == "" {
		return nil
	}

	for _, dangerousChar := range timeDangerousChars {
		if strings.Contains(timeParam, dangerousChar) {
			return fmt.Errorf("time parameter contains dangerous character: %s", dangerousChar)
		}
	}

	if !timeParamRegex.MatchString(timeParam) {
		return fmt.Errorf("invalid time format. Use formats like: 90, 1:30, or 1:30:45")
	}

	return nil
}

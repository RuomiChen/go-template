package utils

import "strings"

func ParseDeviceFromUA(ua string) string {
	device := "Unknown Device"

	switch {
	case strings.Contains(ua, "Windows"):
		device = "Windows PC"
	case strings.Contains(ua, "Macintosh"):
		device = "MacOS"
	case strings.Contains(ua, "iPhone"):
		device = "iPhone"
	case strings.Contains(ua, "iPad"):
		device = "iPad"
	case strings.Contains(ua, "Android"):
		device = "Android"
	}

	browser := "Unknown Browser"
	switch {
	case strings.Contains(ua, "Chrome"):
		browser = "Chrome"
	case strings.Contains(ua, "Safari") && !strings.Contains(ua, "Chrome"):
		browser = "Safari"
	case strings.Contains(ua, "Firefox"):
		browser = "Firefox"
	case strings.Contains(ua, "Edg"):
		browser = "Edge"
	}

	return device + " / " + browser
}

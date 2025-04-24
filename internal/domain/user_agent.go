package domain

import (
	"fmt"
	"strings"
)

type UserAgent = string

// WARNING: important! files in static folder naming with this
const (
	IOS     UserAgent = "ios"
	ANDROID UserAgent = "android"
	WEB     UserAgent = "web"
)

// TODO: затянуть библиотеку по ua
func ParseUserAgent(ua string) UserAgent {
	fmt.Println(ua)
	lowerUA := strings.ToLower(ua)
	if strings.Contains(lowerUA, "iphone") || strings.Contains(lowerUA, "ipad") {
		return IOS
	}
	if strings.Contains(lowerUA, "android") {
		return ANDROID
	}
	return WEB
}

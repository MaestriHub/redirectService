package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseUserAgent_IOS(t *testing.T) {
	tests := []struct {
		name string
		ua   string
		want UserAgent
	}{
		{"iPhone lowercase", "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X)", IOS},
		{"iPad uppercase", "Mozilla/5.0 (iPad; CPU OS 15_0 like Mac OS X)", IOS},
		{"Mixed case", "Mozilla/5.0 (iPhOnE; CPU iPhone OS 15_0 like Mac OS X)", IOS},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseUserAgent(tt.ua)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseUserAgent_ANDROID(t *testing.T) {
	tests := []struct {
		name string
		ua   string
		want UserAgent
	}{
		{"Android lowercase", "Mozilla/5.0 (Linux; Android 11; Pixel 5)", ANDROID},
		{"Android uppercase", "Mozilla/5.0 (Linux; ANDROID 11; Pixel 5)", ANDROID},
		{"Mixed case", "Mozilla/5.0 (Linux; AnDrOiD 11; Pixel 5)", ANDROID},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseUserAgent(tt.ua)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseUserAgent_WEB(t *testing.T) {
	tests := []struct {
		name string
		ua   string
		want UserAgent
	}{
		{"Empty string", "", WEB},
		{"Unknown UA", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)", WEB},
		{"No keywords", "Some random string", WEB},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseUserAgent(tt.ua)
			assert.Equal(t, tt.want, got)
		})
	}
}

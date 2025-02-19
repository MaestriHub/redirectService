package test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"redirectServer/models"
	"redirectServer/routers"
	"testing"
)

func TestWithoutLinkInput(t *testing.T) {
	directLink := NewDirectLinkBuilder().Build(DB)
	fingerprint := NewFingerprintBuilder().SetVersion("124124").SetDirectLink(*directLink).Build(DB)

	input := models.ParticalFingerprint{
		Platform:       fingerprint.Platform,
		Version:        fingerprint.Version,
		Language:       fingerprint.Language,
		Languages:      fingerprint.Languages,
		Cores:          fingerprint.Cores,
		Memory:         fingerprint.Memory,
		ScreenWidth:    fingerprint.ScreenWidth,
		ScreenHeight:   fingerprint.ScreenHeight,
		ColorDepth:     fingerprint.ColorDepth,
		PixelRatio:     fingerprint.PixelRatio,
		ViewportWidth:  fingerprint.ViewportWidth,
		ViewportHeight: fingerprint.ViewportHeight,
		Renderer:       fingerprint.Renderer,
		VendorRender:   fingerprint.VendorRender,
		TimeZone:       fingerprint.TimeZone,
		UniversalLink:  nil,
	}
	body, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Failed to marshal data: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/find/without-link", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	routers.FindFingerprint(w, req)

	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	var result models.DirectLink
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Errorf("Ошибка парсинга JSON:%v", err)
		return
	}

	if result.ID != directLink.ID {
		t.Errorf("Expected %v, got %v", directLink.ID, result.ID)
	}

	if result.Event != directLink.Event {
		t.Errorf("Expected %v, got %v", directLink.Event, result.Event)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %v", res.StatusCode)
	}
}

func TestWithLinkInput(t *testing.T) {
	directLink := NewDirectLinkBuilder().Build(DB)
	fingerprint := NewFingerprintBuilder().SetVersion("124124").SetDirectLink(*directLink).Build(DB)
	link := directLink.ParseToURL()
	input := models.ParticalFingerprint{
		Platform:       fingerprint.Platform,
		Version:        "hueta",
		Language:       fingerprint.Language,
		Languages:      fingerprint.Languages,
		Cores:          fingerprint.Cores,
		Memory:         fingerprint.Memory,
		ScreenWidth:    fingerprint.ScreenWidth,
		ScreenHeight:   fingerprint.ScreenHeight,
		ColorDepth:     fingerprint.ColorDepth,
		PixelRatio:     fingerprint.PixelRatio,
		ViewportWidth:  fingerprint.ViewportWidth,
		ViewportHeight: fingerprint.ViewportHeight,
		Renderer:       fingerprint.Renderer,
		VendorRender:   fingerprint.VendorRender,
		TimeZone:       fingerprint.TimeZone,
		UniversalLink:  &link,
	}
	body, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Failed to marshal data: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/find/without-link", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	routers.FindFingerprint(w, req)

	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	var result models.DirectLink
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Errorf("Ошибка парсинга JSON:%v", err)
		return
	}

	if result.ID != directLink.ID {
		t.Errorf("Expected %v, got %v", directLink.ID, result.ID)
	}

	if result.Event != directLink.Event {
		t.Errorf("Expected %v, got %v", directLink.Event, result.Event)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %v", res.StatusCode)
	}
}

func TestNonFound(t *testing.T) {
	directLink := NewDirectLinkBuilder().Build(DB)
	fingerprint := NewFingerprintBuilder().SetVersion("124124").SetDirectLink(*directLink).Build(DB)
	input := models.ParticalFingerprint{
		Platform:       fingerprint.Platform,
		Version:        "hueta",
		Language:       fingerprint.Language,
		Languages:      fingerprint.Languages,
		Cores:          fingerprint.Cores,
		Memory:         fingerprint.Memory,
		ScreenWidth:    fingerprint.ScreenWidth,
		ScreenHeight:   fingerprint.ScreenHeight,
		ColorDepth:     fingerprint.ColorDepth,
		PixelRatio:     fingerprint.PixelRatio,
		ViewportWidth:  fingerprint.ViewportWidth,
		ViewportHeight: fingerprint.ViewportHeight,
		Renderer:       fingerprint.Renderer,
		VendorRender:   fingerprint.VendorRender,
		TimeZone:       fingerprint.TimeZone,
		UniversalLink:  nil,
	}
	body, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Failed to marshal data: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/find/without-link", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	routers.FindFingerprint(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status 200 OK, got %v", res.StatusCode)
	}
}

package test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"redirectServer/internal/domain"
	"redirectServer/internal/repo/models/payload"
	"redirectServer/internal/transport/dto"
	"redirectServer/internal/transport/dto/params"
	"redirectServer/routers"

	"github.com/google/uuid"
)

func TestCreateSalonInvite(t *testing.T) {
	parsedUUID, _ := uuid.Parse("c6acaff7-29a5-4c60-b8b8-4be02503bd8b")
	payload := parameters.Salon{
		ID: parsedUUID,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal data: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/create/salon", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	routers.CreateSalonInvite(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %v", res.StatusCode)
	}
}
func TestCreateEmployeerInvite(t *testing.T) {
	parsedUUID, _ := uuid.Parse("53bb0f86-a94e-4302-8a07-ea0b083d3bde")
	payload := payload.MasterToSalon{
		EmployeeId: parsedUUID,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal data: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/create/master-to-salon", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	routers.CreateMasterToSalonInvite(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %v", res.StatusCode)
	}
}
func TestCreateCustomerInvite(t *testing.T) {
	parsedUUID, _ := uuid.Parse("c6acaff7-29a5-4c60-b8b8-4be02503bd8b")
	payload := dto.Customer{
		ID: parsedUUID,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal data: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/create/customer", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	routers.CreateCustomerInvite(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %v", res.StatusCode)
	}
}
func TestCreateMasterToSalonInvite(t *testing.T) {
	parsedUUID, _ := uuid.Parse("c6acaff7-29a5-4c60-b8b8-4be02503bd8b")
	parsedUUID2, _ := uuid.Parse("53bb0f86-a94e-4302-8a07-ea0b083d3bde")
	payload := payload.Employeer{
		ID:      parsedUUID2,
		SalonId: parsedUUID,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal data: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/create/employer", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	routers.CreateEmployeerInvite(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %v", res.StatusCode)
	}
}
func TestWithoutLinkInput(t *testing.T) {
	directLink := NewDirectLinkBuilder().Build(DB)
	fingerprint := NewFingerprintBuilder().SetVersion("124124").SetDirectLink(*directLink).Build(DB)

	input := domain.ParticalFingerprint{
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
	var result domain.DirectLink
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
	input := domain.ParticalFingerprint{
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
	var result domain.DirectLink
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
	input := domain.ParticalFingerprint{
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

func TestGetDirectLink(t *testing.T) {
	NewDirectLinkBuilder().Build(DB)

	req := httptest.NewRequest(http.MethodPost, "/find/with-link/?code=YSg6Ugcf", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	routers.GetDirectLinkPayload(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	var result domain.DirectLink
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Errorf("Ошибка парсинга JSON:%v", err)
		return
	}

	if result.ID != "YSg6Ugcf" {
		t.Errorf("Expected %v, got %v", "YSg6Ugcf", result.ID)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %v", res.StatusCode)
	}
}

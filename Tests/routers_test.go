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

func TestOrganicInput(t *testing.T) {
	if DB == nil {
		t.Fatal("DB is not initialized")
	}
	cores := 2
	input := models.ParticalFingerprint{
		Platform:       "iPhone",
		Version:        "12.32",
		Language:       "ru",
		Languages:      []string{"ru"},
		Cores:          &cores,
		Memory:         nil,
		ScreenWidth:    12,
		ScreenHeight:   15,
		ColorDepth:     1,
		PixelRatio:     23,
		ViewportWidth:  21,
		ViewportHeight: 12,
		Renderer:       "ASF",
		VendorRender:   nil,
		TimeZone:       "Africa",
		UniversalLink:  nil,
	}
	body, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Failed to marshal data: %v", err)
	}
	NewFingerprintBuilder().SetVersion("124124").Build(DB)
	req := httptest.NewRequest(http.MethodPost, "/findFingerprint", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	routers.FindFingerprint(w, req)

	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	t.Logf("Response body: %s", data)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %v", res.StatusCode)
	}
}

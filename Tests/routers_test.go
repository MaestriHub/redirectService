package test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"redirectServer/routers"
	"testing"
)

func TestUpperCaseHandler(t *testing.T) {
	if DB == nil {
		t.Fatal("DB is not initialized")
	}
	req := httptest.NewRequest(http.MethodPost, "/findRequester", nil)
	w := httptest.NewRecorder()
	requester, err := NewRequesterBuilder().SetVersion("124124").Build(DB)
	requester.ColorDepth = 12
	routers.FindRequester(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if string(data) != "ABC" {
		t.Errorf("expected ABC got %v", string(data))
	}
}

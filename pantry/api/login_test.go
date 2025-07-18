package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var TestConfig = Config{
	Secret:   "SuperTestSecret",
	DBUrl:    "data/pantry_pal_test.db",
	Port:     "8080",
	Renderer: &MockRenderer{},
}

func TestIndex(t *testing.T) {
	writer := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/login", nil)

	TestConfig.Index(writer, request)

	if writer.Header().Get("HX-Replace-Url") != "/login" {
		t.Errorf("Expected HX-Replace-Url to be set")
	}

	if writer.Result().StatusCode != 200 {
		t.Errorf("Expected 200 status code. Got: %d", writer.Result().StatusCode)
	}
}

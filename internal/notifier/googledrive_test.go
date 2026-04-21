package notifier

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

func newGoogleDriveSecret(days float64) *vault.Secret {
	return &vault.Secret{
		Path:      "secret/gdrive/token",
		ExpiresAt: time.Now().Add(time.Duration(days*24) * time.Hour),
	}
}

func TestNewGoogleDriveNotifier_Valid(t *testing.T) {
	n, err := NewGoogleDriveNotifier("sheet123", "Alerts", "apikey")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n.spreadsheetID != "sheet123" {
		t.Errorf("expected sheet123, got %s", n.spreadsheetID)
	}
	if n.sheetName != "Alerts" {
		t.Errorf("expected Alerts, got %s", n.sheetName)
	}
}

func TestNewGoogleDriveNotifier_DefaultSheetName(t *testing.T) {
	n, err := NewGoogleDriveNotifier("sheet123", "", "apikey")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n.sheetName != "Alerts" {
		t.Errorf("expected default sheet name Alerts, got %s", n.sheetName)
	}
}

func TestNewGoogleDriveNotifier_MissingSpreadsheetID(t *testing.T) {
	_, err := NewGoogleDriveNotifier("", "Alerts", "apikey")
	if err == nil {
		t.Fatal("expected error for missing spreadsheet ID")
	}
}

func TestNewGoogleDriveNotifier_MissingAPIKey(t *testing.T) {
	_, err := NewGoogleDriveNotifier("sheet123", "Alerts", "")
	if err == nil {
		t.Fatal("expected error for missing API key")
	}
}

func TestGoogleDriveNotifier_Notify_ExpiringSoon(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer ts.Close()

	n, _ := NewGoogleDriveNotifier("sheet123", "Alerts", "key")
	n.baseURL = ts.URL

	if err := n.Notify(newGoogleDriveSecret(5)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received == nil {
		t.Error("expected request body to be received")
	}
}

func TestGoogleDriveNotifier_Notify_Expired(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer ts.Close()

	n, _ := NewGoogleDriveNotifier("sheet123", "Alerts", "key")
	n.baseURL = ts.URL

	if err := n.Notify(newGoogleDriveSecret(-1)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGoogleDriveNotifier_Notify_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	n, _ := NewGoogleDriveNotifier("sheet123", "Alerts", "key")
	n.baseURL = ts.URL

	if err := n.Notify(newGoogleDriveSecret(3)); err == nil {
		t.Fatal("expected error on server error response")
	}
}

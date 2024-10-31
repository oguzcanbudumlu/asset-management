package wallet

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupMockServer(handlerFunc http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handlerFunc)
}

func TestOne_Success(t *testing.T) {
	mockServer := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	defer mockServer.Close()

	adapter := NewValidationAdapter(mockServer.URL)
	err := adapter.One("testWallet", "testNetwork")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestOne_Failure(t *testing.T) {
	mockServer := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})
	defer mockServer.Close()

	adapter := NewValidationAdapter(mockServer.URL)
	err := adapter.One("testWallet", "testNetwork")

	if err == nil || err.Error() != "failed to validate wallet" {
		t.Errorf("expected error 'failed to validate wallet', got %v", err)
	}
}

func TestBoth_Success(t *testing.T) {
	mockServer := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	defer mockServer.Close()

	adapter := NewValidationAdapter(mockServer.URL)
	err := adapter.Both("fromWallet", "toWallet", "testNetwork")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestBoth_SourceValidationFailure(t *testing.T) {
	mockServer := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/wallet/testNetwork/fromWallet" {
			http.Error(w, "source validation failed", http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	})
	defer mockServer.Close()

	adapter := NewValidationAdapter(mockServer.URL)
	err := adapter.Both("fromWallet", "toWallet", "testNetwork")

	if err == nil || err.Error() != "source wallet validation failed: failed to validate wallet" {
		t.Errorf("expected error 'source wallet validation failed', got %v", err)
	}
}

func TestBoth_DestinationValidationFailure(t *testing.T) {
	mockServer := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/wallet/testNetwork/toWallet" {
			http.Error(w, "destination validation failed", http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	})
	defer mockServer.Close()

	adapter := NewValidationAdapter(mockServer.URL)
	err := adapter.Both("fromWallet", "toWallet", "testNetwork")

	if err == nil || err.Error() != "destination wallet validation failed: failed to validate wallet" {
		t.Errorf("expected error 'destination wallet validation failed', got %v", err)
	}
}

func TestBoth_ParallelFailures(t *testing.T) {
	mockServer := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "validation failed", http.StatusBadRequest)
	})
	defer mockServer.Close()

	adapter := NewValidationAdapter(mockServer.URL)
	err := adapter.Both("fromWallet", "toWallet", "testNetwork")

	if err == nil {
		t.Errorf("expected error but got none")
	} else if err.Error() != "source wallet validation failed: failed to validate wallet" &&
		err.Error() != "destination wallet validation failed: failed to validate wallet" {
		t.Errorf("unexpected error: %v", err)
	}
}

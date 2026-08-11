package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAuthenticationRejectsFunctionRoutesWithoutUser(t *testing.T) {
	handler := testApp(t, func(*http.Request) (string, error) { return "", nil })

	for _, target := range []string{"/?action=list", "/?action=save"} {
		t.Run(target, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, target, nil))

			if recorder.Code != http.StatusUnauthorized {
				t.Fatalf("status = %d, want %d", recorder.Code, http.StatusUnauthorized)
			}
		})
	}
}

func TestPageAndStaticAssetsLoadBeforeAuthentication(t *testing.T) {
	handler := testApp(t, func(*http.Request) (string, error) { return "", nil })

	for _, target := range []string{"/", "/?asset=i18n", "/images/icon_256.png"} {
		t.Run(target, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, target, nil))

			if recorder.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
			}
		})
	}
}

func TestAuthenticationRejectsCommandError(t *testing.T) {
	handler := testApp(t, func(*http.Request) (string, error) { return "", errors.New("authenticate failed") })
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/?action=list", nil))

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusUnauthorized)
	}
	if !strings.Contains(recorder.Body.String(), `"error":"unauthorized"`) {
		t.Fatalf("API response = %q, want JSON unauthorized error", recorder.Body.String())
	}
}

func TestAuthenticatedAndDevelopmentRequestsSucceed(t *testing.T) {
	for _, authenticate := range []func(*http.Request) (string, error){
		func(*http.Request) (string, error) { return "dsm-user\n", nil },
		nil,
	} {
		handler := testApp(t, authenticate)
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/?action=list", nil))

		if recorder.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
		}
	}
}

func TestDSMTokenPrefersHeaderAndFallsBackToQuery(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/?SynoToken=query-token", nil)
	if token := dsmSynoToken(request); token != "query-token" {
		t.Fatalf("query token = %q, want query-token", token)
	}

	request.Header.Set("X-Syno-Token", "header-token")
	if token := dsmSynoToken(request); token != "header-token" {
		t.Fatalf("header token = %q, want header-token", token)
	}
}

func testApp(t *testing.T, authenticate func(*http.Request) (string, error)) *app {
	t.Helper()
	handler, err := newApp(t.TempDir()+"/devices.json", authenticate)
	if err != nil {
		t.Fatal(err)
	}
	return handler
}

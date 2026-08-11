package main

import (
	"errors"
	"fmt"
	"io"
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

func TestUpdateReportsMatchingGitHubReleaseAsset(t *testing.T) {
	originalVersion := packageVersion
	packageVersion = "1.0.2-0001"
	t.Cleanup(func() { packageVersion = originalVersion })
	handler := testApp(t, nil)
	handler.releaseAPI = "https://api.github.test/releases/latest"
	handler.httpClient = &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.String() != handler.releaseAPI {
			return nil, fmt.Errorf("request URL = %s", r.URL)
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"tag_name":"v1.0.3-0001","html_url":"https://github.com/jerryt92/wol-spk/releases/tag/v1.0.3-0001","assets":[{"name":"WOLManager-1.0.3-0001-x86_64.spk","browser_download_url":"https://example.test/WOLManager-1.0.3-0001-x86_64.spk"}]}`)),
		}, nil
	})}

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/?action=update", nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
	if !strings.Contains(recorder.Body.String(), `"available":true`) || !strings.Contains(recorder.Body.String(), `"latestVersion":"1.0.3-0001"`) {
		t.Fatalf("response = %s, want available update", recorder.Body.String())
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

func TestComparePackageVersions(t *testing.T) {
	for _, test := range []struct {
		left, right string
		want        int
	}{
		{"1.0.2-0001", "1.0.2-0001", 0},
		{"1.0.3-0001", "1.0.2-9999", 1},
		{"1.0.2-0001", "1.0.2-0002", -1},
	} {
		got := comparePackageVersions(test.left, test.right)
		if (got > 0) != (test.want > 0) || (got < 0) != (test.want < 0) {
			t.Errorf("comparePackageVersions(%q, %q) = %d, want sign %d", test.left, test.right, got, test.want)
		}
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

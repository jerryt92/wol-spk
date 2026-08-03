package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"net/http/cgi"
	"os"
	"path/filepath"
	"strings"

	"wolmanager/internal/store"
	"wolmanager/internal/wol"
)

//go:embed page.html
var pageHTML string

//go:embed i18n.generated.js
var i18nJS []byte

//go:embed images/*
var imagesFS embed.FS

type app struct {
	store  *store.Store
	images http.Handler
}

func main() {
	dataDir := os.Getenv("WOLMANAGER_DATA_DIR")
	if dataDir == "" {
		dataDir = "/var/packages/WOLManager/var"
	}

	images, err := fs.Sub(imagesFS, "images")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	handler := &app{
		store:  store.New(filepath.Join(dataDir, "devices.json")),
		images: http.StripPrefix("/images/", http.FileServer(http.FS(images))),
	}
	if os.Getenv("GATEWAY_INTERFACE") != "" {
		if err := cgi.Serve(handler); err != nil {
			fmt.Printf("Content-Type: text/plain\r\n\r\n%v\n", err)
		}
		return
	}

	addr := getenv("WOLMANAGER_ADDR", "127.0.0.1:8088")
	fmt.Printf("WOL Manager dev server listening on http://%s\n", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func (a *app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/images/") {
		a.images.ServeHTTP(w, r)
		return
	}

	if r.URL.Query().Get("asset") == "i18n" {
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(i18nJS)
		return
	}

	action := r.URL.Query().Get("action")
	if action == "" {
		a.page(w)
		return
	}

	switch action {
	case "list":
		a.list(w)
	case "save":
		a.save(w, r)
	case "delete":
		a.delete(w, r)
	case "replace":
		a.replace(w, r)
	case "wake":
		a.wake(w, r)
	default:
		http.Error(w, "unknown action", http.StatusNotFound)
	}
}

func (a *app) page(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	tpl := template.Must(template.New("page").Parse(pageHTML))
	_ = tpl.Execute(w, nil)
}

func (a *app) list(w http.ResponseWriter) {
	devices, err := a.store.List()
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"devices": devices})
}

func (a *app) save(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, fmt.Errorf("method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	var device store.Device
	if err := readJSON(r, &device); err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}
	saved, err := a.store.Save(device)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"device": saved})
}

func (a *app) delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, fmt.Errorf("method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	var body struct {
		ID string `json:"id"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}
	if err := a.store.Delete(body.ID); err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (a *app) replace(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, fmt.Errorf("method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	var body struct {
		Devices []store.Device `json:"devices"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}
	devices, err := a.store.Replace(body.Devices)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"devices": devices})
}

func (a *app) wake(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, fmt.Errorf("method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	var body struct {
		ID string `json:"id"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}
	device, ok, err := a.store.Get(body.ID)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	if !ok {
		writeError(w, fmt.Errorf("device not found"), http.StatusNotFound)
		return
	}

	if err := wol.Wake(device.MAC, device.Broadcast, device.Port); err != nil {
		writeError(w, err, http.StatusBadGateway)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func readJSON(r *http.Request, target any) error {
	defer r.Body.Close()
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, err error, status int) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func getenv(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

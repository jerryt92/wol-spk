package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"wolmanager/internal/wol"
)

type Device struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	MAC       string `json:"mac"`
	Broadcast string `json:"broadcast"`
	Port      int    `json:"port"`
	Notes     string `json:"notes"`
	UpdatedAt string `json:"updatedAt"`
}

type Store struct {
	path string
}

func New(path string) *Store {
	return &Store{path: path}
}

func (s *Store) List() ([]Device, error) {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return []Device{}, nil
	}
	if err != nil {
		return nil, err
	}
	if len(strings.TrimSpace(string(data))) == 0 {
		return []Device{}, nil
	}

	var devices []Device
	if err := json.Unmarshal(data, &devices); err != nil {
		return nil, fmt.Errorf("read devices: %w", err)
	}
	return devices, nil
}

func (s *Store) Save(device Device) (Device, error) {
	devices, err := s.List()
	if err != nil {
		return Device{}, err
	}

	device, err = normalizeDevice(device)
	if err != nil {
		return Device{}, err
	}

	replaced := false
	for i := range devices {
		if devices[i].ID == device.ID {
			devices[i] = device
			replaced = true
			break
		}
	}
	if !replaced {
		devices = append(devices, device)
	}

	if err := s.write(devices); err != nil {
		return Device{}, err
	}
	return device, nil
}

func (s *Store) Replace(devices []Device) ([]Device, error) {
	normalized := make([]Device, 0, len(devices))
	seen := map[string]bool{}
	for _, device := range devices {
		next, err := normalizeDevice(device)
		if err != nil {
			return nil, err
		}
		if seen[next.ID] {
			next.ID = newUniqueID(seen)
		}
		seen[next.ID] = true
		normalized = append(normalized, next)
	}
	if err := s.write(normalized); err != nil {
		return nil, err
	}
	return normalized, nil
}

func (s *Store) Delete(id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return errors.New("id is required")
	}

	devices, err := s.List()
	if err != nil {
		return err
	}
	next := make([]Device, 0, len(devices))
	for _, device := range devices {
		if device.ID != id {
			next = append(next, device)
		}
	}
	return s.write(next)
}

func (s *Store) Get(id string) (Device, bool, error) {
	devices, err := s.List()
	if err != nil {
		return Device{}, false, err
	}
	for _, device := range devices {
		if device.ID == id {
			return device, true, nil
		}
	}
	return Device{}, false, nil
}

func (s *Store) write(devices []Device) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(devices, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(s.path, data, 0600)
}

func normalizeDevice(device Device) (Device, error) {
	normalized, err := wol.NormalizeMAC(device.MAC)
	if err != nil {
		return Device{}, err
	}
	device.MAC = normalized
	device.Name = strings.TrimSpace(device.Name)
	device.Broadcast = strings.TrimSpace(device.Broadcast)
	device.Notes = strings.TrimSpace(device.Notes)
	if device.Name == "" {
		return Device{}, errors.New("name is required")
	}
	if device.Broadcast == "" {
		device.Broadcast = "255.255.255.255"
	}
	if device.Port == 0 {
		device.Port = 9
	}
	if device.Port < 1 || device.Port > 65535 {
		return Device{}, errors.New("port must be between 1 and 65535")
	}
	if strings.TrimSpace(device.ID) == "" {
		device.ID = newID()
	}
	device.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	return device, nil
}

func newID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func newUniqueID(seen map[string]bool) string {
	for {
		id := newID()
		if !seen[id] {
			return id
		}
		time.Sleep(time.Nanosecond)
	}
}

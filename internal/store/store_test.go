package store

import (
	"path/filepath"
	"testing"
)

func TestSaveListDelete(t *testing.T) {
	st := New(filepath.Join(t.TempDir(), "devices.json"))

	device, err := st.Save(Device{Name: "Office PC", MAC: "aa-bb-cc-dd-ee-ff"})
	if err != nil {
		t.Fatal(err)
	}
	if device.ID == "" {
		t.Fatal("expected generated id")
	}
	if device.Broadcast != "255.255.255.255" {
		t.Fatalf("broadcast = %q", device.Broadcast)
	}
	if device.Port != 9 {
		t.Fatalf("port = %d", device.Port)
	}

	devices, err := st.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(devices) != 1 {
		t.Fatalf("device count = %d", len(devices))
	}
	if devices[0].MAC != "AA:BB:CC:DD:EE:FF" {
		t.Fatalf("stored MAC = %q", devices[0].MAC)
	}

	if err := st.Delete(device.ID); err != nil {
		t.Fatal(err)
	}
	devices, err = st.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(devices) != 0 {
		t.Fatalf("device count after delete = %d", len(devices))
	}
}

func TestReplaceNormalizesDevices(t *testing.T) {
	st := New(filepath.Join(t.TempDir(), "devices.json"))

	devices, err := st.Replace([]Device{
		{Name: "PC", MAC: "aa-bb-cc-dd-ee-ff"},
		{Name: "Server", MAC: "01:23:45:67:89:ab", Broadcast: "192.168.1.255", Port: 7},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(devices) != 2 {
		t.Fatalf("device count = %d", len(devices))
	}
	if devices[0].MAC != "AA:BB:CC:DD:EE:FF" || devices[0].Port != 9 {
		t.Fatalf("unexpected first device: %+v", devices[0])
	}
	if devices[1].MAC != "01:23:45:67:89:AB" || devices[1].Port != 7 {
		t.Fatalf("unexpected second device: %+v", devices[1])
	}
}

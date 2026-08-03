package wol

import "testing"

func TestNormalizeMAC(t *testing.T) {
	got, err := NormalizeMAC("aa-bb-cc-dd-ee-ff")
	if err != nil {
		t.Fatal(err)
	}
	if got != "AA:BB:CC:DD:EE:FF" {
		t.Fatalf("NormalizeMAC() = %q", got)
	}
}

func TestNormalizeMACAcceptsCompactInput(t *testing.T) {
	got, err := NormalizeMAC("aabbccddeeff")
	if err != nil {
		t.Fatal(err)
	}
	if got != "AA:BB:CC:DD:EE:FF" {
		t.Fatalf("NormalizeMAC() = %q", got)
	}
}

func TestNormalizeMACRejectsInvalidInput(t *testing.T) {
	if _, err := NormalizeMAC("AA:BB:CC"); err == nil {
		t.Fatal("expected invalid MAC error")
	}
}

func TestMagicPacket(t *testing.T) {
	packet, err := MagicPacket("01:23:45:67:89:AB")
	if err != nil {
		t.Fatal(err)
	}
	if len(packet) != 102 {
		t.Fatalf("packet length = %d", len(packet))
	}
	for i := 0; i < 6; i++ {
		if packet[i] != 0xff {
			t.Fatalf("sync byte %d = %x", i, packet[i])
		}
	}
	want := []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xAB}
	for repeat := 0; repeat < 16; repeat++ {
		offset := 6 + repeat*6
		for i := range want {
			if packet[offset+i] != want[i] {
				t.Fatalf("repeat %d byte %d = %x", repeat, i, packet[offset+i])
			}
		}
	}
}

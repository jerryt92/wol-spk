package wol

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"
	"syscall"
	"time"
)

var (
	macPattern      = regexp.MustCompile(`(?i)^([0-9a-f]{2}[:-]){5}[0-9a-f]{2}$`)
	compactPattern  = regexp.MustCompile(`(?i)^[0-9a-f]{12}$`)
	nonHexSeparator = regexp.MustCompile(`[^0-9a-fA-F]`)
)

func NormalizeMAC(input string) (string, error) {
	mac := strings.TrimSpace(input)
	compact := nonHexSeparator.ReplaceAllString(mac, "")
	if compactPattern.MatchString(compact) {
		compact = strings.ToUpper(compact)
		return strings.Join([]string{
			compact[0:2],
			compact[2:4],
			compact[4:6],
			compact[6:8],
			compact[8:10],
			compact[10:12],
		}, ":"), nil
	}
	if !macPattern.MatchString(mac) {
		return "", errors.New("MAC address must look like AA:BB:CC:DD:EE:FF")
	}
	mac = strings.ToUpper(strings.ReplaceAll(mac, "-", ":"))
	return mac, nil
}

func MagicPacket(mac string) ([]byte, error) {
	normalized, err := NormalizeMAC(mac)
	if err != nil {
		return nil, err
	}

	hw, err := hex.DecodeString(strings.ReplaceAll(normalized, ":", ""))
	if err != nil {
		return nil, fmt.Errorf("decode MAC: %w", err)
	}

	packet := make([]byte, 6+16*len(hw))
	for i := 0; i < 6; i++ {
		packet[i] = 0xff
	}
	for i := 0; i < 16; i++ {
		copy(packet[6+i*len(hw):6+(i+1)*len(hw)], hw)
	}
	return packet, nil
}

func Wake(mac, address string, port int) error {
	if strings.TrimSpace(address) == "" {
		address = "255.255.255.255"
	}
	if port <= 0 || port > 65535 {
		return errors.New("port must be between 1 and 65535")
	}

	packet, err := MagicPacket(mac)
	if err != nil {
		return err
	}

	dialer := net.Dialer{
		Timeout: 3 * time.Second,
		Control: func(network, address string, conn syscall.RawConn) error {
			var controlErr error
			if err := conn.Control(func(fd uintptr) {
				controlErr = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_BROADCAST, 1)
			}); err != nil {
				return err
			}
			return controlErr
		},
	}
	conn, err := dialer.Dial("udp4", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		return fmt.Errorf("open UDP connection: %w", err)
	}
	defer conn.Close()

	if _, err := conn.Write(packet); err != nil {
		return fmt.Errorf("send magic packet: %w", err)
	}
	return nil
}

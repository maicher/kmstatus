package network

import (
	"os"
	"testing"
	"time"

	"github.com/maicher/kmstatus/internal/test"
)

const header = "Inter-|   Receive\n face |bytes    packets errs drop fifo frame compressed multicast|bytes    packets\n"

func writeNetDev(f *os.File, lines string) {
	f.Truncate(0)
	f.Seek(0, 0)
	f.WriteString(header + lines)
}

func Test_Parser_InterfacesChange(t *testing.T) {
	f := test.NewTempFile()
	defer f.Close()

	var d []data
	parser := Parser{file: f, dataBuf: make(map[string]data)}

	writeNetDev(f, "  eth0: 100 1 0 0 0 0 0 0 200 1 0 0 0 0 0 0\n  wlan0:300 1 0 0 0 0 0 0 400 1 0 0 0 0 0 0\n")
	parser.Parse(&d)
	if len(d) != 2 {
		t.Fatalf("Interfaces: %d, want: 2", len(d))
	}
	if d[1].Name != "wlan0" || d[1].RxTotal != 300 || d[1].TxTotal != 400 {
		t.Fatalf("Interface: %+v, want: wlan0 300 400", d[1])
	}

	// wlan0 disappears, usb0 appears.
	parser.parsedAt = time.Now().Add(-time.Second)
	writeNetDev(f, "  eth0: 150 1 0 0 0 0 0 0 200 1 0 0 0 0 0 0\n  usb0: 9000 1 0 0 0 0 0 0 9000 1 0 0 0 0 0 0\n")
	parser.Parse(&d)
	if len(d) != 2 || d[1].Name != "usb0" {
		t.Fatalf("Interfaces: %+v, want: eth0, usb0", d)
	}
	if d[0].Rx != 50 {
		t.Fatalf("eth0 Rx: %d, want: 50", d[0].Rx)
	}
	if d[1].Rx != 0 || d[1].Tx != 0 {
		t.Fatalf("usb0 speed: %d/%d, want: 0/0 for a new interface", d[1].Rx, d[1].Tx)
	}
}

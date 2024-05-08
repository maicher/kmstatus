package network

import (
	"testing"
	"time"

	"github.com/maicher/kmstatus/internal/test"
)

func Test_MemParser(t *testing.T) {
	f := test.NewTempFile()
	defer f.Close()

	f.WriteString("Inter-|   Receive                                                |  Transmit\n")
	f.WriteString("face |bytes    packets errs drop fifo frame compressed multicast|bytes    packets errs drop fifo colls carrier compressed\n")
	f.WriteString("enp6s0: 100 1361812    0    0    0     0          0        19 200  565252    0    0    0     0       0          0\n")

	d := make([]data, 5)
	parser := Parser{file: f}
	parser.dataBuf = make(map[string]data)
	err := parser.Parse(d)
	if err != nil {
		t.Fatalf("Error: %s, want: nil", err)
	}

	if name := d[0].Name; name != "enp6s0" {
		t.Fatalf("Net interface name: %s, want: enp6s0", name)
	}
	if rxTotal := d[0].RxTotal; rxTotal != 100 {
		t.Fatalf("Interface RxTotal: %d, want: 100", rxTotal)
	}
	if txTotal := d[0].TxTotal; txTotal != 200 {
		t.Fatalf("Interface TxTotal: %d, want: 200", txTotal)
	}

	err = f.Truncate(0)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	_, err = f.Seek(0, 0)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	f.WriteString("Inter-|   Receive                                                |  Transmit\n")
	f.WriteString("face |bytes    packets errs drop fifo frame compressed multicast|bytes    packets errs drop fifo colls carrier compressed\n")
	f.WriteString("enp6s0: 150 1361812    0    0    0     0          0        19 260  565252    0    0    0     0       0          0\n")
	if err != nil {
		t.Fatalf("Error: %s", err)
	}

	parser.parsedAt = time.Now().Add(-2 * time.Second)
	err = parser.Parse(d)
	if err != nil {
		t.Fatalf("Error: %s, want: nil", err)
	}
	if rxTotal := d[0].RxTotal; rxTotal != 150 {
		t.Fatalf("Interface RxTotal: %d, want: 150", rxTotal)
	}
	if txTotal := d[0].TxTotal; txTotal != 260 {
		t.Fatalf("Interface TxTotal: %d, want: 260", txTotal)
	}

	if rx := d[0].Rx; rx != 25 {
		t.Fatalf("Interface Rx: %d, want: 25", rx)
	}
	if tx := d[0].Tx; tx != 30 {
		t.Fatalf("Interface Tx: %d, want: 30", tx)
	}
}

package mem

import "testing"

func Test_Data_NoSwap(t *testing.T) {
	d := data{Total: 100, Used: 50}

	if p := d.SwapUsedPercentage(); p != 0 {
		t.Fatalf("SwapUsedPercentage: %f, want: 0", p)
	}
}

package cpu

import "testing"

func Test_Data_FreqGHz(t *testing.T) {
	d := data{Freq: 3000000} // kHz

	if f := d.FreqGHz(); f != 3 {
		t.Fatalf("FreqGHz: %f, want: 3", f)
	}
}

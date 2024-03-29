package mem

import (
	"testing"

	"github.com/maicher/kmst/internal/test"
)

func Test_MemParser(t *testing.T) {
	f := test.NewTempFile()
	test.WriteLine(f, `MemTotal:       32814068 kB
MemFree:          782964 kB
MemAvailable:    6068128 kB
Buffers:           83328 kB
Cached:          5376188 kB
SwapCached:       187088 kB
Active:          7818892 kB
Inactive:        5076488 kB
Active(anon):    5975824 kB
Inactive(anon):  1865080 kB
Active(file):    1843068 kB
Inactive(file):  3211408 kB
Unevictable:         544 kB
Mlocked:             544 kB
SwapTotal:      50331644 kB
SwapFree:       47258620 kB
Zswap:            908400 kB
Zswapped:        2266780 kB
Dirty:            559496 kB
Writeback:         13684 kB
AnonPages:       7382844 kB
Mapped:           902196 kB
Shmem:            404976 kB
KReclaimable:     670988 kB
Slab:             909448 kB
SReclaimable:     670988 kB
SUnreclaim:       238460 kB
KernelStack:       39552 kB
PageTables:        97152 kB
SecPageTables:         0 kB
NFS_Unstable:          0 kB
Bounce:                0 kB
WritebackTmp:          0 kB
CommitLimit:    58333684 kB
Committed_AS:   28171424 kB
VmallocTotal:   34359738367 kB
VmallocUsed:      149628 kB
VmallocChunk:          0 kB
Percpu:            17664 kB
HardwareCorrupted:     0 kB
AnonHugePages:   1744896 kB
ShmemHugePages:        0 kB
ShmemPmdMapped:        0 kB
FileHugePages:     77824 kB
FilePmdMapped:     71680 kB
CmaTotal:              0 kB
CmaFree:               0 kB
HugePages_Total:      16
HugePages_Free:       16
HugePages_Rsvd:        0
HugePages_Surp:        0
Hugepagesize:       2048 kB
Hugetlb:        16809984 kB
DirectMap4k:      285332 kB
DirectMap2M:    10096640 kB
DirectMap1G:    25165824 kB`)

	data := Data{}
	parser := MemParser{file: f}
	err := parser.Parse(&data)

	if err != nil {
		t.Fatalf("Error: %s, want: nil", err)
	}

	if data.Total != 32814068 {
		t.Fatalf("MemTotal: %d, want: %d", data.Total, 32814068)
	}

	if data.Used != 26571588 {
		t.Fatalf("MemUsed: %d, want: %d", data.Used, 26571588)
	}
	if data.SwapTotal != 50331644 {
		t.Fatalf("SwapTotal: %d, want: %d", data.SwapTotal, 50331644)
	}

	if data.SwapUsed != 3073024 {
		t.Fatalf("SwapUsed: %d, want: %d", data.SwapUsed, 3073024)
	}
}

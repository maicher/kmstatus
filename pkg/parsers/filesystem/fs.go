package filesystem

type Drive struct {
	Name      string
	MountedOn string
	Total     int
	Used      int
	Free      int
}

func (d *Drive) UsedPercentage() float64 {
	return 100 * float64(d.Used) / float64(d.Total)
}

func (d *Drive) FreePercentage() float64 {
	return 100 * float64(d.Free) / float64(d.Total)
}

type FS struct {
	Drives []Drive
	ENCFS  bool
}

func NewFS() FS {
	return FS{}
}

func (fs *FS) Find(name string) (drive *Drive) {
	for _, d := range fs.Drives {
		if d.Name == name {
			drive = &d
			return
		}
	}

	return
}

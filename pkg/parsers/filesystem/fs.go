package filesystem

type Drive struct {
	Name      string
	MountedOn string
	Total     int
	Used      int
	Free      int
}

type FS struct {
	Drives []Drive
	ENCFS  bool
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

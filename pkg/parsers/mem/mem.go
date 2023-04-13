package mem

type SpaceKB int

type Mem struct {
	MemTotal  SpaceKB
	MemUsed   SpaceKB
	SwapTotal SpaceKB
	SwapUsed  SpaceKB
}

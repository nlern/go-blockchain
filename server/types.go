package server

type getBlocks struct {
	AddrFrom string
}

type inv struct {
	AddrFrom string
	Type     string
	Items    [][]byte
}

type version struct {
	Version    int
	BestHeight int
	AddrFrom   string
}

package models

type Parameters struct {
	Q      uint64     `json:"q"`
	Count  uint32     `json:"count"`
	Lambda uint32     `json:"lambda"`
	D      [][]uint64 `json:"d"`
}

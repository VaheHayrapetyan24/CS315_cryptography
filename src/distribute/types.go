package main

type Parameters struct {
	Q      uint64     `json:"q"`
	N      uint32     `json:"n"` // probably won't be needing this one if planning to increase dynamically
	Count  uint32     `json:"count"`
	Lambda uint32     `json:"lambda"`
	D      [][]uint64 `json:"d"`
	G      [][]uint64 `json:"g"`
}

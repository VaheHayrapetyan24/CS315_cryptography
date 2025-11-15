package models

type DistributeResponse struct {
	Id   uint32   `json:"id"`
	Gcol []uint64 `json:"g_col"`
	Acol []uint64 `json:"a_col"`
}

type ParametersResponse struct {
	Q uint64     `json:"q"`
	G [][]uint64 `json:"g"`
}

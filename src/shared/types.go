package shared

type DistributeResponse struct {
	Id   uint32   `json:"id"`
	Q    uint64   `json:"q"`
	Gcol []uint64 `json:"g_col"`
	Acol []uint64 `json:"a_col"`
}

type ParametersResponse struct {
	Q uint64     `json:"q"`
	G [][]uint64 `json:"g"`
}

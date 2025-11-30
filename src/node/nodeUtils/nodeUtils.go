package nodeUtils

func GetKey(Gcol []uint64, Acol []uint64, q uint64) uint64 {
	var k uint64 = 0
	for i := 0; i < len(Gcol) && i < len(Acol); i++ {
		k += (Acol[i] * Gcol[i]) % q
	}
	return k % q
}

package util

// BSVToSatoshis converts a bsv coin amount satoshis.
func BSVToSatoshis(f float64) uint64 {
	return uint64(f * 100000000)
}

// SatoshisToBSV an amount of satoshis to bsv coins.
func SatoshisToBSV(s int64) float64 {
	return float64(s) / 100000000
}

// MapBSVToSatoshis converts a string => bsv coin map to string => satoshi.
func MapBSVToSatoshis(vv map[string]float64) map[string]uint64 {
	if vv == nil {
		return nil
	}

	mm := make(map[string]uint64, len(vv))
	for k, v := range vv {
		mm[k] = BSVToSatoshis(v)
	}

	return mm
}

// MapSatoshisToBSV converts a string => satoshi map to string => bsv coin.
func MapSatoshisToBSV(vv map[string]uint64) map[string]float64 {
	if vv == nil {
		return nil
	}

	mm := make(map[string]float64, len(vv))
	for k, v := range vv {
		//nolint:gosec // Bitcoin values never exceed int64 max (21M BTC = 2.1e15 satoshis < 9.2e18)
		mm[k] = SatoshisToBSV(int64(v))
	}

	return mm
}

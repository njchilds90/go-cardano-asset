package cardanoasset

import "fmt"

// bech32Encode encodes data bytes into a bech32 string with the given HRP.
// This is a minimal, zero-dependency bech32 implementation sufficient for
// encoding asset fingerprints per CIP-14.
func bech32Encode(hrp string, data []byte) (string, error) {
	conv, err := convertBits(data, 8, 5, true)
	if err != nil {
		return "", err
	}
	return encodeBech32(hrp, conv)
}

const charset = "qpzry9x8gf2tvdw0s3jn54khce6mua7l"

var gen = []uint32{0x3b6a57b2, 0x26508e6d, 0x1ea119fa, 0x3d4233dd, 0x2a1462b3}

func polymod(values []byte) uint32 {
	chk := uint32(1)
	for _, v := range values {
		top := chk >> 25
		chk = (chk&0x1ffffff)<<5 ^ uint32(v)
		for i := 0; i < 5; i++ {
			if (top>>uint(i))&1 == 1 {
				chk ^= gen[i]
			}
		}
	}
	return chk
}

func hrpExpand(hrp string) []byte {
	result := make([]byte, len(hrp)*2+1)
	for i, c := range hrp {
		result[i] = byte(c >> 5)
		result[i+len(hrp)+1] = byte(c & 31)
	}
	result[len(hrp)] = 0
	return result
}

func createChecksum(hrp string, data []byte) []byte {
	values := append(hrpExpand(hrp), data...)
	values = append(values, []byte{0, 0, 0, 0, 0, 0}...)
	mod := polymod(values) ^ 1
	ret := make([]byte, 6)
	for i := 0; i < 6; i++ {
		ret[i] = byte((mod >> (5 * (5 - i))) & 31)
	}
	return ret
}

func encodeBech32(hrp string, data []byte) (string, error) {
	combined := append(data, createChecksum(hrp, data)...)
	result := hrp + "1"
	for _, b := range combined {
		if int(b) >= len(charset) {
			return "", fmt.Errorf("invalid bech32 data byte: %d", b)
		}
		result += string(charset[b])
	}
	return result, nil
}

func convertBits(data []byte, fromBits, toBits uint, pad bool) ([]byte, error) {
	acc := 0
	bits := uint(0)
	var result []byte
	maxv := (1 << toBits) - 1
	for _, value := range data {
		acc = (acc << fromBits) | int(value)
		bits += fromBits
		for bits >= toBits {
			bits -= toBits
			result = append(result, byte((acc>>bits)&maxv))
		}
	}
	if pad {
		if bits > 0 {
			result = append(result, byte((acc<<(toBits-bits))&maxv))
		}
	} else if bits >= fromBits || ((acc<<(toBits-bits))&maxv) != 0 {
		return nil, fmt.Errorf("invalid padding in bit conversion")
	}
	return result, nil
}

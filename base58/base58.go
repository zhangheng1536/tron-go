package base58

import (
	"crypto/sha256"
	"github.com/btcsuite/btcutil/base58"
)

func From58Check(s string) ([]byte, error) {
	decodeCheck := base58.Decode(s)
	//decodeCheck:=cov(decodeCheckt)
	if len(decodeCheck) <= 4 {
		return nil, nil
	}
	decodeData := decodeCheck[:len(decodeCheck)-4]
	h0 := sha256Do(decodeData)
	h1 := sha256Do(h0)
	if h1[0] == decodeCheck[len(decodeData)] &&
		h1[1] == decodeCheck[len(decodeData)+1] &&
		h1[2] == decodeCheck[len(decodeData)+2] &&
		h1[3] == decodeCheck[len(decodeData)+3] {
		return decodeData, nil
	}
	return nil, nil
}
func To58Check(input []byte) string {
	h0 := sha256Do(input)
	h1 := sha256Do(h0)

	inputCheck := input
	inputCheck = append(inputCheck, h1[:4]...)

	return base58.Encode(inputCheck)
}
func sha256Do(input []byte) []byte {
	h256h0 := sha256.New()
	h256h0.Write(input)
	h0 := h256h0.Sum(nil)
	return h0
}

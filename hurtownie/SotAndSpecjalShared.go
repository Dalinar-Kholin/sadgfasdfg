package hurtownie

import (
	"math/rand"
	"strings"
	"time"
)

func GenerateUUID() string {
	var D [36]byte
	A := "0123456789abcdef"

	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	for B := 0; B < 36; B++ {
		D[B] = A[rand.Intn(16)]
	}
	D[14] = '4'
	D[19] = A[(D[19]&0x3)|0x8]
	D[8] = '-'
	D[13] = '-'
	D[18] = '-'
	D[23] = '-'
	C := strings.Join([]string{string(D[:])}, "")
	return C
}

type SotAndSpecjalTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type SotAndSpecjalResponse struct {
	DataOferty time.Time `json:"dataOferty"`
	Pozycje    []struct {
		Nazwa               string  `json:"nazwa"`
		IlOpkZb             int     `json:"ilOpkZb"`
		CenaNettoOstateczna float64 `json:"cenaNettoOstateczna"`
	} `json:"pozycje"`
	CountPozycji int `json:"countPozycji"`
}

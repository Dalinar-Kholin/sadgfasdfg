package factory

import (
	"errors"
	"optimaHurt/hurtownie"
	"optimaHurt/hurtownie/eurocash"
	"optimaHurt/hurtownie/sot"
	"optimaHurt/hurtownie/specjal"
	"optimaHurt/hurtownie/tedi"
)

func HurtFactory(hurt hurtownie.HurtName) (hurtownie.IHurt, error) {
	switch hurt {
	case hurtownie.Eurocash:
		return &eurocash.EurocashObject{}, nil
	case hurtownie.Specjal:
		return &specjal.Specjal{}, nil
	case hurtownie.Sot:
		return &sot.Sot{}, nil
	case hurtownie.Tedi:
		return &tedi.Tedi{}, nil

	}
	return nil, errors.New("Nie ma takiej hurtowni")
}

// wiatrak w kółku
// 20 pieczenia

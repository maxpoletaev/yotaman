package selfcare

import (
	"math"
	"fmt"
	"strconv"
)

type Tariff struct {
	Amount      int    `json:"amountNumber,string"`
	DaysRemain  int    `json:"remainNumber,string"`
	RawSpeed    string `json:"speedNumber"`
	Code        string `json:"code"`
	IsMax       bool   `json:"isSpeedMaximum"`
	IsLight     bool   `json:"isLight"`
	MoneyEnough bool   `json:"moneyEnough"`
}

// Label returns unique label of tariiff.
func (t *Tariff) Label() string {
	if t.IsMax {
		return "max"
	}
	return t.RawSpeed
}

// Speed in kbps.
func (t *Tariff) Speed() float64 {
	if t.IsMax {
		return math.Inf(1)
	}

	speed, err := strconv.ParseFloat(t.RawSpeed, 64)
	if err != nil {
		panic(err)
	}

	if speed < 20 {
		// if speed in megabytes
		return speed * 1000
	}

	return speed
}

// Repr returns human-readeable representation of trariff.
func (t *Tariff) Repr() string {
	return fmt.Sprintf(
		"%v Kbps, %d Rub, %d days ",
		t.Speed(), t.Amount, t.DaysRemain,
	)
}

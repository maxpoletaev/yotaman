package selfcare

import (
	"encoding/json"
	"net/url"
	"io/ioutil"
	"strconv"
	"regexp"
	"fmt"
)

type Tariff struct {
	Amount      int     `json:"amountNumber,string"`
	DaysRemain  int     `json:"remainNumber,string"`
	RawSpeed    string  `json:"speedNumber"`
	Code        string  `json:"code"`
	IsMax       bool    `json:"isSpeedMaximum"`
	IsLight     bool    `json:"isLight"`
	MoneyEnough bool    `json:"moneyEnough"`
}

func (t *Tariff) Label() string {
	if (t.IsMax) {
		return "max"
	}
	return t.RawSpeed
}

func (t *Tariff) Speed() float64 {
	if t.IsMax {
		return 100000
	}

	speed, err := strconv.ParseFloat(t.RawSpeed, 64)
	if err != nil { panic(err) }

	if speed < 20 {
		// if speed in megabytes
		return speed * 1000
	}

	return speed
}

// Human-readeable representation of trariff.
func (t *Tariff) Repr() string {
	return fmt.Sprintf(
		"%d days, %d Rub, %v Kbps",
		t.DaysRemain, t.Amount, t.Speed(),
	)
}

type Device struct {
	ID            json.Number `json:"productId"`
	Tariffs       []Tariff    `json:"steps"`
	CurrentTariff Tariff      `json:"currentProduct"`
}

func (d *Device) ChangeTariff(t Tariff) {
	form := url.Values{
		"product": {d.ID.String()},
		"offerCode": {t.Code},
		"status": {"custom"},
	}

	if d.CurrentTariff.Code != t.Code {
		_, err := client.PostForm(changeOfferURL, form)
		if err != nil { panic(err) }
	}
}

func (d *Device) IsCurrentTariff(t Tariff) bool {
	return d.CurrentTariff.Code == t.Code
}

func GetDevices() []Device {
	sliderDataRegexp := regexp.MustCompile("var sliderData = (.*);\n")
	matches := sliderDataRegexp.FindSubmatch(LoadPage(devicesURL))

	if len(matches) != 2 {
		panic("Variable sliderData not found on devices page.")
	}

	sliderData := make(map[string]Device)
	err := json.Unmarshal(matches[1], &sliderData)
	if err != nil { panic(err) }

	devices := make([]Device, 0, len(sliderData))
	for _, d := range sliderData {
		devices = append(devices, d)
	}

	return devices
}

func GetCurrentDevice() Device {
	devices := GetDevices()
	if len(devices) < 1 {
		panic("No devices found.")
	}
	return devices[0]
}

func LoadPage(url string) []byte {
	resp, err := client.Get(url)
	if err != nil { panic(err) }
	if resp.StatusCode != 200 {
		panic("Page load error")
	}
	defer resp.Body.Close()
	page, _ := ioutil.ReadAll(resp.Body)
	return page
}

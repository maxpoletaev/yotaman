package selfcare

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"regexp"
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

// Speed in kilobyte per second.
func (t *Tariff) Speed() float64 {
	if t.IsMax {
		// infinitie speed (100 mbps)
		return 100000
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

type Device struct {
	ID            json.Number `json:"productId"`
	Tariffs       []Tariff    `json:"steps"`
	CurrentTariff Tariff      `json:"currentProduct"`
}

func (d *Device) ChangeTariff(t Tariff) error {
	form := url.Values{
		"product":   {d.ID.String()},
		"offerCode": {t.Code},
		"status":    {"custom"},
	}

	if d.CurrentTariff.Code != t.Code {
		_, err := client.PostForm(changeOfferURL, form)
		return err
	}

	return nil
}

func (d *Device) IsCurrentTariff(t Tariff) bool {
	return d.CurrentTariff.Code == t.Code
}

func GetDevices() ([]*Device, error) {
	page, err := LoadPage(devicesURL)
	if err != nil {
		return nil, err
	}

	sliderDataRegexp := regexp.MustCompile("var sliderData = (.*);\n")
	matches := sliderDataRegexp.FindSubmatch(page)

	if len(matches) != 2 {
		return nil, errors.New("variable sliderData not found on devices page")
	}

	sliderData := make(map[string]Device)
	err = json.Unmarshal(matches[1], &sliderData)
	if err != nil {
		return nil, err
	}

	devices := make([]*Device, 0, len(sliderData))
	for _, d := range sliderData {
		devices = append(devices, &d)
	}

	return devices, nil
}

func GetCurrentDevice() (*Device, error) {
	devices, err := GetDevices()
	if err != nil {
		return nil, err
	}
	if len(devices) < 1 {
		return nil, errors.New("no devices found")
	}
	return devices[0], nil
}

func LoadPage(url string) ([]byte, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("page load error")
	}
	defer resp.Body.Close()
	page, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return page, nil
}

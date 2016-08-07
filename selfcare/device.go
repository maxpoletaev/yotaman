package selfcare

import (
	"encoding/json"
	"errors"
	"net/url"
	"regexp"
)

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

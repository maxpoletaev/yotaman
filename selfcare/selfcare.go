package selfcare

import (
	"errors"
	"strings"
	"net/url"
	"net/http"
	"net/http/cookiejar"
)

const (
	changeOfferURL = "https://my.yota.ru/selfcare/devices/changeOffer"
	devicesURL = "https://my.yota.ru/selfcare/devices"
	loginURL = "https://login.yota.ru/UI/Login"
	autoLoginURL = "https://my.yota.ru/selfcare/mydevices"
	loginSuccessURL = "https://my.yota.ru/selfcare/loginSuccess"
	loginErrorURL = "https://my.yota.ru:443/selfcare/loginError"
)

func Login(username string, password string) error {
	form := url.Values{
		"gotoOnFail": {loginErrorURL},
		"goto": {loginSuccessURL},
		"ForceAuth": {"true"},
		"org": {"customer"},
		"IDToken1": {username},
		"IDToken2": {password},
	}

	resp, err := client.PostForm(loginURL, form)
	if err != nil { return err }

	if resp.StatusCode != http.StatusOK {
		return errors.New("Login error")
	}
	return nil
}

func AutoLogin() error {
	resp, err := client.Get(autoLoginURL)
	if err != nil { return err }

	if resp.StatusCode != http.StatusOK {
		return errors.New("login error")
	}

	if strings.Contains(resp.Request.URL.Path, "login") {
		return errors.New("autologin error, used proxy or non-yota provider")
	}

	return nil
}

func createClient() http.Client {
	jar, _ := cookiejar.New(nil)
	client := http.Client{Jar: jar}
	return client
}

var client = createClient()

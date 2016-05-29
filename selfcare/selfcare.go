package selfcare

import (
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

func Login(username string, password string) {
	form := url.Values{
		"gotoOnFail": {loginErrorURL},
		"goto": {loginSuccessURL},
		"ForceAuth": {"true"},
		"org": {"customer"},
		"IDToken1": {username},
		"IDToken2": {password},
	}
	_, err := client.PostForm(loginURL, form)
	if err != nil { panic(err) }
}

func AutoLogin() {
	_, err := client.Get(autoLoginURL)
	if err != nil { panic(err) }
}

func createClient() http.Client {
	jar, _ := cookiejar.New(nil)
	client := http.Client{Jar: jar}
	return client
}

var client = createClient()

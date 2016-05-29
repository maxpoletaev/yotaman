package selfcare

import (
    "net/url"
    "net/http"
    "net/http/cookiejar"
)

const (
    CHANGE_OFFER_URL string = "https://my.yota.ru/selfcare/devices/changeOffer"
    DEVICES_URL string = "https://my.yota.ru/selfcare/devices"
    LOGIN_URL string = "https://login.yota.ru/UI/Login"
    AUTO_LOGIN_URL string = "https://my.yota.ru/selfcare/mydevices"
    LOGIN_SUCCESS_URL string = "https://my.yota.ru/selfcare/loginSuccess"
    LOGIN_ERROR_URL string = "https://my.yota.ru:443/selfcare/loginError"
)

func Login(username string, password string) {
    _, err := client.PostForm(LOGIN_URL, url.Values{
        "gotoOnFail": {LOGIN_ERROR_URL},
        "goto": {LOGIN_SUCCESS_URL},
        "ForceAuth": {"true"},
        "org": {"customer"},
        "IDToken1": {username},
        "IDToken2": {password},
    })
    if err != nil { panic(err) }
}

func AutoLogin() {
    _, err := client.Get(AUTO_LOGIN_URL)
    if err != nil { panic(err) }
}

func CreateClient() http.Client {
    jar, _ := cookiejar.New(nil)
    client := http.Client{Jar: jar}
    return client
}

var client http.Client = CreateClient()

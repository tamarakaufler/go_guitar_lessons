package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	//	"time"
)

const recaptchaURL string = "https://www.google.com/recaptcha/api/siteverify"

// CaptchaAuth ... object holding Google reCaptcha credentials
type CaptchaAuth struct {
	Secret string
}

type captchaResp struct {
	Success bool
	//	ChallengeTs time.Time `json:"challenge_ts,omitempty"` // timestamp of the challenge load (ISO format yyyy-MM-dd'T'HH:mm:ssZZ)
	//	Hostname    string    `json:"omitempty"`              // the hostname of the site where the reCAPTCHA was solved
	ErrorCodes []string `json:"error-codes,omitempty"`
}

// Validate ... validates recaptcha token sent on form submission as an extra form field
func (c CaptchaAuth) Validate(token string) (bool, error) {

	vdata := url.Values{"secret": {c.Secret}, "response": {token}}

	client := &http.Client{}
	resp, err := client.PostForm(recaptchaURL, vdata)

	defer resp.Body.Close()

	if err != nil {
		fmt.Printf("%v", err)
		return false, err
	}

	body, _ := ioutil.ReadAll(resp.Body)

	recaptchaResp := new(captchaResp)
	err = json.Unmarshal(body, recaptchaResp)

	if err != nil {
		return false, err
	}

	//fmt.Printf("recaptcha validation response: %+v \n", recaptchaResp)

	return recaptchaResp.Success, err
}

package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

var emailHost string = os.Getenv("SMTP_host")
var emailPort string = os.Getenv("SMTP_port")
var emailPass string = os.Getenv("SMTP_pass")

var emailUser string = "xxxxx@btinternet.com"
var emailFrom string = "noreply@guitar-lessons.co.uk"
var emailTo string = "xxxxx@gmail.com"

var reCaptchaKey string = os.Getenv("GOOGLE_RECAPTCHA_KEY")
var reCaptchaSecret string = os.Getenv("GOOGLE_RECAPTCHA_SECRET")
var gapiKey string = os.Getenv("GOOGLE_API_KEY")

type templateHandler struct {
	once     sync.Once
	templ    *template.Template
	filename string
}

type IntroData struct {
	Message string
	Error   bool
}

type AuthData struct {
	ReCaptchaKey string
	GapiKey      string
}

var keyAuth = AuthData{ReCaptchaKey: reCaptchaKey, GapiKey: gapiKey}
var recaptchaAuth = CaptchaAuth{Secret: reCaptchaSecret}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})

	t.templ.Execute(w, keyAuth)
}

func main() {

	// Sanity checks
	//----------------------------------------

	if reCaptchaKey == "" {
		log.Fatal("ERROR: email authentication (password) must be supplied for GOOGLE_RECAPTCHA_KEY")
	}

	if gapiKey == "" {
		log.Fatal("ERROR: email authentication (password) must be supplied for GOOGLE_API_KEY")
	}

	if emailPass == "" {
		log.Fatal("ERROR: email authentication (password) must be supplied for SMTP BTInternet server and xxxxx@btinternet.com")
	}
	if emailHost == "" {
		emailHost = "mail.btinternet.com"
	}
	if emailPort == "" {
		emailPort = "465"
	}

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Handlers
	//----------------------------------------

	// show the main page
	http.Handle("/", &templateHandler{filename: "intro.html"})

	// process form submission
	http.HandleFunc("/contact", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Email server check
		//----------------------------------------

		email := &Email{}

		err := email.Init(Auth{Host: emailHost, Port: emailPort, User: emailUser, Pass: emailPass})
		if err != nil {
			log.Panic("Wrong email authorization: ", err.Error())
			return
		}

		e := r.FormValue("email")

		if e == "" {
			failureMsg := IntroData{Message: "Redirect", Error: true}
			jsonData, err := json.Marshal(failureMsg)

			if err != nil {
				fmt.Printf("contact handler error: %v \n\n", err)
			}
			w.Write(jsonData)
			return
		}

		// reCaptcha validation
		//----------------------------------------

		recaptchaToken := r.Form.Get("g-recaptcha-response")
		isHuman, err := recaptchaAuth.Validate(recaptchaToken)
		if err != nil || !isHuman {

			failureMsg := IntroData{Message: "Redirect", Error: true}
			jsonData, err := json.Marshal(failureMsg)

			if err != nil {
				fmt.Printf("contact handler error: %v \n\n", err)
			}
			w.Write(jsonData)
			return
		}

		message, err := createEmailMessage(emailFrom, emailTo, r)
		if err != nil {
			failureMsg := IntroData{Message: "Form was not successfully submitted. Please try again.", Error: true}
			jsonData, _ := json.Marshal(failureMsg)
			w.Write(jsonData)
			return
		}

		err = email.Send(message)
		if err != nil {
			failureMsg := IntroData{Message: "Form was not successfully submitted. Please try again.", Error: true}
			jsonData, _ := json.Marshal(failureMsg)
			w.Write(jsonData)
			return
		}

		successMsg := IntroData{Message: "Thank you. Form was successfully submitted. I shall contact you shortly.", Error: false}
		jsonData, err := json.Marshal(successMsg)

		if err != nil {
			fmt.Printf("contact handler error: %v \n\n", err)
		}
		w.Write(jsonData)

	})

	// ------- WEB SERVER -------
	fmt.Println("Starting web server ...")
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal("Server problems: ", err)
	}
}

// createEmailMessage ... creates complete email message to be sent across the wire
func createEmailMessage(from string, to string, r *http.Request) (Message, error) {
	err := r.ParseForm()

	if err != nil {
		log.Println(err)

		return Message{}, err
	}

	fmt.Printf(">>> %+v\n\n", r.Form)

	firstName := r.Form.Get("first_name")
	lastName := r.Form.Get("last_name")
	email := r.Form.Get("note")
	note := r.Form.Get("note")

	fmt.Printf("\t\tForm data : firstName=%s, lastName=%s,email=%s,note=%s\n", firstName, lastName, email, note)
	subject := fmt.Sprintf("Guitar lessons interest (%s %s - %s)", firstName, lastName, email)

	m := "First name: %s \r\nLast name: %s \r\n\r\nEmail: %s \r\n\r\nNote: %s\r\n"
	body := fmt.Sprintf(m, firstName, lastName, email, note)

	message := Message{
		From:    from,
		To:      to,
		Subject: subject,
		Body:    body,
	}
	return message, nil
}

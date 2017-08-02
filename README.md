# go_guitar_lessons
Website for advertising guitar tuition

## Synopsis

The website has one main page, with tuition information, a Google map and a form for sending
contact information. The form is protected by captcha to avoid spamming.

## Technical details

- Go codebase
- email module
- captcha module (Google reCaptcha v2)
- javascript: jQuery
- CSS: Bootstrap
- Google maps

The sensitive information is not stored in the codebase, but provided on the command line, when
starting the application.

### Sensitive information

- SMTP_pass (email password to be able to send emails)l
- GOOGLE_RECAPTCHA_KEY (allows using Google reCaptcha)
- GOOGLE_RECAPTCHA_SECRET(allows using Google reCaptcha)
- GOOGLE_API_KEY (allows using various Google APIs, including Google maps)  


## Usage

a) change the email authentication detailsfirst:
b) change the email recipient
c) go build -o guitar_lessons
d) Optional: change the server port

### CAVEAT
    if port < 1024 => the application needs to be started with root privileges or, if running as a normal user, the following command needs
    to be run first:
        sudo setcap CAP_NET_BIND_SERVICE=+eip /path/to/program

e) export SMTP_pass=xxxxxxx && export GOOGLE_RECAPTCHA_KEY=yyyyyy && export GOOGLE_API_KEY=zzzzzz && export GOOGLE_RECAPTCHA_SECRET=qqqqqq && ./guitar_lessons > /dev/null 2>&1 &

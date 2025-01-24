package src

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"gopkg.in/gomail.v2"
)

type SubmitForm struct {
	name    string
	mail    string
	message string
}

var contactTemplate = `
<link rel="stylesheet" href="../static/styles/contact.css" />
<h1>Contact Me</h1>
<form id="contactForm" class="contactForm" hx-post="/submit" hx-swap="innerHTML" hx-target="#content">
	<label for="name">Name:</label>
	<input type="text" id="name" name="name" value="">
	<br>
	<label for="email">Email:</label>
	<input type="email" id="mail" name="mail" value="">
	<br>
	<label for="message">Message:</label>
	<textarea id="message" name="message"></textarea>
	<br>
	<input class="contactButton" type="submit" value="Submit">
</form>
`

var contactTemplateSubmitted = `
<h1>Contact Me</h1>
<p>Thank you for submitting the form!</p>
<p>We will answer you in the most briefly delay</p>
`

var contactTemplateFalse = `
<h1>Contact Me</h1>
<p class="contactError"> Something went wrong please try again !</p>
<form id="contactForm" class="contactForm" hx-post="/submit" hx-swap="innerHTML" hx-target="#content">
	<label for="name">Name:</label>
	<input type="text" id="name" name="name" value="">
	<br>
	<label for="email">Email:</label>
	<input type="email" id="mail" name="mail" value="">
	<br>
	<label for="message">Message:</label>
	<textarea id="message" name="message"></textarea>
	<br>
	<input class="contactButton" type="submit" value="Submit">
</form>
`

func ContactSend(f SubmitForm) bool {
	mail := os.Getenv("MAIL")
	pass := os.Getenv("PASS_MAIL")

	m := gomail.NewMessage()
	m.SetHeader("From", mail)
	m.SetHeader("To", mail)
	m.SetHeader("Subject", "Gomail test subject")
	m.SetBody("text/plain", fmt.Sprintf("Name: %s\nMail: %s\n\nMessage: %s\n", f.name, f.mail, f.message))
	d := gomail.NewDialer("smtp.gmail.com", 587, mail, pass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func ContactHandler(w http.ResponseWriter, r *http.Request) string {
	return contactTemplate
}

func HandleSubmit(w http.ResponseWriter, r *http.Request) string {
	if r.Method == http.MethodPost {

		form := SubmitForm{
			name:    r.FormValue("name"),
			mail:    r.FormValue("mail"),
			message: r.FormValue("message"),
		}
		w.Header().Set("Content-Type", "text/html")
		if ContactSend(form) {
			return contactTemplateSubmitted
		} else {
			return contactTemplateFalse
		}
		
	}
	http.Redirect(w, r, "/contact", http.StatusSeeOther)
	return ""
}

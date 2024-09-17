package src

import (
	"crypto/tls"
	"fmt"
	"os"

	"gopkg.in/gomail.v2"
)

type SubmitForm struct{
	name string
	mail string
	message string
}

var contactTemplate = `
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
<button class="contactButton" hx-get="/about" hx-swap="innerHTML" hx-target="#content" >Go back to homepage</button>
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

func ContactSend(f SubmitForm)bool {
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


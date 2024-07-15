package src

type SubmitForm struct{
	name string
	mail string
	message string
}

var contactTemplate = `
<h1>Contact Us</h1>
<form id="contactForm" class="contactForm" hx-post="/submit" hx-swap="innerHTML" hx-target="#content">
	<label for="name">Name:</label>
	<input type="text" id="name" name="name" value="">
	<br>
	<label for="email">Email:</label>
	<input type="email" id="email" name="email" value="">
	<br>
	<label for="message">Message:</label>
	<textarea id="message" name="message"></textarea>
	<br>
	<input type="submit" value="Submit">
</form>
`

var contactTemplateSubmitted = `
<h1>Contact Us</h1>
<p>Thank you for submitting the form!</p>
<form id="contactForm" class="contactForm" hx-post="/submit" hx-swap="innerHTML" hx-target="#content">
	<label for="name">Name:</label>
	<input type="text" id="name" name="name" value="">
	<br>
	<label for="email">Email:</label>
	<input type="email" id="email" name="email" value="">
	<br>
	<label for="message">Message:</label>
	<textarea id="message" name="message"></textarea>
	<br>
	<input type="submit" value="Submit">
</form>
`

var contactTemplateFalse = `
<h1>Contact Us</h1>
<p>ALERT ALERT!</p>
<form id="contactForm" class="contactForm" hx-post="/submit" hx-swap="innerHTML" hx-target="#content">
	<label for="name">Name:</label>
	<input type="text" id="name" name="name" value="">
	<br>
	<label for="email">Email:</label>
	<input type="email" id="email" name="email" value="">
	<br>
	<label for="message">Message:</label>
	<textarea id="message" name="message"></textarea>
	<br>
	<input type="submit" value="Submit">
</form>
`

func ContactSend(f SubmitForm)bool{
	return false
}
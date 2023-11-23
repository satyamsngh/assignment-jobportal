package otputil

import (
	"fmt"
	"math/rand"
	"net/smtp"
	"time"
)

func (o *Otp) GenerateOtp(email string) string {
	// Sender's email address and password
	from := "your-email@gmail.com"
	password := "your-email-password"

	// Recipient's email address
	to := email

	// SMTP server details for Gmail
	smtpServer := "smtp.gmail.com"
	smtpPort := 587

	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Generate a random six-digit number
	randomNumber := rand.Intn(900000) + 100000

	// Message content
	message := []byte(fmt.Sprintf("Subject: Your Verification Code\n\nYour verification code is: %06d", randomNumber))

	// Authentication information
	auth := smtp.PlainAuth("", from, password, smtpServer)

	// SMTP connection
	smtpAddr := fmt.Sprintf("%s:%d", smtpServer, smtpPort)

	// Send email
	err := smtp.SendMail(smtpAddr, auth, from, []string{to}, message)
	if err != nil {
		fmt.Println("Error sending email:", err)
		return ""
	}

	fmt.Println("Email sent successfully!")
	return fmt.Sprintf("%06d", randomNumber)
}

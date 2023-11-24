package otputil

import (
	"fmt"
	"math/rand"
	"net/smtp"
	"time"
)

func (o Otp) GenerateOtp(email string) (string, error) {
	// Sender's email address and password
	from := "satyamdg18577@gmail.com"
	password := "ixlj zbcm fgwk tjta"

	// Recipient's email address
	fmt.Println("++++++++++++++++++++++", email)
	to := email

	// SMTP server details for Gmail
	smtpServer := "smtp.gmail.com"
	fmt.Print("==============", o.Rd)
	smtpPort := o.Rd

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
		return "", fmt.Errorf("failed to send OTP email: %v", err)
	}

	fmt.Println("Email sent successfully!")
	return fmt.Sprintf("%06d", randomNumber), nil
}

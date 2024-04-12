package main

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
	"time"

	"github.com/google/uuid"
)

func SendVerificationMail(email string) (string, error) {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	from := os.Getenv("SMTP_MAIL")
	password := os.Getenv("SMTP_PASSWORD")

	subject := "email verification by rohanyh"
	//	example@example.com		EXAMPLE_PASSWORD	smtp.example.com
	auth := smtp.PlainAuth(
		"",
		from,
		password,
		host,
	)

	verificationKey := fmt.Sprintf(
		"%x", sha256.Sum256([]byte(email + "-" + uuid.New().String())[:]),
	)

	verificationLinkBase := "http://localhost:3000/verify-email?token="

	body := fmt.Sprintf(`
 	<html>
 	<a href="%v%v" target="_blank">CLICK</a>
 	<p>This link will expire in 24 hours.</p>
 	</html>
 	`, verificationLinkBase, verificationKey)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	client, err := smtp.Dial(host + ":" + port)
	if err != nil {
		fmt.Println("error while connecting to smtp...")
		return "", err
	}

	client.StartTLS(tlsConfig)

	if err := client.Auth(auth); err != nil {
		fmt.Println("error while authenticating to SMTP...")
		return "", err
	}

	if err := client.Mail(from); err != nil {
		fmt.Println("error while initiating mail transaction...")
		return "", err
	}

	if err := client.Rcpt(email); err != nil {
		fmt.Println("error while issuing RCPT(CMD) to SMTP...")
		return "", err
	}

	w, err := client.Data()
	if err != nil {
		fmt.Println("error while issuing DATA(CMD) to SMTP...")
		return "", err
	}

	_, err = w.Write([]byte(
		fmt.Sprintf("MIME-Version: %v\r\n", "1.0") +
			fmt.Sprintf("Content-type: %v\r\n", "text/html; charset=UTF-8") +
			fmt.Sprintf("From: %v\r\n", from) +
			fmt.Sprintf("To: %v\r\n", email) +
			fmt.Sprintf("Subject: %v\r\n", subject) +
			fmt.Sprintf("%v\r\n", body),
	))

	if err != nil {
		fmt.Println("error while writing data to the header/body")
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	statusCMD := redisClient.Set(ctx, verificationKey, fmt.Sprintf("%v", email), time.Duration(time.Hour*24))
	if statusCMD.Err() != nil {
		return "", statusCMD.Err()
	}

	err = w.Close()
	if err != nil {
		fmt.Println("error closing writer")
		return "", err
	}

	client.Quit()
	return fmt.Sprintf("%v%v", verificationLinkBase, verificationKey), nil
}

package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/sylank/lavender-commons-go/messaging"
)

// SendMailRequest ..
type SendMailRequest struct {
	ToAddress string `json:"to_address"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
}

// SendTransactionalMail ...
func SendTransactionalMail(emailAddress string, subject string, message string) error {
	mailReques := &SendMailRequest{ToAddress: emailAddress, Subject: subject, Body: message}
	jsonData, err := json.Marshal(mailReques)
	if err != nil {
		log.Println("Failed to marshall transactional email request")

		return err
	}

	queueName := os.Getenv("TRANSACTIONAL_EMAIL_QUEUE_NAME")
	return messaging.SendTransactionalEmail(string(jsonData), queueName)
}

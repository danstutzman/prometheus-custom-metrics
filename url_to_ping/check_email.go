package url_to_ping

import (
	pop3 "github.com/bytbox/go-pop3"
	"log"
	"net/mail"
	"strings"
	"time"
)

func main() {
	MailboxHasMailWithSubject("fake.alert.emails@gmail.com", "99require23provide",
		time.Now().Add(-1*time.Hour), "[FIRING:1] FakeAlertToVerifyEndToEnd")
}

func MailboxHasMailWithSubject(username string, password string, since time.Time,
	expectedSubject string) bool {

	log.Printf("Attempting DialTLS to pop.gmail.com:995...")
	client, err := pop3.DialTLS("pop.gmail.com:995")
	if err != nil {
		log.Fatalf("DialTLS failed: %s", err)
	}

	// Why put "recent:" in front of user name?
	// Otherwise we can only download messages from POP client one time
	// https://support.google.com/mail/answer/7104828?hl=en&authuser=1&ref_topic=3398031

	log.Printf("Sending USER...")
	if err = client.User("recent:" + username); err != nil {
		log.Fatalf("User failed: %s", err)
	}

	log.Printf("Sending PASS...")
	if err = client.Pass(password); err != nil {
		log.Fatalf("Pass failed: %s", err)
	}

	var msgNums []int
	foundExpectedSubject := false
	log.Printf("Sending LIST...")
	if msgNums, _, err = client.ListAll(); err != nil {
		log.Fatalf("ListAll failed: %s", err)
	}
	log.Printf("Listed %d message indices", len(msgNums))
	for _, msgNum := range msgNums {
		var text string
		log.Printf("Sending RETR for msg num %d...", msgNum)
		if text, err = client.Retr(msgNum); err != nil {
			log.Fatalf("Retr(msgNum=%d) failed: %s", msgNum, err)
		}

		// Add extra newline in case there's no body
		message, err := mail.ReadMessage(strings.NewReader(text + "\n"))
		if err != nil {
			log.Fatalf("Error from ReadMessage: %s", err)
		}
		date, err := message.Header.Date()
		if err != nil {
			log.Fatalf("Error from Date(): %s", err)
		}

		subject := message.Header.Get("Subject")
		if date.Before(since) {
			log.Printf("Deleting %s - %s", date, subject)
			if err := client.Dele(msgNum); err != nil {
				log.Fatalf("Tried to delete msgNum=%d but got %s", msgNum, err)
			}
		} else {
			if subject == expectedSubject {
				log.Printf("Found %s - %s", date, subject)
				foundExpectedSubject = true
			}
		}
	}

	log.Printf("Sending QUIT...")
	if err = client.Quit(); err != nil {
		log.Fatalf("Quit failed: %s", err)
	}

	return foundExpectedSubject
}

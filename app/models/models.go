package models

import (
	"fmt"
	"time"
)

var (
	// data about the outside contact
	// proposed relationship = a single contact can have many inmates but no other contacts
	Contact struct {
		firstName, lastName, email, addressLine1, addressLine2, city, state string
		id, zipCode, dayPhone, mobilePhone                                  int
	}

	// data about inmate
	// proposed relationship = a single inmate can have many contacts but no other inmates
	Inmate struct {
		firstName, lastName, correctionalAgency, institution, deliveryStatus string
		id, inmateID                                                         int
		emailAlert, block                                                    bool
	}

	// inbound and outbound messages between contact and inmate
	// proposed relationship = many inboxMessages can be associated with a single inboxBatch (id)
	inboxMessage struct {
		from, subject, message string
		id, messageID          int
		isUnreadMessage        bool
		sentTimeStamp          time.Time
	}

	// data about batch process (crawler service worker)
	// proposed relationship = a single inboxBatch can have many inboxMessages
	inboxBatch struct {
		status, errorMessage         string
		id                           int
		batchError                   bool
		startTimeStamp, endTimeStamp time.Time
	}
)

func main() {
	fmt.Println("vim-go")
}

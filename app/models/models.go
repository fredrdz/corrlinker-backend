package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	// data about the outside contact
	// proposed relationship = a single contact can have many inmates but no other contacts
	Contact struct {
		ID                                                                  primitive.ObjectID `json:"id" bson:"_id,omitempty"`
		FirstName, LastName, Email, AddressLine1, AddressLine2, City, State string
		ZipCode, DayPhone, MobilePhone                                      uint64
	}

	// data about inmate
	// proposed relationship = a single inmate can have many contacts but no other inmates
	Inmate struct {
		ID                                                                   primitive.ObjectID `json:"id" bson:"_id,omitempty"`
		FirstName, LastName, CorrectionalAgency, Institution, DeliveryStatus string
		InmateID                                                             uint64
		EmailAlert, Block                                                    bool
	}

	// inbound and outbound messages between contact and inmate
	// proposed relationship = many inboxMessages can be associated with a single inboxBatch (id)
	InboxMessage struct {
		ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
		From            string             `json:"from" bson:"from,omitempty"`
		Subject         string             `json:"subject" bson:"subject,omitempty"`
		Message         string             `json:"message" bson:"message,omitempty"`
		MessageID       uint64             `json:"message_id" bson:"message_id,omitempty"`
		IsUnreadMessage bool               `json:"is_unread_message" bson:"is_unread_message,omitempty"`
		SentDate        string             `json:"sent_date" bson:"sent_date,omitempty"` // temporary, needs conversion to time.Time
	}

	// data about batch process (crawler service worker)
	// proposed relationship = a single inboxBatch can have many inboxMessages
	InboxBatch struct {
		ID                           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
		InboxMessage                 primitive.ObjectID `json:"inbox_message" bson:"inbox_message,omitempty"`
		Status, ErrorMessage         string
		BatchError                   bool
		StartTimeStamp, EndTimeStamp time.Time
	}
)

package models

import (
	"context"
	"errors"
	"go-poke-app/pkg/database"
	"go-poke-app/pkg/email"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Poke struct {
	ID            primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	SenderName    string             `json:"SenderName" bson:"sender_name,omitempty"`
	Message       string             `json:"Message" bson:"message,omitempty"`
	ReceiverEmail string             `json:"ReceiverEmail" bson:"receiver_email,omitempty"`
	ReceiverName  string             `json:"ReceiverName" bson:"receiver_name,omitempty"`
	IsPublic      bool               `json:"IsPublic" bson:"is_public,omitempty"`
	DateSent      time.Time          `json:"DateSent" bson:"date_sent,omitempty"`
}

type Pokes struct {
	List []Poke
}

// create poke
func (p *Poke) Create() (err error) {
	if p.ReceiverEmail == "" {
		err = errors.New("poke needs email")
		return
	}

	if p.SenderName == "" {
		p.SenderName = "Anonymous"
	}
	if p.ReceiverName == "" {
		p.ReceiverName = "Anonymous"
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	collection := database.Collection()
	_, err = collection.InsertOne(ctx, p)
	return err
}

// send email
func (p *Poke) SendEmail() (err error) {
	email.Send(email.Email{
		SenderName:    p.SenderName,
		Message:       p.Message,
		ReceiverEmail: p.ReceiverEmail,
		ReceiverName:  p.ReceiverName,
	})

	return err
}

// Get 10 most recent pokes
func GetPokes() *Pokes {

	pokes := &Pokes{}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"date_sent", -1}})
	findOptions.SetLimit(10)
	filter := bson.D{{"is_public", true}}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	collection := database.Collection()
	cur, err := collection.Find(ctx, filter, findOptions)
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var poke Poke
		err = cur.Decode(&poke)
		if err != nil {
			log.Fatal(err)
		}

		pokes.List = append(pokes.List, poke)
	}

	return pokes
}

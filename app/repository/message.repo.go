package repository

import (
	"lift-fitness-gym/app/db"
	"lift-fitness-gym/app/model"

	"github.com/jmoiron/sqlx"
)


type Message struct {
	db * sqlx.DB
}

func NewMessageRepository () Message{
	return Message{
		db: db.GetConnection(),
	}
}
func (repo * Message)NewMessage(contactUs model.ContactUs) error {
	_, err := repo.db.Exec("INSERT INTO message(name, value , email)  VALUES (?, ? , ?)", contactUs.Name, contactUs.Message, contactUs.Email)
	return err
}

func(repo * Message) GetMessages()([]model.ContactUs, error) {
	messages := make([]model.ContactUs, 0)
	err := repo.db.Select(&messages, "SELECT id, value, email, created_at, name from message order by created_at desc")
	return messages,err
}
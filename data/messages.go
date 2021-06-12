package data

import (
	"fmt"
	"github.com/mdev5000/qlik_message/postgres"
)

const Schema = `
create table if not exists messages (
	id serial,
	version integer not null,
    message text not null
);
`

const messageDbTableName = "messages"

type IdMissingError struct {
	Type string
	Id   int64
}

func (e IdMissingError) Error() string {
	return fmt.Sprintf("no rows in result set for %s with id %d", e.Type, e.Id)
}

func (e IdMissingError) Is(target error) bool {
	switch target.(type) {
	case IdMissingError:
		return true
	default:
		return false
	}
}

type MessageId = int64
type MessageVersion = int

type Message struct {
	Id      MessageId      `sharedDbInstance:"id"`
	Version MessageVersion `sharedDbInstance:"version"`
	Message string         `sharedDbInstance:"message"`
}

type CreateMessage struct {
	Message string `sharedDbInstance:"message"`
}

type MessageRepository struct {
	db *postgres.DB
}

func NewMessageRepository(db *postgres.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (mr *MessageRepository) DeleteById(id MessageId) error {
	r, err := mr.db.Exec(`delete from messages where id = $1`, id)
	if err != nil {
		return err
	}
	affected, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		return fmt.Errorf("expected 1 row to be deleted by id but was %d", affected)
	}
	return nil
}

func (mr *MessageRepository) Create(cm CreateMessage) (MessageId, error) {
	rows, err := mr.db.Query(`insert into messages (version, message) values (1, $1) returning id`, cm.Message)
	if err != nil {
		return 0, err
	}
	if !rows.Next() {
		return 0, fmt.Errorf("expected 1 creation rows, but none present")
	}
	var id MessageId
	err = rows.Scan(&id)
	if err != nil {
		return id, err
	}
	numRow := 1
	for rows.Next() {
		numRow += 1
	}
	if numRow != 1 {
		return id, fmt.Errorf("unexpected number of rows expected %d, but was %d", 1, numRow)
	}
	return id, nil
}

func (mr *MessageRepository) GetAll(messages *[]*Message) error {
	return mr.db.Select(messages, `select id, version, message from messages`)
}

func (mr *MessageRepository) GetById(id MessageId, m *Message) error {
	if err := mr.db.Get(m, `select id, version, message from messages where id=$1`, id); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return IdMissingError{
				Type: messageDbTableName,
				Id:   id,
			}
		}
		return err
	}
	return nil
}

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

const RepositoryIdentifierMessages = "messages"

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

type MessagesRepository struct {
	db *postgres.DB
}

func NewMessageRepository(db *postgres.DB) *MessagesRepository {
	return &MessagesRepository{db: db}
}

func (mr *MessagesRepository) DeleteById(id MessageId) error {
	r, err := mr.db.Exec(`delete from messages where id = $1`, id)
	if err != nil {
		return err
	}
	affected, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		return repoError(
			RepositoryIdentifierMessages,
			fmt.Errorf("expected delete by id to delete 1 row but %d were deleted", affected))
	}
	return nil
}

func (mr *MessagesRepository) Create(cm CreateMessage) (MessageId, error) {
	rows, err := mr.db.Query(`insert into messages (version, message) values (1, $1) returning id`, cm.Message)
	if err != nil {
		return 0, err
	}
	if !rows.Next() {
		return 0, repoError(RepositoryIdentifierMessages,
			fmt.Errorf("create message expected 1 row returned by was 0"))
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
		return id, repoError(RepositoryIdentifierMessages,
			fmt.Errorf("unexpected number of rows expected %d, but was %d", 1, numRow))
	}
	return id, nil
}

func (mr *MessagesRepository) GetAll(messages *[]*Message) error {
	return mr.db.Select(messages, `select id, version, message from messages`)
}

func (mr *MessagesRepository) GetById(id MessageId, m *Message) error {
	if err := mr.db.Get(m, `select id, version, message from messages where id=$1`, id); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return idMissingError(RepositoryIdentifierMessages, id)
		}
		return err
	}
	return nil
}

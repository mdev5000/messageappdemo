package data

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/mdev5000/qlik_message/postgres"
	"github.com/pkg/errors"
	"time"
)

const RepositoryIdentifierMessages = "messages"

var NoChangesInUpdateError = errors.New("no changes made in update call")

type MessageId = int64
type MessageVersion = int

type Message struct {
	Id        MessageId      `db:"id"`
	Version   MessageVersion `db:"version"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
	Message   string         `db:"message"`
}

type CreateMessage struct {
	Message   string    `db:"message"`
	CreatedAt time.Time `db:"created_at"`
}

type MessagesRepository struct {
	db *postgres.DB
}

func NewMessageRepository(db *postgres.DB) *MessagesRepository {
	return &MessagesRepository{db: db}
}

const repoName = "MessagesRepository"

func (mr *MessagesRepository) DeleteById(id MessageId) error {
	const op = repoName + ".DeleteById"
	r, err := mr.db.Exec(`delete from messages where id = $1`, id)
	if err != nil {
		return repoError(op, fmt.Errorf("failed to delete message with id %d: \n%w", id, err), err)
	}
	affected, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		return repoError2(op, fmt.Errorf("expected delete by id to delete 1 row but %d were deleted", affected))
	}
	return nil
}

// Create creates a new message. Note that CreatedAt should be in the UTC-0 timezone.
func (mr *MessagesRepository) Create(cm CreateMessage) (MessageId, error) {
	const op = repoName + ".Create"
	rows, err := mr.db.Query(
		`
insert into messages (version, created_at, updated_at, message)
values (1, $1, $1, $2) returning id
`, cm.CreatedAt, cm.Message)
	if err != nil {
		return 0, repoError(op, fmt.Errorf("failed to create message: %w", err), err)
	}

	if !rows.Next() {
		return 0, repoError2(op, fmt.Errorf("create message expected 1 row returned by was 0"))
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
		return id, repoError2(op, fmt.Errorf("unexpected number of rows expected %d, but was %d", 1, numRow))
	}
	return id, nil
}

func (mr *MessagesRepository) GetAll(messages *[]*Message) error {
	return mr.db.Select(messages, `select id, version, created_at, updated_at, message from messages`)
}

func (mr *MessagesRepository) GetById(id MessageId, m *Message) error {
	const op = repoName + ".GetById"
	if err := mr.db.Get(m, `select id, version, created_at, updated_at, message from messages where id=$1`, id); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return repoError2(op, idMissingError(RepositoryIdentifierMessages, id))
		}
		return err
	}
	return nil
}

func (mr *MessagesRepository) UpdateById(id MessageId, m Message) (MessageVersion, error) {
	const op = repoName + ".UpdateById"
	q := sq.Update("messages").
		PlaceholderFormat(sq.Dollar).
		Set("version", sq.Expr("version + 1")).
		Set("updated_at", NowUTC()).
		Where(sq.Eq{"id": id})

	hasSetField := false
	if m.Message != "" {
		hasSetField = true
		q = q.Set("message", m.Message)
	}
	// No changes actually made, so don't do anything.
	if !hasSetField {
		return 0, repoError2(op, NoChangesInUpdateError)
	}
	sql, args, err := q.ToSql()
	if err != nil {
		return 0, repoError(op, fmt.Errorf("failed to generate update query: %w", err), err)
	}
	row := mr.db.QueryRow(sql+" returning version", args...)
	if err := row.Err(); err != nil {
		return 0, repoError(op, fmt.Errorf("failed to update row: %w", err), err)
	}
	var version MessageVersion
	if err := row.Scan(&version); err != nil {
		return 0, repoError(op, fmt.Errorf("failed to scan version number: %w", err), err)
	}
	return version, nil
}

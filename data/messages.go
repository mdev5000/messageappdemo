package data

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/mdev5000/messageappdemo/apperrors"
	"github.com/mdev5000/messageappdemo/messages"
	"github.com/mdev5000/messageappdemo/postgres"
	errors2 "github.com/pkg/errors"
)

type MessageId = messages.MessageId
type MessageVersion = messages.MessageVersion
type Message = messages.Message
type CreateMessage = messages.CreateMessage
type ModifyMessage = messages.ModifyMessage
type MessageQuery = messages.MessageQuery

// MessagesRepository is the repository implementation for the messages Repository.
type MessagesRepository struct {
	db *postgres.DB
}

func NewMessageRepository(db *postgres.DB) *MessagesRepository {
	return &MessagesRepository{db: db}
}

const repoName = "MessagesRepository"

// Map of fields the user is allowed to query in format {field: table_col}
var queryableFields = map[string]string{
	messages.FieldId:        "id",
	messages.FieldVersion:   "version",
	messages.FieldCreatedAt: "created_at",
	messages.FieldUpdatedAt: "updated_at",
	messages.FieldMessage:   "message",
}

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
		return repoError2(op, idMissingError(op, id))
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
	const op = repoName + ".GetAll"
	if err := mr.db.Select(messages,
		`select id, version, created_at, updated_at, message from messages`); err != nil {
		return repoError(op, fmt.Errorf("failed to get messages: %w", err), err)
	}
	return nil
}

func (mr *MessagesRepository) GetAllQuery(query MessageQuery, messages *[]*Message) error {
	const op = repoName + ".GetAllQuery"

	var cols []string
	if len(query.Fields) == 0 {
		cols = make([]string, 0, len(queryableFields))
		for _, col := range queryableFields {
			cols = append(cols, col)
		}
	} else {
		cols = make([]string, 0, len(query.Fields))
		var notFound []string
		for field := range query.Fields {
			col, found := queryableFields[field]
			if found {
				cols = append(cols, col)
			} else {
				notFound = append(notFound, field)
			}
		}
		if len(notFound) != 0 {
			err := fmt.Errorf("invalid messages fields: %s", strings.Join(notFound, ", "))
			aErr := apperrors.Error{
				EType: apperrors.ETInvalid,
				Op:    op,
				Err:   err,
				Stack: errors2.WithStack(err),
			}
			aErr.AddResponse(apperrors.ErrorResponse(err.Error()))
			return &aErr
		}
	}

	q := sq.Select(cols...).From("messages")
	if query.Limit > 0 {
		q = q.Limit(query.Limit)
	}
	if query.Offset > 0 {
		q = q.Offset(query.Offset)
	}

	sqlS, args, err := q.ToSql()
	if err != nil {
		return repoError(op, fmt.Errorf("failed to generate messages query:\n%w", err), err)
	}
	if err := mr.db.Select(messages, sqlS, args...); err != nil {
		return repoError(op,
			fmt.Errorf("failed to run messages query\nquery: %s\nargs: %+v\n%w", sqlS, args, err), err)
	}
	return nil
}

func (mr *MessagesRepository) GetById(id MessageId, m *Message) error {
	const op = repoName + ".GetById"
	if err := mr.db.Get(m, `select id, version, created_at, updated_at, message from messages where id=$1`, id); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return repoError2(op, idMissingError(op, id))
		}
		return err
	}
	return nil
}

func (mr *MessagesRepository) UpdateById(id MessageId, m ModifyMessage) (MessageVersion, error) {
	const op = repoName + ".UpdateById"

	q := sq.Update("messages").
		PlaceholderFormat(sq.Dollar).
		Set("version", sq.Expr("version + 1")).
		Set("updated_at", nowUTC()).
		Where(sq.Eq{"id": id})

	q = q.Set("message", m.Message)

	sqlS, args, err := q.ToSql()
	if err != nil {
		return 0, repoError(op, fmt.Errorf("failed to generate update query: %w", err), err)
	}

	row := mr.db.QueryRow(sqlS+" returning version", args...)
	if err := row.Err(); err != nil {
		return 0, repoError(op, fmt.Errorf("failed to update row: %w", err), err)
	}

	var version MessageVersion
	if err := row.Scan(&version); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, repoError2(op, idMissingError(op, id))
		}
		return 0, repoError(op, fmt.Errorf("failed to scan version number: %w", err), err)
	}
	return version, nil
}

func nowUTC() time.Time {
	loc, err := time.LoadLocation("UTC")
	if err != nil {
		panic(fmt.Errorf("failed to load UTC timezone: %w", err))
	}
	return time.Now().In(loc).Round(time.Millisecond)
}

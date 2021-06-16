package data

import (
	"errors"
	"github.com/mdev5000/qlik_message/messages"
	"github.com/mdev5000/qlik_message/postgres"
	"github.com/stretchr/testify/require"
	"sort"
	"testing"
)

func tMessageRepository(db *postgres.DB) *MessagesRepository {
	return NewMessageRepository(db)
}

func TestMessageRepository_Create_canCreateNewMessages(t *testing.T) {
	db, closeDb := acquireDb()
	defer closeDb()
	mr := tMessageRepository(db)

	now := nowUTC()

	id, err := mr.Create(CreateMessage{
		Message:   "my message",
		CreatedAt: now,
	})
	require.NoError(t, err)

	var m Message
	require.NoError(t, mr.GetById(id, &m))
	require.Equal(t, 1, m.Version)
	require.Equal(t, "my message", m.Message)
	require.True(t, now.Equal(m.CreatedAt))
	require.True(t, now.Equal(m.UpdatedAt))
}

func TestMessageRepository_GetById(t *testing.T) {
	db, closeDb := acquireDb()
	defer closeDb()
	mr := tMessageRepository(db)

	// Add more records to make sure it's actually returning the correct one by id
	var err error
	_, err = mr.Create(CreateMessage{Message: "first"})
	require.NoError(t, err)
	id, err := mr.Create(CreateMessage{Message: "find this one"})
	require.NoError(t, err)
	_, err = mr.Create(CreateMessage{Message: "first"})
	require.NoError(t, err)

	var m Message
	require.NoError(t, mr.GetById(id, &m))
	require.Equal(t, id, m.Id)
	require.Equal(t, 1, m.Version)
	require.Equal(t, "find this one", m.Message)
}

func TestMessageRepository_DeleteById_canDeleteMessagesById(t *testing.T) {
	db, closeDb := acquireDb()
	defer closeDb()
	mr := tMessageRepository(db)

	id, err := mr.Create(CreateMessage{Message: "my message"})
	require.NoError(t, err)

	require.NoError(t, mr.DeleteById(id))

	var m Message
	err = mr.GetById(id, &m)
	require.True(t, errors.Is(err, messages.IdMissingError{}))
	require.Error(t, err, "MessagesRepository.DeleteById: no rows in result for get by id with id %d", id)
}

func TestMessageRepository_DeleteById_onlyDeletesTheSpecifiedId(t *testing.T) {
	db, closeDb := acquireDb()
	defer closeDb()
	mr := tMessageRepository(db)

	id1, err := mr.Create(CreateMessage{Message: "message 1"})
	require.NoError(t, err)
	id2, err := mr.Create(CreateMessage{Message: "message 2"})
	require.NoError(t, err)
	id3, err := mr.Create(CreateMessage{Message: "message 3"})
	require.NoError(t, err)

	require.NoError(t, mr.DeleteById(id2))

	var messages []*Message
	require.NoError(t, mr.GetAll(&messages))

	require.Len(t, messages, 2)

	messageIds := make([]MessageId, len(messages))
	for i, message := range messages {
		messageIds[i] = message.Id
	}
	sort.Slice(messageIds, func(i, j int) bool { return i < j })

	require.Equal(t, []MessageId{id1, id3}, messageIds, "only ids of existing rows are still present")
}

func TestMessageRepository_DeleteById_returnsErrorWhenNoRowsAreDeleted(t *testing.T) {
	db, closeDb := acquireDb()
	defer closeDb()
	mr := tMessageRepository(db)

	err := mr.DeleteById(5)
	require.EqualError(t, err,
		"Error [internal] (MessagesRepository.DeleteById): MessagesRepository.DeleteById: no row in result with id 5")
}

func TestMessagesRepository_UpdateById(t *testing.T) {
	db, closeDb := acquireDb()
	defer closeDb()
	mr := tMessageRepository(db)

	// Add more records to make sure it's actually returning the correct one by id
	id, err := mr.Create(CreateMessage{Message: "first message"})
	require.NoError(t, err)

	var message Message
	require.NoError(t, mr.GetById(id, &message))

	v, err := mr.UpdateById(id, ModifyMessage{Message: "new message"})
	require.NoError(t, err)

	var messageChanged Message
	require.NoError(t, mr.GetById(id, &messageChanged))

	require.Equal(t, message.Version+1, v, "version is not the same as the previous")
	require.Equal(t, messageChanged.Version, v, "version has been updated")
	require.Equal(t, messageChanged.Message, "new message", "message is changed")
	require.True(t, message.UpdatedAt.Before(messageChanged.UpdatedAt), "updated_at has been updated")
}

func TestMessagesRepository_UpdateById_errorWhenMissing(t *testing.T) {
	db, closeDb := acquireDb()
	defer closeDb()
	mr := tMessageRepository(db)

	_, err := mr.UpdateById(5, ModifyMessage{Message: "new message"})

	require.Equal(t,
		idMissingError("MessagesRepository.UpdateById", 5),
		errors.Unwrap(err))
}

func TestMessagesRepository_UpdateById_whenNoChanges(t *testing.T) {
	db, closeDb := acquireDb()
	defer closeDb()
	mr := tMessageRepository(db)

	_, err := mr.UpdateById(5, ModifyMessage{Message: "new message"})

	require.Equal(t,
		idMissingError("MessagesRepository.UpdateById", 5),
		errors.Unwrap(err))
}

func TestMessagesRepository_GetAllQuery_canGetAllFields(t *testing.T) {
	db, closeDb := acquireDb()
	defer closeDb()
	mr := tMessageRepository(db)

	now := nowUTC()
	id, err := mr.Create(CreateMessage{Message: "first", CreatedAt: now})
	require.NoError(t, err)

	var messages []*Message
	require.NoError(t, mr.GetAllQuery(MessageQuery{}, &messages))

	require.Len(t, messages, 1)
	require.Equal(t, id, messages[0].Id)
	require.Equal(t, "first", messages[0].Message)
	require.True(t, now.Equal(messages[0].CreatedAt))
	require.True(t, now.Equal(messages[0].UpdatedAt))
}

func TestMessagesRepository_GetAllQuery_canLimitQueriedFields(t *testing.T) {
	db, closeDb := acquireDb()
	defer closeDb()
	mr := tMessageRepository(db)

	id, err := mr.Create(CreateMessage{Message: "a message"})
	require.NoError(t, err)

	var golden Message
	require.NoError(t, mr.GetById(id, &golden))

	cases := []struct {
		name     string
		fields   []string
		expected Message
	}{
		{name: "id", fields: []string{"id"}, expected: Message{Id: golden.Id}},
		{name: "version", fields: []string{"version"}, expected: Message{Version: golden.Version}},
		{name: "createdAt", fields: []string{"createdAt"}, expected: Message{CreatedAt: golden.CreatedAt}},
		{name: "updatedAt", fields: []string{"updatedAt"}, expected: Message{UpdatedAt: golden.UpdatedAt}},
		{name: "message", fields: []string{"message"}, expected: Message{Message: golden.Message}},
		{
			name:   "version, message",
			fields: []string{"version", "message"},
			expected: Message{
				Version: golden.Version,
				Message: golden.Message,
			}},
		{
			name:   "version, updatedAt, message",
			fields: []string{"version", "updatedAt", "message"},
			expected: Message{
				Version:   golden.Version,
				UpdatedAt: golden.UpdatedAt,
				Message:   golden.Message,
			}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			fields := map[string]struct{}{}
			for _, field := range c.fields {
				fields[field] = struct{}{}
			}
			var messages []*Message
			require.NoError(t, mr.GetAllQuery(MessageQuery{Fields: fields}, &messages))
			require.Len(t, messages, 1)
			require.Equal(t, &c.expected, messages[0])
		})
	}
}

func TestMessagesRepository_getAllQuery(t *testing.T) {
	db, closeDb := acquireDb()
	defer closeDb()
	mr := tMessageRepository(db)

	id1, err := mr.Create(CreateMessage{Message: "first"})
	require.NoError(t, err)
	id2, err := mr.Create(CreateMessage{Message: "second"})
	require.NoError(t, err)
	id3, err := mr.Create(CreateMessage{Message: "third"})
	require.NoError(t, err)

	t.Run("can retrieve all records", func(t *testing.T) {
		q := MessageQuery{Fields: map[string]struct{}{"id": {}}}
		var messages []*Message
		require.NoError(t, mr.GetAllQuery(q, &messages))

		require.Len(t, messages, 3)
		require.Equal(t, id1, messages[0].Id)
		require.Equal(t, id2, messages[1].Id)
		require.Equal(t, id3, messages[2].Id)

	})

	t.Run("can limited and offset values", func(t *testing.T) {
		q := MessageQuery{
			Fields: map[string]struct{}{"id": {}},
			Limit:  2,
			Offset: 1,
		}
		var messages []*Message
		require.NoError(t, mr.GetAllQuery(q, &messages))

		require.Len(t, messages, 2)
		require.Equal(t, id2, messages[0].Id)
		require.Equal(t, id3, messages[1].Id)
	})

	t.Run("bad offset returns empty", func(t *testing.T) {
		q := MessageQuery{
			Fields: map[string]struct{}{"id": {}},
			Offset: 500,
		}
		var messages []*Message
		require.NoError(t, mr.GetAllQuery(q, &messages))
		require.Len(t, messages, 0)
	})
}

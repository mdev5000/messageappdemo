package data

import (
	"errors"
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

	now := NowUTC()

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
	require.True(t, errors.Is(err, IdMissingError{}))
	require.Error(t, err, "messages repository: no rows in result for get by id with id %d", id)
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
		"Error [internal] (MessagesRepository.DeleteById): expected delete by id to delete 1 row but 0 were deleted")
}

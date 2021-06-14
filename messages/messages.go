package messages

import (
	"github.com/mdev5000/qlik_message/apperrors"
	"github.com/mdev5000/qlik_message/data"
	"github.com/mdev5000/qlik_message/logging"
	"github.com/pkg/errors"
)

const InvalidMessageId = 0

type CreateMessage struct {
	Message string
}

type MessageId = data.MessageId
type MessageVersion = data.MessageVersion

type Message = data.Message

type Service struct {
	log  *logging.Logger
	repo *data.MessagesRepository
}

func NewService(log *logging.Logger, repo *data.MessagesRepository) *Service {
	return &Service{
		log:  log,
		repo: repo,
	}
}

func (ms *Service) Create(message CreateMessage) (MessageId, error) {
	const op = "MessagesService.Create"

	if message.Message == "" {
		re := apperrors.Error{Op: op, EType: apperrors.ETInvalid}
		re.AddResponse(apperrors.FieldErrorResponse{
			Field: "message",
			Error: "Message field cannot be blank.",
		})
		return InvalidMessageId, &re
	}

	now := data.NowUTC()

	id, err := ms.repo.Create(data.CreateMessage{
		Message:   message.Message,
		CreatedAt: now,
	})
	if err != nil {
		return id, err
	}
	return id, err
}

func (ms *Service) Read(id MessageId) (*Message, error) {
	const op = "MessagesService.Read"
	var message Message
	err := ms.repo.GetById(id, &message)
	if errors.Is(err, data.IdMissingError{}) {
		return nil, &apperrors.Error{Op: op, EType: apperrors.ETNotFound, Err: err}
	}
	return &message, err
}

func (ms *Service) Delete(id MessageId) error {
	return ms.repo.DeleteById(id)
}

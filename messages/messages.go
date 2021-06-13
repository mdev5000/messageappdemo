package messages

import (
	"github.com/mdev5000/qlik_message/apperrors"
	"github.com/mdev5000/qlik_message/data"
	"github.com/mdev5000/qlik_message/logging"
)

const InvalidMessageId = 0

type CreateMessage struct {
	Message string
}

type MessageId = data.MessageId

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
		re := apperrors.Error{Op: op}
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

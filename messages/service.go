package messages

import (
	"errors"

	"github.com/mdev5000/messageappdemo/apperrors"
	"github.com/mdev5000/messageappdemo/logging"
)

// Service is the primary interface between the domain and the outside layers of the application. All interaction with
// this package should be done via this struct for non-domain packages (ex. server).
type Service struct {
	log  *logging.Logger
	repo Repository
}

func NewService(log *logging.Logger, repo Repository) *Service {
	return &Service{
		log:  log,
		repo: repo,
	}
}

// Create creates a new message. The message body cannot be empty and have a character limit of MaxMessageCharLength.
func (ms *Service) Create(message ModifyMessage) (MessageId, error) {
	const op = "MessagesService.Create"

	if err := validateMessage(op, message); err != nil {
		return noOp, err
	}

	now := nowUTC()

	id, err := ms.repo.Create(CreateMessage{
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
	if errors.Is(err, IdMissingError{}) {
		return nil, &apperrors.Error{Op: op, EType: apperrors.ETNotFound, Err: err}
	}
	return &message, err
}

func (ms *Service) Delete(id MessageId) error {
	return ms.repo.DeleteById(id)
}

// Update updates a message. The message body cannot be empty and have a character limit of MaxMessageCharLength.
func (ms *Service) Update(id MessageId, message ModifyMessage) (MessageVersion, error) {
	const op = "MessagesService.Update"

	if err := validateMessage(op, message); err != nil {
		return noOp, err
	}

	version, err := ms.repo.UpdateById(id, ModifyMessage{
		Message: message.Message,
	})
	if err != nil {
		return version, err
	}
	return version, err
}

func (ms *Service) List(query MessageQuery) ([]*Message, error) {
	var messagesRaw []*Message
	if err := ms.repo.GetAllQuery(query, &messagesRaw); err != nil {
		return nil, err
	}

	out := make([]*Message, len(messagesRaw))
	for i, mRaw := range messagesRaw {
		m := Message{
			Id:        mRaw.Id,
			Version:   mRaw.Version,
			CreatedAt: mRaw.CreatedAt,
			UpdatedAt: mRaw.UpdatedAt,
			Message:   mRaw.Message,
		}
		out[i] = &m
	}

	return out, nil
}

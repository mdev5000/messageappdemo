package messages

import (
	"errors"
	"github.com/mdev5000/qlik_message/apperrors"
	"github.com/mdev5000/qlik_message/logging"
)

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

func (ms *Service) Create(message ModifyMessage) (MessageId, error) {
	const op = "MessagesService.Create"

	if err := validateMessage(message, op); err != nil {
		return NoOp, err
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

func (ms *Service) Update(id MessageId, message ModifyMessage) (MessageVersion, error) {
	const op = "MessagesService.Update"

	if err := validateMessage(message, op); err != nil {
		return NoOp, err
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

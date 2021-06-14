package messages

import (
	"fmt"
	"github.com/mdev5000/qlik_message/apperrors"
	"github.com/mdev5000/qlik_message/data"
	"github.com/mdev5000/qlik_message/logging"
	"github.com/pkg/errors"
	"strings"
	"time"
)

const NoOp = 0

type ModifyMessage struct {
	Message string
}

type MessageId = data.MessageId
type MessageVersion = data.MessageVersion

type Message struct {
	Id           MessageId
	Version      MessageVersion
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Message      string
	IsPalindrome bool
}

type MessageQuery = data.MessageQuery

type isDbField = bool

var queryableFields = map[string]isDbField{
	"id":           true,
	"version":      true,
	"createdAt":    true,
	"updatedAt":    true,
	"message":      true,
	"isPalindrome": false,
}

var dynamicFields = map[string]func(message *Message) error{
	"isPalindrome": func(message *Message) error {
		message.IsPalindrome = isPalindrome(message.Message)
		return nil
	},
}

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

func (ms *Service) Create(message ModifyMessage) (MessageId, error) {
	const op = "MessagesService.Create"

	if err := validateMessage(message, op); err != nil {
		return NoOp, err
	}

	now := data.NowUTC()

	id, err := ms.repo.Create(data.ModifyMessage{
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

	var message data.Message
	err := ms.repo.GetById(id, &message)
	if errors.Is(err, data.IdMissingError{}) {
		return nil, &apperrors.Error{Op: op, EType: apperrors.ETNotFound, Err: err}
	}
	return &Message{
		Id:           message.Id,
		Version:      message.Version,
		CreatedAt:    message.CreatedAt,
		UpdatedAt:    message.UpdatedAt,
		Message:      message.Message,
		IsPalindrome: isPalindrome(message.Message),
	}, err
}

func (ms *Service) Delete(id MessageId) error {
	return ms.repo.DeleteById(id)
}

func (ms *Service) Update(id MessageId, message ModifyMessage) (MessageVersion, error) {
	const op = "MessagesService.Update"

	if err := validateMessage(message, op); err != nil {
		return NoOp, err
	}

	version, err := ms.repo.UpdateById(id, data.ModifyMessage{
		Message: message.Message,
	})
	if err != nil {
		return version, err
	}
	return version, err
}

func (ms *Service) List(query MessageQuery) ([]*Message, error) {
	// @todo currently filtering fields is flawed find a better way to do this.
	const op = "MessagesService.List"

	var dynamFields []string
	var invalidFields []string
	dbFields := map[string]struct{}{}
	for field := range query.Fields {
		isDbField, found := queryableFields[field]
		if !found {
			invalidFields = append(invalidFields, field)
			continue
		}
		if isDbField {
			dbFields[field] = struct{}{}
		} else {
			dynamFields = append(dynamFields, field)
		}
	}

	if len(invalidFields) != 0 {
		fields := strings.Join(invalidFields, ", ")
		err := fmt.Errorf("invalid fields %s", fields)
		aErr := apperrors.Error{
			EType: apperrors.ETInvalid,
			Op:    op,
			Err:   err,
			Stack: errors.WithStack(err),
		}
		aErr.AddResponse(apperrors.ErrorResponse(err.Error()))
		return nil, &aErr
	}

	query.Fields = dbFields
	var messagesRaw []*data.Message
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
		for _, field := range dynamFields {
			fn, found := dynamicFields[field]
			if !found {
				err := fmt.Errorf("invalid dynamic field '%s' specified", field)
				return nil, &apperrors.Error{
					EType: apperrors.ETInternal,
					Op:    op,
					Err:   err,
					Stack: errors.WithStack(err),
				}
			}
			if err := fn(&m); err != nil {
				return nil, &apperrors.Error{
					EType: apperrors.ETInternal,
					Op:    op,
					Err:   fmt.Errorf("dynamic field error for field '%s' (id=%d); %w", field, m.Id, err),
					Stack: errors.WithStack(err),
				}
			}
			out[i] = &m
		}
	}

	return out, nil
}

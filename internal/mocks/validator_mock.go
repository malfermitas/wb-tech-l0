package mocks

import (
	"wb-tech-l0/internal/validator"

	"github.com/stretchr/testify/mock"
)

// ValidatorMock реализует интерфейс validator.Validator.
type ValidatorMock struct {
	mock.Mock
}

var _ validator.Validator = (*ValidatorMock)(nil)

func (m *ValidatorMock) Validate(v interface{}) error {
	args := m.Called(v)
	return args.Error(0)
}

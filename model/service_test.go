package model_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/rknizzle/rkmesh/domain"
	"github.com/rknizzle/rkmesh/domain/mocks"
	"github.com/rknizzle/rkmesh/model"
)

func TestServiceGetAll(t *testing.T) {
	mockModelRepo := new(mocks.ModelRepository)
	mockFilestore := new(mocks.Filestore)
	mockModel := domain.Model{Name: "test.stl"}

	mockListModel := make([]domain.Model, 0)
	mockListModel = append(mockListModel, mockModel)

	t.Run("success", func(t *testing.T) {
		mockModelRepo.On("GetAll", mock.Anything).Return(mockListModel, nil).Once()

		u := model.NewModelService(mockModelRepo, mockFilestore, time.Second*2)

		list, err := u.GetAll(context.TODO())
		assert.NoError(t, err)
		assert.Len(t, list, len(mockListModel))

		mockModelRepo.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockModelRepo.On("GetAll", mock.Anything).Return(nil, errors.New("Unexpexted Error")).Once()

		s := model.NewModelService(mockModelRepo, mockFilestore, time.Second*2)
		list, err := s.GetAll(context.TODO())

		assert.Error(t, err)
		assert.Len(t, list, 0)
		mockModelRepo.AssertExpectations(t)
	})

}

func TestServiceGetByID(t *testing.T) {
	mockModelRepo := new(mocks.ModelRepository)
	mockFilestore := new(mocks.Filestore)
	mockModel := domain.Model{
		ID:   1,
		Name: "test.stl",
	}

	t.Run("success", func(t *testing.T) {
		mockModelRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockModel, nil).Once()
		s := model.NewModelService(mockModelRepo, mockFilestore, time.Second*2)

		m, err := s.GetByID(context.TODO(), mockModel.ID)

		assert.NoError(t, err)
		assert.NotNil(t, m)

		mockModelRepo.AssertExpectations(t)
	})
	t.Run("error-failed", func(t *testing.T) {
		mockModelRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(domain.Model{}, errors.New("Unexpected")).Once()

		s := model.NewModelService(mockModelRepo, mockFilestore, time.Second*2)

		m, err := s.GetByID(context.TODO(), mockModel.ID)

		assert.Error(t, err)
		assert.Equal(t, domain.Model{}, m)

		mockModelRepo.AssertExpectations(t)
	})

}

func TestServiceStore(t *testing.T) {
	mockModelRepo := new(mocks.ModelRepository)
	mockFilestore := new(mocks.Filestore)
	mockModel := domain.Model{
		ID:   1,
		Name: "test.stl",
	}

	t.Run("success", func(t *testing.T) {
		tempMockModel := mockModel
		tempMockModel.ID = 0
		mockModelRepo.On("Store", mock.Anything, mock.AnythingOfType("*domain.Model")).Return(nil).Once()
		mockFilestore.On("Upload", mock.Anything, mock.Anything, "test.stl").Return("", nil).Once()

		s := model.NewModelService(mockModelRepo, mockFilestore, time.Second*2)

		err := s.Store(context.TODO(), &tempMockModel, strings.NewReader("test"), "test.stl", 1)

		assert.NoError(t, err)
		assert.Equal(t, mockModel.Name, tempMockModel.Name)
		mockModelRepo.AssertExpectations(t)
	})
}

func TestServiceDelete(t *testing.T) {
	mockModelRepo := new(mocks.ModelRepository)
	mockFilestore := new(mocks.Filestore)
	mockModel := domain.Model{Name: "test.stl"}

	t.Run("success", func(t *testing.T) {
		mockModelRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockModel, nil).Once()

		mockModelRepo.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(nil).Once()

		s := model.NewModelService(mockModelRepo, mockFilestore, time.Second*2)

		err := s.Delete(context.TODO(), mockModel.ID)

		assert.NoError(t, err)
		mockModelRepo.AssertExpectations(t)
	})
	t.Run("model-does-not-exist", func(t *testing.T) {
		mockModelRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(domain.Model{}, nil).Once()

		s := model.NewModelService(mockModelRepo, mockFilestore, time.Second*2)

		err := s.Delete(context.TODO(), mockModel.ID)

		assert.Error(t, err)
		mockModelRepo.AssertExpectations(t)
	})
	t.Run("error-happens-in-db", func(t *testing.T) {
		mockModelRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(domain.Model{}, errors.New("Unexpected Error")).Once()

		s := model.NewModelService(mockModelRepo, mockFilestore, time.Second*2)

		err := s.Delete(context.TODO(), mockModel.ID)

		assert.Error(t, err)
		mockModelRepo.AssertExpectations(t)
	})
}

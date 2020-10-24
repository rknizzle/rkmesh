package service

import (
	"context"
	"time"

	//"github.com/sirupsen/logrus"
	//"golang.org/x/sync/errgroup"

	"github.com/rknizzle/rkmesh/domain"
)

type modelService struct {
	modelRepo      domain.ModelRepository
	contextTimeout time.Duration
}

func NewModelService(m domain.ModelRepository, timeout time.Duration) domain.ModelService {
	return &modelService{
		modelRepo:      m,
		contextTimeout: timeout,
	}
}

func (m *modelService) GetAll(c context.Context) (res []domain.Model, err error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	res, err = m.modelRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return
}

func (m *modelService) GetByID(c context.Context, id int64) (res domain.Model, err error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	res, err = m.modelRepo.GetByID(ctx, id)
	if err != nil {
		return
	}

	return
}

func (m *modelService) GetByName(c context.Context, name string) (res domain.Model, err error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()
	res, err = m.modelRepo.GetByName(ctx, name)
	if err != nil {
		return
	}

	return
}

func (m *modelService) Store(c context.Context, model *domain.Model) (err error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()
	/*
		// What is this doing / checking for?
		existedModel, _ := m.GetByName(ctx, model.Name)
		// if it doesnt get an empty result then you know it hasnt been created yet
		if existedModel != (domain.Model{}) {
			return domain.ErrConflict
		}
	*/

	err = m.modelRepo.Store(ctx, model)
	return
}

func (m *modelService) Delete(c context.Context, id int64) (err error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()
	existedModel, err := m.modelRepo.GetByID(ctx, id)
	if err != nil {
		return
	}
	if existedModel == (domain.Model{}) {
		return domain.ErrNotFound
	}
	return m.modelRepo.Delete(ctx, id)
}

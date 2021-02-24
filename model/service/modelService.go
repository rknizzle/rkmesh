package service

import (
	"context"
	"time"

	//"github.com/sirupsen/logrus"
	//"golang.org/x/sync/errgroup"

	"github.com/rknizzle/rkmesh/domain"
	"io"
)

type modelService struct {
	modelRepo      domain.ModelRepository
	s3Repo         domain.FileRepository
	contextTimeout time.Duration
}

func NewModelService(m domain.ModelRepository, s domain.FileRepository, timeout time.Duration) domain.ModelService {
	return &modelService{
		modelRepo:      m,
		s3Repo:         s,
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

func (m *modelService) Store(c context.Context, model *domain.Model, file io.Reader, filename string) (err error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	_, err = m.s3Repo.Upload(ctx, file, filename)
	if err != nil {
		return err
	}

	model.Name = filename

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

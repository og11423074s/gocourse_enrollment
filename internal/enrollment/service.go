package enrollment

import (
	"context"
	"github.com/og11423074s/gocourse_enrollment/internal/domain"
	"log"
)

type (
	Filters struct {
		UserID   string
		CourseID string
	}

	Service interface {
		Create(ctx context.Context, userID, courseID string) (*domain.Enrollment, error)
		GetAll(ctx context.Context, filters Filters, offSet, limit int) ([]domain.Enrollment, error)
		Update(ctx context.Context, id string, status *string) error
		Count(ctx context.Context, filters Filters) (int, error)
	}

	service struct {
		log  *log.Logger
		repo Repository
	}
)

func NewService(log *log.Logger, repo Repository) Service {
	return &service{
		log:  log,
		repo: repo,
	}
}

func (s *service) Create(ctx context.Context, userID, courseID string) (*domain.Enrollment, error) {

	enroll := &domain.Enrollment{
		UserID:   userID,
		CourseID: courseID,
		Status:   "P",
	}

	if err := s.repo.Create(ctx, enroll); err != nil {
		return nil, err
	}

	return enroll, nil
}

func (s *service) GetAll(ctx context.Context, filters Filters, offSet, limit int) ([]domain.Enrollment, error) {
	enrollments, err := s.repo.GetAll(ctx, filters, offSet, limit)
	if err != nil {
		return nil, err
	}

	return enrollments, nil
}
func (s *service) Update(ctx context.Context, id string, status *string) error {
	if err := s.repo.Update(ctx, id, status); err != nil {
		return err
	}

	return nil
}

func (s *service) Count(ctx context.Context, filters Filters) (int, error) {
	return s.repo.Count(ctx, filters)
}

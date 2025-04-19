package enrollment

import (
	"context"
	"github.com/og11423074s/gocourse_domain/domain"
	"log"

	courseSdk "github.com/og11423074s/go_course_sdk/course"
	userSdk "github.com/og11423074s/go_course_sdk/user"
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
		log       *log.Logger
		userTrans userSdk.Transport
		courseSdk courseSdk.Transport
		repo      Repository
	}
)

func NewService(log *log.Logger, repo Repository, courseTrans courseSdk.Transport, userTrans userSdk.Transport) Service {
	return &service{
		log:       log,
		userTrans: userTrans,
		courseSdk: courseTrans,
		repo:      repo,
	}
}

func (s *service) Create(ctx context.Context, userID, courseID string) (*domain.Enrollment, error) {

	enroll := &domain.Enrollment{
		UserID:   userID,
		CourseID: courseID,
		Status:   domain.Pending,
	}
	if _, err := s.userTrans.Get(userID); err != nil {
		return nil, err
	}

	if _, err := s.courseSdk.Get(courseID); err != nil {
		return nil, err
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

	if status != nil {
		switch domain.EnrollmentStatus(*status) {
		case domain.Pending, domain.Active, domain.Studying, domain.Inactive:
		default:
			return ErrorInvalidStatus{*status}
		}
	}

	if err := s.repo.Update(ctx, id, status); err != nil {
		return err
	}

	return nil
}

func (s *service) Count(ctx context.Context, filters Filters) (int, error) {
	return s.repo.Count(ctx, filters)
}

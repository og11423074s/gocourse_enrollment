package enrollment

import (
	"context"
	"github.com/og11423074s/gocourse_domain/domain"
	"gorm.io/gorm"
	"log"
)

type (
	Repository interface {
		Create(ctx context.Context, enroll *domain.Enrollment) error
		GetAll(ctx context.Context, filters Filters, offSet, limit int) ([]domain.Enrollment, error)
		Update(ctx context.Context, id string, status *string) error
		Count(ctx context.Context, filters Filters) (int, error)
	}

	repo struct {
		db  *gorm.DB
		log *log.Logger
	}
)

func NewRepo(log *log.Logger, db *gorm.DB) Repository {
	return &repo{
		db:  db,
		log: log,
	}
}

func (r *repo) Create(ctx context.Context, enroll *domain.Enrollment) error {
	if err := r.db.WithContext(ctx).Create(enroll).Error; err != nil {
		r.log.Printf("error: %v", err)
		return err
	}
	r.log.Println("enrollment created with id: ", enroll.ID)
	return nil
}

func (r *repo) GetAll(ctx context.Context, filters Filters, offSet, limit int) ([]domain.Enrollment, error) {
	var enrollments []domain.Enrollment

	tx := r.db.WithContext(ctx).Model(&enrollments)
	tx = applyFilters(tx, filters)
	tx = tx.Limit(limit).Offset(offSet)

	result := tx.Order("created_at desc").Find(&enrollments)

	if result.Error != nil {
		r.log.Println(result.Error)
		return nil, result.Error
	}

	return enrollments, nil
}

func (r *repo) Update(ctx context.Context, id string, status *string) error {
	values := make(map[string]interface{})

	if status != nil {
		values["status"] = *status
	}

	result := r.db.WithContext(ctx).Model(&domain.Enrollment{}).Where("id = ?", id).Updates(values)

	if result.Error != nil {
		r.log.Println(result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		r.log.Println("enrollment %s does not exist", id)
		return ErrNotFound{id}
	}

	return nil
}

func (r *repo) Count(ctx context.Context, filters Filters) (int, error) {
	var count int64
	tx := r.db.WithContext(ctx).Model(&domain.Enrollment{})
	tx = applyFilters(tx, filters)

	if err := tx.Count(&count).Error; err != nil {
		r.log.Printf("error: %v", err)
		return 0, err
	}

	return int(count), nil
}

func applyFilters(tx *gorm.DB, filters Filters) *gorm.DB {
	if filters.UserID != "" {
		tx = tx.Where("user_id = ?", filters.UserID)
	}
	if filters.CourseID != "" {
		tx = tx.Where("course_id = ?", filters.CourseID)
	}
	return tx
}

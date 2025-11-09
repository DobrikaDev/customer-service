package sql

import (
	"DobrikaDev/customer-service/internal/domain"
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	"go.uber.org/zap"
)

const (
	pgErrUniqueViolation     = "23505"
	pgErrForeignKeyViolation = "23503"
)

func (s *SqlStorage) GetFeedbackByID(ctx context.Context, id string) (*domain.Feedback, error) {
	query, args := sq.Select(
		"f.id",
		"f.user_id",
		"f.task_id",
		"f.customer_id",
		"f.rating",
		"f.comment",
		"f.created_at",
		"f.updated_at",
	).
		From("feedbacks f").
		Where(sq.Eq{"f.id": id}).
		PlaceholderFormat(sq.Dollar).
		MustSql()
	var feedback domain.Feedback
	err := s.trf.Transaction(ctx).GetContext(ctx, &feedback, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrFeedbackNotFound
		}
		s.logger.Error("failed to get feedback by id", zap.Error(err), zap.String("id", id))
		return nil, ErrFeedbackInternal
	}
	return &feedback, nil
}

func (s *SqlStorage) CreateFeedback(ctx context.Context, feedback *domain.Feedback) (*domain.Feedback, error) {
	feedback.ID = uuid.NewString()
	query, args := sq.Insert("feedbacks").
		Columns("id", "user_id", "task_id", "rating", "comment").
		Values(feedback.ID, feedback.UserID, feedback.TaskID, feedback.Rating, feedback.Comment).
		PlaceholderFormat(sq.Dollar).
		MustSql()
	_, err := s.trf.Transaction(ctx).ExecContext(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgErrUniqueViolation:
				return nil, ErrFeedbackAlreadyExists
			case pgErrForeignKeyViolation:
				return nil, ErrFeedbackInvalid
			}
		}
		s.logger.Error("failed to create feedback", zap.Error(err), zap.String("user_id", feedback.UserID), zap.String("task_id", feedback.TaskID))
		return nil, ErrFeedbackInternal
	}
	return feedback, nil
}

type GetFeedbacksOptions func(sb sq.SelectBuilder) sq.SelectBuilder

func WithTaskID(taskID string) GetFeedbacksOptions {
	return func(sb sq.SelectBuilder) sq.SelectBuilder {
		if taskID != "" {
			sb = sb.Where(sq.Eq{"f.task_id": taskID})
		}
		return sb
	}
}

func WithUserID(userID string) GetFeedbacksOptions {
	return func(sb sq.SelectBuilder) sq.SelectBuilder {
		if userID != "" {
			sb = sb.Where(sq.Eq{"f.user_id": userID})
		}
		return sb
	}
}

func WithCustomerID(customerID string) GetFeedbacksOptions {
	return func(sb sq.SelectBuilder) sq.SelectBuilder {
		if customerID != "" {
			sb = sb.Where(sq.Eq{"f.customer_id": customerID})
		}
		return sb
	}
}

func WithLimit(limit int) GetFeedbacksOptions {
	return func(sb sq.SelectBuilder) sq.SelectBuilder {
		if limit > 0 {
			sb = sb.Limit(uint64(limit))
		}
		return sb
	}
}

func WithOffset(offset int) GetFeedbacksOptions {
	return func(sb sq.SelectBuilder) sq.SelectBuilder {
		return sb.Offset(uint64(offset))
	}
}

func (s *SqlStorage) GetFeedbacks(ctx context.Context, opts ...GetFeedbacksOptions) ([]*domain.Feedback, int, error) {
	sb := sq.Select(
		"f.id",
		"f.user_id",
		"f.task_id",
		"f.rating",
		"f.comment",
	).
		From("feedbacks f").
		PlaceholderFormat(sq.Dollar).
		OrderBy("f.created_at DESC")

	if len(opts) > 0 {
		for _, opt := range opts {
			sb = opt(sb)
		}
	}

	query, args := sb.MustSql()
	feedbacks := make([]*domain.Feedback, 0, 10)
	err := s.trf.Transaction(ctx).SelectContext(ctx, &feedbacks, query, args...)
	if err != nil {
		return nil, 0, ErrFeedbackInternal
	}

	count, err := s.CountFeedbacks(ctx, opts...)
	if err != nil {
		return nil, 0, ErrFeedbackInternal
	}

	return feedbacks, count, nil
}

func (s *SqlStorage) CountFeedbacks(ctx context.Context, opts ...GetFeedbacksOptions) (int, error) {
	sb := sq.Select("COUNT(*)").From("feedbacks").PlaceholderFormat(sq.Dollar)
	if len(opts) > 0 {
		for _, opt := range opts {
			sb = opt(sb)
		}
	}
	query, args := sb.MustSql()
	var count int
	err := s.trf.Transaction(ctx).GetContext(ctx, &count, query, args...)
	if err != nil {
		return 0, ErrFeedbackInternal
	}
	return count, nil
}

package customer

import (
	"DobrikaDev/customer-service/internal/domain"
	"DobrikaDev/customer-service/internal/storage/sql"
	"context"
	"errors"

	"go.uber.org/zap"
)

func (s *CustomerService) GetCustomerByMaxID(ctx context.Context, maxID string) (*domain.Customer, error) {
	customer, err := s.storage.GetCustomerByMaxID(ctx, maxID)
	if err != nil {
		if errors.Is(err, sql.ErrCustomerNotFound) {
			return nil, ErrCustomerNotFound
		}
		s.logger.Error("failed to get customer by max id", zap.Error(err), zap.String("max_id", maxID))
		return nil, ErrCustomerInternal
	}
	return customer, nil
}

func (s *CustomerService) GetCustomers(ctx context.Context, maxID string, limit int, offset int) ([]*domain.Customer, int, error) {
	opts := []sql.GetCustomersOption{
		sql.WithCustomerMaxID(maxID),
		sql.WithCustomerLimit(limit),
		sql.WithCustomerOffset(offset),
	}
	customers, count, err := s.storage.GetCustomers(ctx, opts...)
	if err != nil {
		return nil, 0, ErrCustomerInternal
	}
	return customers, count, nil
}

func (s *CustomerService) CountCustomers(ctx context.Context, opts ...sql.GetCustomersOption) (int, error) {
	count, err := s.storage.CountCustomers(ctx, opts...)
	if err != nil {
		return 0, ErrCustomerInternal
	}
	return count, nil
}

func (s *CustomerService) CreateCustomer(ctx context.Context, customer *domain.Customer) (*domain.Customer, error) {
	customer, err := s.storage.CreateCustomer(ctx, customer)
	if err != nil {
		if errors.Is(err, sql.ErrCustomerAlreadyExists) {
			return nil, ErrCustomerAlreadyExists
		}
		s.logger.Error("failed to create customer", zap.Error(err), zap.Any("customer", customer))
		return nil, ErrCustomerInternal
	}
	return customer, nil
}

func (s *CustomerService) UpdateCustomer(ctx context.Context, customer *domain.Customer) (*domain.Customer, error) {
	customer, err := s.storage.UpdateCustomer(ctx, customer)
	if err != nil {
		if errors.Is(err, sql.ErrCustomerNotFound) {
			return nil, ErrCustomerNotFound
		}
		s.logger.Error("failed to update customer", zap.Error(err), zap.Any("customer", customer))
		return nil, ErrCustomerInternal
	}
	return customer, nil
}

func (s *CustomerService) DeleteCustomer(ctx context.Context, maxID string) error {
	err := s.storage.DeleteCustomer(ctx, maxID)
	if err != nil {
		if errors.Is(err, sql.ErrCustomerNotFound) {
			return ErrCustomerNotFound
		}
		s.logger.Error("failed to delete customer", zap.Error(err), zap.String("max_id", maxID))
		return ErrCustomerInternal
	}
	return nil
}

func (s *CustomerService) GetFeedbacks(ctx context.Context, taskID string, userID string, limit int, offset int) ([]*domain.Feedback, int, error) {
	opts := []sql.GetFeedbacksOptions{
		sql.WithTaskID(taskID),
		sql.WithUserID(userID),
		sql.WithLimit(limit),
		sql.WithOffset(offset),
	}
	feedbacks, count, err := s.storage.GetFeedbacks(ctx, opts...)
	if err != nil {
		return nil, 0, ErrFeedbackInternal
	}
	return feedbacks, count, nil
}

func (s *CustomerService) GetFeedbackByID(ctx context.Context, id string) (*domain.Feedback, error) {
	feedback, err := s.storage.GetFeedbackByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrFeedbackNotFound) {
			return nil, ErrFeedbackNotFound
		}
		s.logger.Error("failed to get feedback by id", zap.Error(err), zap.String("id", id))
		return nil, ErrFeedbackInternal
	}
	return feedback, nil
}

func (s *CustomerService) CreateFeedback(ctx context.Context, feedback *domain.Feedback) (*domain.Feedback, error) {
	if feedback.Rating < 1 || feedback.Rating > 5 {
		return nil, ErrFeedbackInvalid
	}
	if feedback.CustomerID == feedback.UserID {
		return nil, ErrFeedbackInvalid
	}
	if feedback.TaskID == "" {
		return nil, ErrFeedbackInvalid
	}
	feedback, err := s.storage.CreateFeedback(ctx, feedback)
	if err != nil {
		if errors.Is(err, sql.ErrFeedbackAlreadyExists) {
			return nil, ErrFeedbackAlreadyExists
		}
		s.logger.Error("failed to create feedback", zap.Error(err), zap.Any("feedback", feedback))
		return nil, ErrFeedbackInternal
	}
	return feedback, nil
}

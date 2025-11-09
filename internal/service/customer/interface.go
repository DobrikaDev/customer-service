package customer

import (
	"DobrikaDev/customer-service/internal/domain"
	"DobrikaDev/customer-service/internal/storage/sql"
	"DobrikaDev/customer-service/utils/config"
	"context"

	"go.uber.org/zap"
)

type storage interface {
	GetCustomerByMaxID(ctx context.Context, maxID string) (*domain.Customer, error)
	GetCustomers(ctx context.Context, opts ...sql.GetCustomersOption) ([]*domain.Customer, int, error)
	CountCustomers(ctx context.Context, opts ...sql.GetCustomersOption) (int, error)
	CreateCustomer(ctx context.Context, customer *domain.Customer) (*domain.Customer, error)
	UpdateCustomer(ctx context.Context, customer *domain.Customer) (*domain.Customer, error)
	DeleteCustomer(ctx context.Context, maxID string) error

	GetFeedbacks(ctx context.Context, opts ...sql.GetFeedbacksOptions) ([]*domain.Feedback, int, error)
	GetFeedbackByID(ctx context.Context, id string) (*domain.Feedback, error)
	CreateFeedback(ctx context.Context, feedback *domain.Feedback) (*domain.Feedback, error)
}

type CustomerService struct {
	storage storage
	cfg     *config.Config
	logger  *zap.Logger
}

func NewCustomerService(storage storage, cfg *config.Config, logger *zap.Logger) *CustomerService {
	return &CustomerService{storage: storage, cfg: cfg, logger: logger}
}

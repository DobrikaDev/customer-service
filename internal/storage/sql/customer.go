package sql

import (
	"DobrikaDev/customer-service/internal/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

const customerTableName = "customers"

var customerSelectColumns = []string{
	"c.max_id",
	"c.name",
	"c.about",
	"c.type",
	"c.created_at",
	"c.updated_at",
}

type (
	customerOption interface {
		applySelect(sq.SelectBuilder) sq.SelectBuilder
		applyCount(sq.SelectBuilder) sq.SelectBuilder
	}

	customerOptionFunc struct {
		selectFn func(sq.SelectBuilder) sq.SelectBuilder
		countFn  func(sq.SelectBuilder) sq.SelectBuilder
	}
)

func (f customerOptionFunc) applySelect(sb sq.SelectBuilder) sq.SelectBuilder {
	if f.selectFn != nil {
		return f.selectFn(sb)
	}
	return sb
}

func (f customerOptionFunc) applyCount(sb sq.SelectBuilder) sq.SelectBuilder {
	if f.countFn != nil {
		return f.countFn(sb)
	}
	return sb
}

type GetCustomersOption interface {
	customerOption
}

func WithCustomerMaxID(maxID string) GetCustomersOption {
	return customerOptionFunc{
		selectFn: func(sb sq.SelectBuilder) sq.SelectBuilder {
			if maxID != "" {
				sb = sb.Where(sq.Eq{"c.max_id": maxID})
			}
			return sb
		},
		countFn: func(sb sq.SelectBuilder) sq.SelectBuilder {
			if maxID != "" {
				sb = sb.Where(sq.Eq{"c.max_id": maxID})
			}
			return sb
		},
	}
}

func WithCustomerName(name string) GetCustomersOption {
	return customerOptionFunc{
		selectFn: func(sb sq.SelectBuilder) sq.SelectBuilder {
			if name != "" {
				sb = sb.Where(sq.Eq{"c.name": name})
			}
			return sb
		},
		countFn: func(sb sq.SelectBuilder) sq.SelectBuilder {
			if name != "" {
				sb = sb.Where(sq.Eq{"c.name": name})
			}
			return sb
		},
	}
}

func WithCustomerNameLike(pattern string) GetCustomersOption {
	return customerOptionFunc{
		selectFn: func(sb sq.SelectBuilder) sq.SelectBuilder {
			if pattern != "" {
				sb = sb.Where(sq.Expr("c.name ILIKE ?", fmt.Sprintf("%%%s%%", pattern)))
			}
			return sb
		},
		countFn: func(sb sq.SelectBuilder) sq.SelectBuilder {
			if pattern != "" {
				sb = sb.Where(sq.Expr("c.name ILIKE ?", fmt.Sprintf("%%%s%%", pattern)))
			}
			return sb
		},
	}
}

func WithCustomerType(customerType domain.CustomerType) GetCustomersOption {
	return customerOptionFunc{
		selectFn: func(sb sq.SelectBuilder) sq.SelectBuilder {
			if customerType != "" {
				sb = sb.Where(sq.Eq{"c.type": customerType})
			}
			return sb
		},
		countFn: func(sb sq.SelectBuilder) sq.SelectBuilder {
			if customerType != "" {
				sb = sb.Where(sq.Eq{"c.type": customerType})
			}
			return sb
		},
	}
}

func WithCustomerLimit(limit int) GetCustomersOption {
	return customerOptionFunc{
		selectFn: func(sb sq.SelectBuilder) sq.SelectBuilder {
			if limit > 0 {
				sb = sb.Limit(uint64(limit))
			}
			return sb
		},
	}
}

func WithCustomerOffset(offset int) GetCustomersOption {
	return customerOptionFunc{
		selectFn: func(sb sq.SelectBuilder) sq.SelectBuilder {
			if offset > 0 {
				sb = sb.Offset(uint64(offset))
			}
			return sb
		},
	}
}

func (s *SqlStorage) GetCustomerByMaxID(ctx context.Context, maxID string) (*domain.Customer, error) {
	query, args := sq.Select(customerSelectColumns...).
		From(fmt.Sprintf("%s c", customerTableName)).
		Where(sq.Eq{"c.max_id": maxID}).
		PlaceholderFormat(sq.Dollar).
		MustSql()

	var customer domain.Customer
	err := s.trf.Transaction(ctx).GetContext(ctx, &customer, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCustomerNotFound
		}
		s.logger.Error("failed to get customer by max_id", zap.Error(err), zap.String("max_id", maxID))
		return nil, ErrCustomerInternal
	}

	return &customer, nil
}

func (s *SqlStorage) GetCustomers(ctx context.Context, opts ...GetCustomersOption) ([]*domain.Customer, int, error) {
	sb := sq.Select(customerSelectColumns...).
		From(fmt.Sprintf("%s c", customerTableName)).
		OrderBy("c.created_at DESC").
		PlaceholderFormat(sq.Dollar)

	if len(opts) > 0 {
		for _, opt := range opts {
			if opt == nil {
				continue
			}
			sb = opt.applySelect(sb)
		}
	}

	query, args := sb.MustSql()

	customers := make([]*domain.Customer, 0)
	err := s.trf.Transaction(ctx).SelectContext(ctx, &customers, query, args...)
	if err != nil {
		s.logger.Error("failed to get customers", zap.Error(err))
		return nil, 0, ErrCustomerInternal
	}

	count, err := s.CountCustomers(ctx, opts...)
	if err != nil {
		s.logger.Error("failed to count customers", zap.Error(err))
		return nil, 0, ErrCustomerInternal
	}

	return customers, count, nil
}

func (s *SqlStorage) CountCustomers(ctx context.Context, opts ...GetCustomersOption) (int, error) {
	sb := sq.Select("COUNT(*)").
		From(fmt.Sprintf("%s c", customerTableName)).
		PlaceholderFormat(sq.Dollar)

	if len(opts) > 0 {
		for _, opt := range opts {
			if opt == nil {
				continue
			}
			sb = opt.applyCount(sb)
		}
	}

	query, args := sb.MustSql()

	var count int
	err := s.trf.Transaction(ctx).GetContext(ctx, &count, query, args...)
	if err != nil {
		return 0, ErrCustomerInternal
	}

	return count, nil
}

func (s *SqlStorage) CreateCustomer(ctx context.Context, customer *domain.Customer) (*domain.Customer, error) {
	query, args := sq.Insert(customerTableName).
		Columns("max_id", "name", "about", "type").
		Values(customer.MaxID, customer.Name, customer.About, customer.Type).
		Suffix("RETURNING max_id, name, about, type, created_at, updated_at").
		PlaceholderFormat(sq.Dollar).
		MustSql()

	var created domain.Customer
	err := s.trf.Transaction(ctx).GetContext(ctx, &created, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgErrUniqueViolation:
				return nil, ErrCustomerAlreadyExists
			case pgErrForeignKeyViolation:
				return nil, ErrCustomerInvalid
			}
		}
		s.logger.Error("failed to create customer", zap.Error(err), zap.String("max_id", customer.MaxID))
		return nil, ErrCustomerInternal
	}

	return &created, nil
}

func (s *SqlStorage) UpdateCustomer(ctx context.Context, customer *domain.Customer) (*domain.Customer, error) {
	query, args := sq.Update(customerTableName).
		Set("name", customer.Name).
		Set("about", customer.About).
		Set("type", customer.Type).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"max_id": customer.MaxID}).
		Suffix("RETURNING max_id, name, about, type, created_at, updated_at").
		PlaceholderFormat(sq.Dollar).
		MustSql()

	var updated domain.Customer
	err := s.trf.Transaction(ctx).GetContext(ctx, &updated, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCustomerNotFound
		}

		s.logger.Error("failed to update customer", zap.Error(err), zap.String("max_id", customer.MaxID))
		return nil, ErrCustomerInternal
	}

	return &updated, nil
}

func (s *SqlStorage) DeleteCustomer(ctx context.Context, maxID string) error {
	query, args := sq.Delete(customerTableName).
		Where(sq.Eq{"max_id": maxID}).
		PlaceholderFormat(sq.Dollar).
		MustSql()

	result, err := s.trf.Transaction(ctx).ExecContext(ctx, query, args...)
	if err != nil {
		s.logger.Error("failed to delete customer", zap.Error(err), zap.String("max_id", maxID))
		return ErrCustomerInternal
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		s.logger.Error("failed to check affected rows when deleting customer", zap.Error(err), zap.String("max_id", maxID))
		return ErrCustomerInternal
	}

	if rowsAffected == 0 {
		return ErrCustomerNotFound
	}

	return nil
}

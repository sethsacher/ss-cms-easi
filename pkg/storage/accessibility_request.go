package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/cmsgov/easi-app/pkg/appcontext"
	"github.com/cmsgov/easi-app/pkg/apperrors"
	"github.com/cmsgov/easi-app/pkg/graph/model"
	"github.com/cmsgov/easi-app/pkg/models"
)

// CreateAccessibilityRequest adds a new accessibility request in the database
func (s *Store) CreateAccessibilityRequest(ctx context.Context, request *model.AccessibilityRequest) (*model.AccessibilityRequest, error) {
	if request.ID == uuid.Nil {
		request.ID = uuid.New()
	}
	createAt := s.clock.Now()
	if request.CreatedAt == nil {
		request.CreatedAt = &createAt
	}
	if request.UpdatedAt == nil {
		request.UpdatedAt = &createAt
	}
	const createRequestSQL = `
		INSERT INTO accessibility_request (
			id,
			name,
			intake_id,
			created_at,
			updated_at
		)
		VALUES (
			:id,
			:name,
			:intake_id,
		    :created_at,
			:updated_at
		)`
	_, err := s.db.NamedExec(
		createRequestSQL,
		request,
	)
	if err != nil {
		appcontext.ZLogger(ctx).Error("Failed to create accessibility request", zap.Error(err))
		return nil, err
	}
	return s.FetchAccessibilityRequestByID(ctx, request.ID)
}

// FetchAccessibilityRequestByID queries the DB for an accessibility matching the given ID
func (s *Store) FetchAccessibilityRequestByID(ctx context.Context, id uuid.UUID) (*model.AccessibilityRequest, error) {
	request := model.AccessibilityRequest{}

	err := s.db.Get(&request, `SELECT * FROM accessibility_request WHERE id=$1`, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &apperrors.ResourceNotFoundError{Err: err, Resource: models.SystemIntake{}}
		}
		appcontext.ZLogger(ctx).Error("Failed to fetch accessibility request", zap.Error(err), zap.String("id", id.String()))
		return nil, &apperrors.QueryError{
			Err:       err,
			Model:     id,
			Operation: apperrors.QueryFetch,
		}
	}

	return &request, nil
}

// FetchAccessibilityRequests queries the DB for an accessibility requests.
// TODO implement cursor pagination
func (s *Store) FetchAccessibilityRequests(ctx context.Context) ([]model.AccessibilityRequest, error) {
	requests := []model.AccessibilityRequest{}

	err := s.db.Select(&requests, `SELECT * FROM accessibility_request`)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return requests, nil
		}
		appcontext.ZLogger(ctx).Error("Failed to fetch accessibility requests", zap.Error(err))
		return nil, &apperrors.QueryError{
			Err:       err,
			Operation: apperrors.QueryFetch,
		}
	}

	return requests, nil
}

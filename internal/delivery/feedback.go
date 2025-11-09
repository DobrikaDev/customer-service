package delivery

import (
	"DobrikaDev/customer-service/internal/domain"
	customerpb "DobrikaDev/customer-service/internal/generated/proto/customer"
	"context"

	"github.com/dr3dnought/gospadi"
	"go.uber.org/zap"
)

func (s *Server) CreateFeedback(ctx context.Context, req *customerpb.CreateFeedbackRequest) (*customerpb.CreateFeedbackResponse, error) {
	if req.Feedback == nil {
		return &customerpb.CreateFeedbackResponse{
			Error: &customerpb.Error{
				Code:    customerpb.ErrorCode_ERROR_CODE_VALIDATION,
				Message: "feedback is required",
			},
		}, nil
	}
	feedback, err := s.customerService.CreateFeedback(ctx, &domain.Feedback{
		CustomerID: req.Feedback.CustomerId,
		UserID:     req.Feedback.UserId,
		Rating:     int(req.Feedback.Rating),
		Comment:    req.Feedback.Comment,
		TaskID:     req.Feedback.TaskId,
	})
	if err != nil {
		return &customerpb.CreateFeedbackResponse{
			Error: convertErrorToProto(err),
		}, nil
	}
	s.logger.Info("feedback created", zap.Any("feedback", feedback))
	return &customerpb.CreateFeedbackResponse{
		Feedback: convertFeedbackToProto(feedback),
	}, nil
}

func (s *Server) GetFeedbacks(ctx context.Context, req *customerpb.GetFeedbacksRequest) (*customerpb.GetFeedbacksResponse, error) {
	if req.TaskId == "" {
		return &customerpb.GetFeedbacksResponse{
			Error: &customerpb.Error{
				Code:    customerpb.ErrorCode_ERROR_CODE_VALIDATION,
				Message: "task id is required",
			},
		}, nil
	}
	feedbacks, count, err := s.customerService.GetFeedbacks(ctx, req.TaskId, req.UserId, int(req.Limit), int(req.Offset))
	if err != nil {
		return &customerpb.GetFeedbacksResponse{
			Error: convertErrorToProto(err),
		}, nil
	}
	s.logger.Info("feedbacks fetched", zap.Any("feedbacks", feedbacks), zap.Int("count", count))
	return &customerpb.GetFeedbacksResponse{
		Feedbacks: gospadi.Map(feedbacks, convertFeedbackToProto),
		Total:     int32(count),
	}, nil
}

func (s *Server) GetFeedbackByID(ctx context.Context, req *customerpb.GetFeedbackByIDRequest) (*customerpb.GetFeedbackByIDResponse, error) {
	if req.Id == "" {
		return &customerpb.GetFeedbackByIDResponse{
			Error: &customerpb.Error{
				Code:    customerpb.ErrorCode_ERROR_CODE_VALIDATION,
				Message: "id is required",
			},
		}, nil
	}
	feedback, err := s.customerService.GetFeedbackByID(ctx, req.Id)
	if err != nil {
		return &customerpb.GetFeedbackByIDResponse{
			Error: convertErrorToProto(err),
		}, nil
	}
	s.logger.Info("feedback fetched", zap.Any("feedback", feedback))
	return &customerpb.GetFeedbackByIDResponse{
		Feedback: convertFeedbackToProto(feedback),
	}, nil
}

func convertFeedbackToProto(feedback *domain.Feedback) *customerpb.Feedback {
	return &customerpb.Feedback{
		Id:        feedback.ID,
		Rating:    int32(feedback.Rating),
		Comment:   feedback.Comment,
		TaskId:    feedback.TaskID,
		CreatedAt: int32(feedback.CreatedAt.Unix()),
		UpdatedAt: int32(feedback.UpdatedAt.Unix()),
	}
}

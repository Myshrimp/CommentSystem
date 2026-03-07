package biz

import (
	"context"
	v1 "review-service/api/review/v1"
	"review-service/internal/data/model"
	"review-service/pkg/snowflake"

	"github.com/go-kratos/kratos/v2/log"
)

// Review is a Review model.
type Review struct {
	Hello string
}

// ReviewRepo is a Greater repo.
type ReviewRepo interface {
	SaveReview(context.Context, *model.ReviewInfo) (*model.ReviewInfo, error)
	GetReviewsByOrderID(context.Context, int64) ([]*model.ReviewInfo, error)
	GetReview(context.Context, int64) (*model.ReviewInfo, error)
	SaveReply(context.Context, *model.ReviewReplyInfo) (*model.ReviewReplyInfo, error)
	GetReviewReply(context.Context, int64) ([]*model.ReviewReplyInfo, error)
	AuditReview(context.Context, *AuditParam) error
	AppealReview(context.Context, *AppealParam) error
	AuditAppeal(context.Context, *AuditAppealParam) error
	ListReviewByUserID(ctx context.Context, userID int64, offset, limit int) ([]*model.ReviewInfo, error)
}

// ReviewUsecase is a Review usecase.
type ReviewUsecase struct {
	repo ReviewRepo
	log  *log.Helper
}

// NewReviewUsecase new a Review usecase.
func NewReviewUsecase(repo ReviewRepo, logger log.Logger) *ReviewUsecase {
	return &ReviewUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

// CreateReview creates a Review, and returns the new Review.
func (uc *ReviewUsecase) CreateReview(ctx context.Context, review *model.ReviewInfo) (*model.ReviewInfo, error) {
	uc.log.WithContext(ctx).Debugf("[biz] CreateReview: %v", review)
	// 1. validate review
	reviews, err := uc.repo.GetReviewsByOrderID(ctx, review.OrderID)
	if err != nil {
		return nil, v1.ErrorDbFailed("Failed to query reviews by order ID: %v", err)
	}
	if len(reviews) > 0 {
		return nil, v1.ErrorOrderReviewed("Order:%d is already reviewed", review.OrderID)
	}
	// 2. generate review ID
	id := snowflake.GenerateID()
	review.ReviewID = id
	// we can use snowflake or UUID service to generate review ID.
	// in this project, we will use snowflake to generate review ID.
	// 3. query order and product snap-shot info
	// 4. save review
	return uc.repo.SaveReview(ctx, review)
}

func (uc *ReviewUsecase) CreateReply(ctx context.Context, param *ReplyParam) (*model.ReviewReplyInfo, error) {
	// Implementation for creating a review reply
	uc.log.WithContext(ctx).Debugf("[biz] CreateReply: %v", param)
	// 1. validate reply
	// 2. generate reply ID
	id := snowflake.GenerateID()
	reply := &model.ReviewReplyInfo{
		ReviewID: param.ReviewID,
		StoreID:  param.StoreID,
		Content:  param.Content,
		ReplyID:  id,
	}
	// 3. save reply
	return uc.repo.SaveReply(ctx, reply)
}
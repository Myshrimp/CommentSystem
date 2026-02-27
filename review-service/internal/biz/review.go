package biz

import (
	"context"
	"errors"
	"fmt"
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
}

// ReviewUsecase is a Review usecase.
type ReviewUsecase struct {
	repo ReviewRepo
	log *log.Helper
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
		return nil, errors.New("Failed to query database")
	}
	if len(reviews) > 0 {
		return nil, fmt.Errorf("Order:%d is already reviewed", review.OrderID)
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

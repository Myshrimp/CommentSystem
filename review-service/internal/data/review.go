package data

import (
	"context"
	"errors"
	"review-service/internal/biz"
	"review-service/internal/data/model"
	"review-service/internal/data/query"

	"github.com/go-kratos/kratos/v2/log"
)

type reviewRepo struct {
	data *Data
	log  *log.Helper
}

// NewReviewRepo .
func NewReviewRepo(data *Data, logger log.Logger) biz.ReviewRepo {
	return &reviewRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *reviewRepo) SaveReview(ctx context.Context, review *model.ReviewInfo) (*model.ReviewInfo, error) {
	err := r.data.query.ReviewInfo.WithContext(ctx).Save(review)
	if err != nil {
		return nil, err
	}
	return review, nil
}

func (r *reviewRepo) GetReviewsByOrderID(ctx context.Context, orderID int64) ([]*model.ReviewInfo, error) {
	return r.data.query.ReviewInfo.WithContext(ctx).
		Where(r.data.query.ReviewInfo.OrderID.Eq(orderID)).
		Find()
}

func (r *reviewRepo) GetReview(ctx context.Context, reviewID int64) (*model.ReviewInfo, error) {
	return r.data.query.ReviewInfo.WithContext(ctx).
		Where(r.data.query.ReviewInfo.ID.Eq(reviewID)).
		First()
}

func (r *reviewRepo) SaveReply(ctx context.Context, reply *model.ReviewReplyInfo) (*model.ReviewReplyInfo, error) {
	// 1. validate data
	// 1.1. check if reply already exists
	review, err := r.data.query.ReviewInfo.WithContext(ctx).
		Where(r.data.query.ReviewInfo.HasReply.Eq(1)).
		First()
	if err != nil {
		return nil, err
	}

	if review.HasReply == 1 {
		return nil, errors.New("review already has reply")
	}
	// 1.2. rights check, only store owner can reply the review. 
	if review.StoreID != reply.StoreID {
		return nil, errors.New("only store owner can reply the review")
	}
	// 2. update database (update reply table and review table simultaneously, we can use transaction or eventually consistent strategy)
	err = r.data.query.Transaction(func(tx *query.Query) error {
		if _, err := tx.ReviewInfo.WithContext(ctx).Where(tx.ReviewInfo.ReviewID.Eq(reply.ReviewID)).UpdateSimple(tx.ReviewInfo.HasReply.Value(1)); err != nil {
			r.log.WithContext(ctx).Errorf("Failed to update review while saving reply: %v", err)
			return err
		}
		if err := tx.ReviewReplyInfo.WithContext(ctx).Save(reply); err != nil {
			r.log.WithContext(ctx).Errorf("Failed to save review reply: %v", err)
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func (r *reviewRepo) GetReviewReply(ctx context.Context, reviewID int64) ([]*model.ReviewReplyInfo, error) {
	return r.data.query.ReviewReplyInfo.WithContext(ctx).
		Where(r.data.query.ReviewReplyInfo.ReviewID.Eq(reviewID)).
		Find()
}

func (r *reviewRepo) AuditReview(ctx context.Context, param *biz.AuditParam) error {
	_, err := r.data.query.ReviewInfo.WithContext(ctx).
		Where(r.data.query.ReviewInfo.ID.Eq(param.ReviewID)).
		UpdateSimple(r.data.query.ReviewInfo.Status.Value(param.Status))
	return err
}

func (r *reviewRepo) AppealReview(ctx context.Context, param *biz.AppealParam) error {
	return r.SaveAppeal(ctx, &model.ReviewAppealInfo{
		ReviewID: param.ReviewID,
		Content:  param.Content,
		Status:   0, // pending
	})
}

func (r *reviewRepo) SaveAppeal(ctx context.Context, appeal *model.ReviewAppealInfo) error {
	err := r.data.query.ReviewAppealInfo.WithContext(ctx).Save(appeal)
	if err != nil {
		r.log.Errorf("Failed to save review appeal: %v", err)
		return err
	}
	return nil
}

func (r *reviewRepo) AuditAppeal(ctx context.Context, param *biz.AuditAppealParam) error {
	_, err := r.data.query.ReviewAppealInfo.WithContext(ctx).
		Where(r.data.query.ReviewAppealInfo.ReviewID.Eq(param.ReviewID)).
		UpdateSimple(r.data.query.ReviewAppealInfo.Status.Value(param.Status), r.data.query.ReviewAppealInfo.Reason.Value(param.Reason))
	return err
}

func (r *reviewRepo) ListReviewByUserID(ctx context.Context, userID int64, offset, limit int) ([]*model.ReviewInfo, error) {
	return r.data.query.ReviewInfo.WithContext(ctx).
		Where(r.data.query.ReviewInfo.UserID.Eq(userID)).
		Offset(offset).
		Limit(limit).
		Find()
}

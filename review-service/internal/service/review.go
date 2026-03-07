package service

import (
	"context"

	pb "review-service/review-api/review/v1"
	"review-service/internal/biz"
)

type ReviewService struct {
	pb.UnimplementedReviewServer
	uc *biz.ReviewUsecase
}

func NewReviewService(uc *biz.ReviewUsecase) *ReviewService {
	return &ReviewService{
		uc: uc,
	}
}

func (s *ReviewService) CreateReview(ctx context.Context, req *pb.CreateReviewRequest) (*pb.CreateReviewReply, error) {
	return &pb.CreateReviewReply{}, nil
}
func (s *ReviewService) AppealReview(ctx context.Context, req *pb.AppealReviewRequest) (*pb.AppealReviewReply, error) {
	return &pb.AppealReviewReply{}, nil
}
func (s *ReviewService) ReplyReview(ctx context.Context, req *pb.ReplyReviewRequest) (*pb.ReplyReviewReply, error) {
	reply, err := s.uc.CreateReply(ctx, &biz.ReplyParam{
		ReviewID: req.GetReviewID(),
		StoreID:  req.GetStoreID(),
		Content:  req.GetContent(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.ReplyReviewReply{
		ReplyID: reply.ReplyID,
		Succeed: true,
	}, nil
}
func (s *ReviewService) GetReviewByOrderID(ctx context.Context, req *pb.GetReviewByOrderIDRequest) (*pb.GetReviewByOrderIDReply, error) {
	return &pb.GetReviewByOrderIDReply{}, nil
}
func (s *ReviewService) ListReviewByUserID(ctx context.Context, req *pb.ListReviewByUserIDRequest) (*pb.ListReviewByUserIDReply, error) {
	return &pb.ListReviewByUserIDReply{}, nil
}
func (s *ReviewService) AuditAppeal(ctx context.Context, req *pb.AuditAppealRequest) (*pb.AuditAppealReply, error) {
	return &pb.AuditAppealReply{}, nil
}
func (s *ReviewService) AuditReview(ctx context.Context, req *pb.AuditReviewRequest) (*pb.AuditReviewReply, error) {
	return &pb.AuditReviewReply{}, nil
}
func (s *ReviewService) GetReviewReplyList(ctx context.Context, req *pb.GetReviewReplyRequest) (*pb.GetReviewReplyReply, error) {
	return &pb.GetReviewReplyReply{}, nil
}
func (s *ReviewService) GetReview(ctx context.Context, req *pb.GetReviewRequest) (*pb.GetReviewReply, error) {
	return &pb.GetReviewReply{}, nil
}

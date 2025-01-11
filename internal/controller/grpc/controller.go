package grpc

import (
	"context"
	"github.com/KirillGreenev/crypto-rate-service/pkg/logger"
	"go.uber.org/zap"
	"time"

	"github.com/KirillGreenev/crypto-rate-service/internal/controller/grpc/proto"
	"github.com/KirillGreenev/crypto-rate-service/internal/models"
	"github.com/KirillGreenev/crypto-rate-service/internal/service"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type RatesServiceGRPC struct {
	proto.UnimplementedRatesServiceServer
	rs service.RatesService
}

func NewRatesServiceGRPC(rs service.RatesService) *RatesServiceGRPC {
	return &RatesServiceGRPC{rs: rs}
}

func (r *RatesServiceGRPC) GetRates(ctx context.Context, _ *emptypb.Empty) (*proto.RatesResponse, error) {
	log := logger.Logger()
	result, err := r.rs.GetRates(ctx)
	if err != nil {
		log.Error("service.RatesService returned an error", zap.Error(err))
		return nil, models.ErrInternalServer
	}

	return ConvertToGRPC(result), nil
}

func ConvertToGRPC(r models.ResponseService) *proto.RatesResponse {
	timestamp := timestamppb.New(time.Unix(r.Timestamp, 0))
	ask := &proto.Ask{
		Price:  r.Ask.Price,
		Volume: r.Ask.Volume,
		Amount: r.Ask.Amount,
		Factor: r.Ask.Factor,
		Type:   r.Ask.Type,
	}

	bid := &proto.Bid{
		Price:  r.Bid.Price,
		Volume: r.Bid.Volume,
		Amount: r.Bid.Amount,
		Factor: r.Bid.Factor,
		Type:   r.Bid.Type,
	}

	return &proto.RatesResponse{Timestamp: timestamp,
		Ask: ask, Bid: bid}
}

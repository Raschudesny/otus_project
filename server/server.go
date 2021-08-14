package server

import (
	"context"
	"errors"
	"time"

	"github.com/Raschudesny/otus_project/v1/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/Raschudesny/otus_project/v1/server/pb"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
)

//go:generate protoc --proto_path=../api --go_out=pb --go-grpc_out=pb ../api/rotation_service.proto
var _ pb.BannerRotationServiceServer = (*RotationService)(nil)

type RotationService struct {
	pb.UnimplementedBannerRotationServiceServer
	app Application
}

func (r *RotationService) AddSlot(ctx context.Context, request *pb.AddSlotRequest) (*pb.AddSlotResponse, error) {
	if request.GetDescription() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "description argument must be not empty")
	}

	slot, err := r.app.AddSlot(ctx, request.GetDescription())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.AddSlotResponse{Slot: MapSlotToPb(slot)}, nil
}

func (r *RotationService) DeleteSlot(ctx context.Context, request *pb.DeleteSlotRequest) (*pb.DeleteSlotResponse, error) {
	if request.GetSlotId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "slot id argument must be not empty")
	}

	err := r.app.DeleteSlot(ctx, request.GetSlotId())
	switch {
	case errors.Is(err, storage.ErrSlotNotFound):
		return nil, status.Errorf(codes.NotFound, "slot with provided id not found")
	case err != nil:
		return nil, status.Errorf(codes.Internal, err.Error())
	default:
		return &pb.DeleteSlotResponse{}, nil
	}
}

func (r *RotationService) AddBannerToSlot(ctx context.Context, request *pb.AddBannerToSlotRequest) (*pb.AddBannerToSlotResponse, error) {
	panic("implement me")
}

func (r *RotationService) DeleteBannerFromSlot(ctx context.Context, request *pb.DeleteBannerFromSlotRequest) (*pb.DeleteBannerFromSlotResponse, error) {
	panic("implement me")
}

func (r *RotationService) AddBanner(ctx context.Context, request *pb.AddBannerRequest) (*pb.AddBannerResponse, error) {
	panic("implement me")
}

func (r *RotationService) DeleteBanner(ctx context.Context, request *pb.DeleteBannerRequest) (*pb.DeleteBannerResponse, error) {
	panic("implement me")
}

func (r *RotationService) AddGroup(ctx context.Context, request *pb.AddGroupRequest) (*pb.AddGroupResponse, error) {
	panic("implement me")
}

func (r *RotationService) DeleteGroup(ctx context.Context, request *pb.DeleteGroupRequest) (*pb.DeleteGroupResponse, error) {
	panic("implement me")
}

func (r *RotationService) PersistClick(ctx context.Context, request *pb.PersistClickRequest) (*pb.PersistClickResponse, error) {
	panic("implement me")
}

func (r *RotationService) NextBanner(ctx context.Context, request *pb.NextBannerRequest) (*pb.NextBannerResponse, error) {
	panic("implement me")
}

type Server struct {
	Srv  *grpc.Server
	Port int
}

func InitServer() {
	srv := grpc.NewServer(
		grpc.ConnectionTimeout(5*time.Second),
		grpc.UnaryInterceptor(grpc_zap.UnaryServerInterceptor(zap.L())),
		grpc.StreamInterceptor(grpc_zap.StreamServerInterceptor(zap.L())),
	)
	pb.RegisterBannerRotationServiceServer(srv, &RotationService{})
}

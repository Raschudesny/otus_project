package server

import (
	"context"
	"errors"
	"net"
	"strings"
	"time"

	"github.com/Raschudesny/otus_project/v1/server/pb"
	"github.com/Raschudesny/otus_project/v1/storage"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:generate protoc --proto_path=../api --go_out=pb --go-grpc_out=pb ../api/rotation_service.proto
var _ pb.BannerRotationServiceServer = (*RotationService)(nil)

type RotationService struct {
	pb.UnimplementedBannerRotationServiceServer
	app Application
}

func (r *RotationService) AddSlot(ctx context.Context, req *pb.AddSlotRequest) (*pb.AddSlotResponse, error) {
	description := strings.TrimSpace(req.GetDescription())
	if description == "" {
		return nil, status.Errorf(codes.InvalidArgument, "description argument must be not empty")
	}

	slot, err := r.app.AddSlot(ctx, description)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.AddSlotResponse{Slot: MapSlotToPb(slot)}, nil
}

func (r *RotationService) DeleteSlot(ctx context.Context, req *pb.DeleteSlotRequest) (*pb.DeleteSlotResponse, error) {
	slotId := strings.TrimSpace(req.GetSlotId())
	if slotId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "slot id argument must be not empty")
	}

	err := r.app.DeleteSlot(ctx, slotId)
	switch {
	case errors.Is(err, storage.ErrSlotNotFound):
		return nil, status.Errorf(codes.NotFound, "slot with provided id not found")
	case err != nil:
		return nil, status.Errorf(codes.Internal, err.Error())
	default:
		return &pb.DeleteSlotResponse{}, nil
	}
}

//nolint:lll
func (r *RotationService) AddBannerToSlot(ctx context.Context, req *pb.AddBannerToSlotRequest) (*pb.AddBannerToSlotResponse, error) {
	slotId := strings.TrimSpace(req.GetSlotId())
	bannerId := strings.TrimSpace(req.GetBannerId())
	if slotId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "slot id is empty")
	}
	if bannerId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "banner id is empty")
	}
	if err := r.app.AddBannerToSlot(ctx, slotId, bannerId); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add banner to slot: %s", err.Error())
	}
	return &pb.AddBannerToSlotResponse{}, nil
}

//nolint:lll
func (r *RotationService) DeleteBannerFromSlot(ctx context.Context, req *pb.DeleteBannerFromSlotRequest) (*pb.DeleteBannerFromSlotResponse, error) {
	if req.GetSlotId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "slot id is empty")
	}
	if req.GetBannerId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "banner id is empty")
	}
	if err := r.app.DeleteBannerFromSlot(ctx, req.GetSlotId(), req.GetBannerId()); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete banner from slot: %s", err.Error())
	}
	return &pb.DeleteBannerFromSlotResponse{}, nil
}

func (r *RotationService) AddBanner(ctx context.Context, req *pb.AddBannerRequest) (*pb.AddBannerResponse, error) {
	description := strings.TrimSpace(req.GetDescription())
	if description == "" {
		return nil, status.Errorf(codes.InvalidArgument, "description is empty")
	}
	banner, err := r.app.AddBanner(ctx, description)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add banner: %s", err.Error())
	}
	return &pb.AddBannerResponse{Banner: MapBannerToPb(banner)}, nil
}

//nolint:lll
func (r *RotationService) DeleteBanner(ctx context.Context, req *pb.DeleteBannerRequest) (*pb.DeleteBannerResponse, error) {
	bannerId := strings.TrimSpace(req.GetBannerId())
	if bannerId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "banner id is empty")
	}
	if err := r.app.DeleteBanner(ctx, bannerId); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete banner: %s", err.Error())
	}
	return &pb.DeleteBannerResponse{}, nil
}

func (r *RotationService) AddGroup(ctx context.Context, req *pb.AddGroupRequest) (*pb.AddGroupResponse, error) {
	description := strings.TrimSpace(req.GetDescription())
	if description == "" {
		return nil, status.Errorf(codes.InvalidArgument, "description is empty")
	}
	group, err := r.app.AddGroup(ctx, description)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add group: %s", err.Error())
	}
	return &pb.AddGroupResponse{Group: MapGroupToPb(group)}, nil
}

//nolint:lll
func (r *RotationService) DeleteGroup(ctx context.Context, req *pb.DeleteGroupRequest) (*pb.DeleteGroupResponse, error) {
	groupId := strings.TrimSpace(req.GetGroupId())
	if groupId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "group id is empty")
	}
	if err := r.app.DeleteGroup(ctx, groupId); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete group: %s", err.Error())
	}
	return &pb.DeleteGroupResponse{}, nil
}

//nolint:lll
func (r *RotationService) PersistClick(ctx context.Context, req *pb.PersistClickRequest) (*pb.PersistClickResponse, error) {
	slotId := strings.TrimSpace(req.GetSlotId())
	groupId := strings.TrimSpace(req.GetGroupId())
	bannerId := strings.TrimSpace(req.GetBannerId())
	if slotId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "slot id is empty")
	}
	if groupId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "group id is empty")
	}
	if bannerId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "banner id is empty")
	}
	if err := r.app.PersistClick(ctx, slotId, groupId, bannerId); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to persist click: %s", err.Error())
	}
	return &pb.PersistClickResponse{}, nil
}

func (r *RotationService) NextBanner(ctx context.Context, req *pb.NextBannerRequest) (*pb.NextBannerResponse, error) {
	slotId := strings.TrimSpace(req.GetSlotId())
	groupId := strings.TrimSpace(req.GetGroupId())
	if slotId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "slot id is empty")
	}
	if groupId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "group id is empty")
	}
	nextBannerID, err := r.app.NextBannerID(ctx, slotId, groupId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get next banner to show: %s", err.Error())
	}
	return &pb.NextBannerResponse{BannerId: nextBannerID}, nil
}

type Server struct {
	Srv *grpc.Server
}

func InitServer(app Application) *Server {
	srv := grpc.NewServer(
		grpc.ConnectionTimeout(5*time.Second),
		grpc.UnaryInterceptor(grpc_zap.UnaryServerInterceptor(zap.L())),
		grpc.StreamInterceptor(grpc_zap.StreamServerInterceptor(zap.L())),
	)
	pb.RegisterBannerRotationServiceServer(srv, &RotationService{app: app})
	return &Server{
		Srv: srv,
	}
}

// Start function is starting grpc api server on the given port.
// This function is blocking so it must be called in separate goroutine.
// If server start fails, CancelFunc will be called.
func (s Server) Start(host string, port string, cancelFunc context.CancelFunc) {
	// manually calling server shutdown
	defer cancelFunc()

	zap.L().Info("GRPC server starting...", zap.String("address", net.JoinHostPort(host, port)))
	lsn, err := net.Listen("tcp", net.JoinHostPort(host, port))
	if err != nil {
		zap.L().Error("Failed to start grpc server", zap.Error(err))
		return
	}
	if err := s.Srv.Serve(lsn); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		zap.L().Error("Failed to start grpc server", zap.Error(err))
		return
	}
}

func (s Server) Stop() {
	zap.L().Info("GRPC server stopping...")
	s.Srv.GracefulStop()
	zap.L().Info("GRPC server stopped")
}

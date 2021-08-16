package server

import (
	"context"

	storage2 "github.com/Raschudesny/otus_project/v1/internal/storage"
)

type Application interface {
	AddSlot(ctx context.Context, description string) (storage2.Slot, error)
	DeleteSlot(ctx context.Context, slotID string) error
	AddBannerToSlot(ctx context.Context, slotID, bannerID string) error
	DeleteBannerFromSlot(ctx context.Context, bannerID, slotID string) error
	AddBanner(ctx context.Context, description string) (storage2.Banner, error)
	DeleteBanner(ctx context.Context, bannerID string) error
	AddGroup(ctx context.Context, description string) (storage2.SocialGroup, error)
	DeleteGroup(ctx context.Context, groupID string) error
	PersistClick(ctx context.Context, slotID, groupID, bannerID string) error
	NextBannerID(ctx context.Context, slotID, groupID string) (string, error)
}

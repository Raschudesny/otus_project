package server

import (
	"context"

	"github.com/Raschudesny/otus_project/v1/storage"
)

type Application interface {
	AddSlot(ctx context.Context, description string) (storage.Slot, error)
	DeleteSlot(ctx context.Context, slotID string) error
	AddBannerToSlot(ctx context.Context, bannerID, slotID string) error
	DeleteBannerFromSlot(ctx context.Context, bannerID, slotID string) error
	AddBanner(ctx context.Context, description string) (storage.Banner, error)
	DeleteBanner(ctx context.Context, bannerID string)
	AddGroup(ctx context.Context, description string) (storage.SocialGroup, error)
	DeleteGroup(ctx context.Context, groupID string) error
	PersistClick(ctx context.Context, slotID, groupID, bannerID string) error
	NextBanner(ctx context.Context, slotID, groupID string) storage.Banner
}

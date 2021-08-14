package internal

import (
	"context"
	"fmt"
	"math"
	"math/big"

	"github.com/Raschudesny/otus_project/v1/server"
	"github.com/Raschudesny/otus_project/v1/storage"
)

type Repository interface {
	AddSlot(ctx context.Context, description string) (string, error)
	GetSlotByID(ctx context.Context, id string) (storage.Slot, error)
	DeleteSlot(ctx context.Context, id string) error
	AddBanner(ctx context.Context, description string) (string, error)
	GetBannerByID(ctx context.Context, id string) (storage.Banner, error)
	FindBannersBySlot(ctx context.Context, slotID, groupID string) ([]storage.Banner, error)
	DeleteBanner(ctx context.Context, id string) error
	AddBannerToSlot(ctx context.Context, bannerID, slotID string) error
	DeleteBannerFromSlot(ctx context.Context, bannerID, slotID string) error
	AddGroup(ctx context.Context, description string) (string, error)
	DeleteGroup(ctx context.Context, id string) error
	InitStatForBanner(ctx context.Context, slotID, groupID, bannerID string) error
	PersistClick(ctx context.Context, slotID, groupID, bannerID string) error
	PersistShow(ctx context.Context, slotID, groupID, bannerID string) error
	GetShowsAmount(ctx context.Context, slotID, groupID, bannerID string) (int, error)
	GetClicksAmount(ctx context.Context, slotID, groupID, bannerID string) (int, error)
	CountTotalShowsAmount(ctx context.Context, slotID, groupID string) (uint, error)
}

var _ server.Application = (*RotationService)(nil)

type RotationService struct {
	repo Repository
}

func NewRotationService(repository Repository) RotationService {
	return RotationService{repo: repository}
}

func (r RotationService) AddSlot(ctx context.Context, description string) (storage.Slot, error) {
	if description == "" {
		return storage.Slot{}, fmt.Errorf("description param is empty")
	}
	slotID, err := r.repo.AddSlot(ctx, description)
	if err != nil {
		return storage.Slot{}, fmt.Errorf("error during slot creation: %w", err)
	}
	slot, err := r.repo.GetSlotByID(ctx, slotID)
	if err != nil {
		return storage.Slot{}, fmt.Errorf("error during getting slot by id: %w", err)
	}
	return slot, nil
}

func (r RotationService) DeleteSlot(ctx context.Context, slotID string) error {
	if slotID == "" {
		return fmt.Errorf("slotID param is empty")
	}
	if err := r.repo.DeleteSlot(ctx, slotID); err != nil {
		return fmt.Errorf("error during slot deleting: %w", err)
	}
	return nil
}

func (r RotationService) AddBannerToSlot(ctx context.Context, bannerID, slotID string) error {
	if bannerID == "" {
		return fmt.Errorf("bannerID param is empty")
	}
	if slotID == "" {
		return fmt.Errorf("slotID param is empty")
	}
	panic("implement me")
}

func (r RotationService) DeleteBannerFromSlot(ctx context.Context, bannerID, slotID string) error {
	panic("implement me")
}

func (r RotationService) AddBanner(ctx context.Context, description string) (storage.Banner, error) {
	panic("implement me")
}

func (r RotationService) DeleteBanner(ctx context.Context, bannerID string) {
	panic("implement me")
}

func (r RotationService) AddGroup(ctx context.Context, description string) (storage.SocialGroup, error) {
	panic("implement me")
}

func (r RotationService) DeleteGroup(ctx context.Context, groupID string) error {
	panic("implement me")
}

func (r RotationService) PersistClick(ctx context.Context, slotID, groupID, bannerID string) error {
	panic("implement me")
}

func (r RotationService) NextBanner(ctx context.Context, slotID, groupID string) storage.Banner {
	allBanners, err := r.repo.FindBannersBySlot(ctx, slotID, groupID)
	if err != nil {
		panic("world collides")
	}

	for _, banner := range allBanners {
		shows, err := r.repo.GetShowsAmount(ctx, slotID, groupID, banner.Id)
		if err != nil {
			panic("world collides")
		}
		if shows == 0 {
			return banner
		}
	}

	maxTargetValue := 0.0
	maxBannerIndex := 0
	for index, banner := range allBanners {
		bannerClicks, err := r.repo.GetClicksAmount(ctx, slotID, groupID, banner.Id)
		if err != nil {
			panic("world collide")
		}
		bannerShows, err := r.repo.GetShowsAmount(ctx, slotID, groupID, banner.Id)
		if err != nil {
			panic("world collide")
		}
		totalBannerShows, err := r.repo.CountTotalShowsAmount(ctx, slotID, groupID)
		if err != nil {
			panic("world collide")
		}
		targetValue := targetFunction(float64(bannerClicks), float64(bannerShows), float64(totalBannerShows))
		if big.NewFloat(targetValue).Cmp(big.NewFloat(maxTargetValue)) > 0 {
			maxTargetValue = targetValue
			maxBannerIndex = index
		}
	}
	return allBanners[maxBannerIndex]
}

// targetFunction is a maximizing on each step in UCB1 algo function value
// it should be used to evaluate value for each banner.
func targetFunction(clickCount, showCount, totalShowCount float64) float64 {
	avgBannerIncome := clickCount / showCount
	return avgBannerIncome + math.Sqrt((2.0*math.Log(totalShowCount))/showCount)
}

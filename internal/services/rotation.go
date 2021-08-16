package services

import (
	"context"
	"fmt"
	"math"
	"math/big"

	server2 "github.com/Raschudesny/otus_project/v1/internal/server"
	storage2 "github.com/Raschudesny/otus_project/v1/internal/storage"
)

type Repository interface {
	AddSlot(ctx context.Context, description string) (string, error)
	GetSlotByID(ctx context.Context, id string) (storage2.Slot, error)
	DeleteSlot(ctx context.Context, id string) error
	AddBanner(ctx context.Context, description string) (string, error)
	GetBannerByID(ctx context.Context, id string) (storage2.Banner, error)
	DeleteBanner(ctx context.Context, id string) error
	AddBannerToSlot(ctx context.Context, slotID, bannerID string) error
	DeleteBannerFromSlot(ctx context.Context, slotID, bannerID string) error
	AddGroup(ctx context.Context, description string) (string, error)
	GetGroupByID(ctx context.Context, groupID string) (storage2.SocialGroup, error)
	DeleteGroup(ctx context.Context, id string) error
	PersistClick(ctx context.Context, slotID, groupID, bannerID string) error
	PersistShow(ctx context.Context, slotID, groupID, bannerID string) error
	CountTotalShowsAmount(ctx context.Context, slotID, groupID string) (int64, error)
	FindSlotBannerStats(ctx context.Context, slotID, groupID string) ([]storage2.SlotBannerStat, error)
}

var _ server2.Application = (*RotationService)(nil)

type RotationService struct {
	repo Repository
}

func NewRotationService(repository Repository) RotationService {
	return RotationService{repo: repository}
}

func (r RotationService) AddSlot(ctx context.Context, description string) (storage2.Slot, error) {
	if description == "" {
		return storage2.Slot{}, fmt.Errorf("description param is empty")
	}
	slotID, err := r.repo.AddSlot(ctx, description)
	if err != nil {
		return storage2.Slot{}, fmt.Errorf("error during slot creation: %w", err)
	}
	slot, err := r.repo.GetSlotByID(ctx, slotID)
	if err != nil {
		return storage2.Slot{}, fmt.Errorf("error during getting slot by id: %w", err)
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

func (r RotationService) AddBannerToSlot(ctx context.Context, slotID, bannerID string) error {
	if bannerID == "" {
		return fmt.Errorf("bannerID param is empty")
	}
	if slotID == "" {
		return fmt.Errorf("slotID param is empty")
	}
	err := r.repo.AddBannerToSlot(ctx, slotID, bannerID)
	if err != nil {
		return fmt.Errorf("error during adding banner to slot: %w", err)
	}
	return nil
}

func (r RotationService) DeleteBannerFromSlot(ctx context.Context, bannerID, slotID string) error {
	if bannerID == "" {
		return fmt.Errorf("bannerID param is empty")
	}
	if slotID == "" {
		return fmt.Errorf("slotID param is empty")
	}
	if err := r.repo.DeleteBannerFromSlot(ctx, bannerID, slotID); err != nil {
		return fmt.Errorf("errod during deleting banner from slot: %w", err)
	}
	return nil
}

func (r RotationService) AddBanner(ctx context.Context, description string) (storage2.Banner, error) {
	if description == "" {
		return storage2.Banner{}, fmt.Errorf("description param is empty")
	}
	bannerID, err := r.repo.AddBanner(ctx, description)
	if err != nil {
		return storage2.Banner{}, fmt.Errorf("error during creating banner: %w", err)
	}
	banner, err := r.repo.GetBannerByID(ctx, bannerID)
	if err != nil {
		return storage2.Banner{}, fmt.Errorf("error during retrieving created banner: %w", err)
	}
	return banner, nil
}

func (r RotationService) DeleteBanner(ctx context.Context, bannerID string) error {
	if bannerID == "" {
		return fmt.Errorf("bannerID param is empty")
	}
	if err := r.repo.DeleteBanner(ctx, bannerID); err != nil {
		return fmt.Errorf("error during deleting banner")
	}
	return nil
}

func (r RotationService) AddGroup(ctx context.Context, description string) (storage2.SocialGroup, error) {
	if description == "" {
		return storage2.SocialGroup{}, fmt.Errorf("description is empty")
	}
	groupID, err := r.repo.AddGroup(ctx, description)
	if err != nil {
		return storage2.SocialGroup{}, fmt.Errorf("error during adding group: %w", err)
	}
	group, err := r.repo.GetGroupByID(ctx, groupID)
	if err != nil {
		return storage2.SocialGroup{}, fmt.Errorf("error during retrieving froup by id: %w", err)
	}
	return group, nil
}

func (r RotationService) DeleteGroup(ctx context.Context, groupID string) error {
	if groupID == "" {
		return fmt.Errorf("gorupID param is empty")
	}
	if err := r.repo.DeleteGroup(ctx, groupID); err != nil {
		return fmt.Errorf("error during deleting group by id: %w", err)
	}
	return nil
}

func (r RotationService) PersistClick(ctx context.Context, slotID, groupID, bannerID string) error {
	if slotID == "" {
		return fmt.Errorf("slotId param is empty")
	}
	if groupID == "" {
		return fmt.Errorf("groupId param is empty")
	}
	if bannerID == "" {
		return fmt.Errorf("bannerId param is empty")
	}

	if err := r.repo.PersistClick(ctx, slotID, groupID, bannerID); err != nil {
		return fmt.Errorf("failed to persist banner click stats: %w", err)
	}
	return nil
}

// NextBannerID function returns id of a banner which should be shown next
// Implementation function is based on UCB1 algo for a multiarmed bandit problem
// Al the theory behind the scenes can be found in paper below:
// DOI:10.1023/A:1013689704352, Authors: Auer et al., 2002.
// Original link: https://link.springer.com/content/pdf/10.1023/A:1013689704352.pdf.
func (r RotationService) NextBannerID(ctx context.Context, slotID, groupID string) (res string, err error) {
	// saving selected banner id show in db procedure
	defer func() {
		if err == nil {
			if persistErr := r.repo.PersistShow(ctx, slotID, groupID, res); persistErr != nil {
				res, err = "", fmt.Errorf("failed to store banner show: %w", persistErr)
			}
		}
	}()

	bannerStats, err := r.repo.FindSlotBannerStats(ctx, slotID, groupID)
	if err != nil {
		return "", fmt.Errorf("failed to get banner statistics for a slot: %w", err)
	}
	if len(bannerStats) == 0 {
		return "", storage2.ErrNoBannersFoundForSlot
	}

	// UCB1 algo implementation below
	for _, bannerStat := range bannerStats {
		if bannerStat.GetShows() == 0 {
			return bannerStat.BannerID, nil
		}
	}
	totalBannerShows, err := r.repo.CountTotalShowsAmount(ctx, slotID, groupID)
	if err != nil {
		return "", fmt.Errorf("failed to get total banner shows amount: %w", err)
	}

	maxTargetValue := 0.0
	maxBannerID := bannerStats[0].BannerID
	for _, bannerStat := range bannerStats {
		bannerClicks, bannerShows := bannerStat.GetClicks(), bannerStat.GetShows()
		targetValue := targetFunction(float64(bannerClicks), float64(bannerShows), float64(totalBannerShows))
		if big.NewFloat(targetValue).Cmp(big.NewFloat(maxTargetValue)) > 0 {
			maxTargetValue = targetValue
			maxBannerID = bannerStat.BannerID
		}
	}

	return maxBannerID, nil
}

// targetFunction is a maximizing on each step in UCB1 algo function value
// it should be used to evaluate value for each banner.
func targetFunction(clickCount, showCount, totalShowCount float64) float64 {
	avgBannerIncome := clickCount / showCount
	return avgBannerIncome + math.Sqrt((2.0*math.Log(totalShowCount))/showCount)
}
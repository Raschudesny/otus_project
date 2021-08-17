package services_test

import (
	"context"
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/Raschudesny/otus_project/v1/internal/services"
	"github.com/Raschudesny/otus_project/v1/internal/storage"
	"github.com/bxcodec/faker/v3"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type RotationSuite struct {
	suite.Suite
	ctx             context.Context
	cancelFunc      context.CancelFunc
	ctl             *gomock.Controller
	mockRepo        *MockRepository
	rotationService services.RotationService
}

func TestRotationService(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(RotationSuite))
}

func (s *RotationSuite) SetupTest() {
	s.ctx, s.cancelFunc = context.WithTimeout(context.Background(), time.Second*10)
	s.ctl = gomock.NewController(s.T())
	s.mockRepo = NewMockRepository(s.ctl)
	s.rotationService = services.NewRotationService(s.mockRepo)

	// init random function
	seed := time.Now().UnixNano()
	rand.Seed(seed)
	s.T().Logf("rand seed = %d", seed)
}

func (s *RotationSuite) TearDownTest() {
}

func (s RotationSuite) TestAddBanner() {
	testBanner, err := fakeBanner()
	s.Require().NoError(err)

	s.mockRepo.EXPECT().AddBanner(s.ctx, testBanner.Description).Return(testBanner.ID, nil)
	s.mockRepo.EXPECT().GetBannerByID(s.ctx, testBanner.ID).Return(testBanner, nil)

	banner, err := s.rotationService.AddBanner(s.ctx, testBanner.Description)
	s.Require().NoError(err)
	s.Require().Equal(testBanner, banner)
}

func (s RotationSuite) TestDeleteBanner() {
	testID := faker.UUIDHyphenated()
	s.mockRepo.EXPECT().DeleteBanner(s.ctx, testID).Return(nil)

	err := s.rotationService.DeleteBanner(s.ctx, testID)
	s.Require().NoError(err)
}

func (s RotationSuite) TestAddSlot() {
	testSlot, err := fakeSlot()
	s.Require().NoError(err)

	s.mockRepo.EXPECT().AddSlot(s.ctx, testSlot.Description).Return(testSlot.ID, nil)
	s.mockRepo.EXPECT().GetSlotByID(s.ctx, testSlot.ID).Return(testSlot, nil)

	slot, err := s.rotationService.AddSlot(s.ctx, testSlot.Description)
	s.Require().NoError(err)
	s.Require().Equal(testSlot, slot)
}

func (s RotationSuite) TestDeleteSlot() {
	testID := faker.UUIDHyphenated()
	s.mockRepo.EXPECT().DeleteSlot(s.ctx, testID).Return(nil)

	err := s.rotationService.DeleteSlot(s.ctx, testID)
	s.Require().NoError(err)
}

func (s RotationSuite) TestAddGroup() {
	testGroup, err := fakeGroup()
	s.Require().NoError(err)

	s.mockRepo.EXPECT().AddGroup(s.ctx, testGroup.Description).Return(testGroup.ID, nil)
	s.mockRepo.EXPECT().GetGroupByID(s.ctx, testGroup.ID).Return(testGroup, nil)

	group, err := s.rotationService.AddGroup(s.ctx, testGroup.Description)
	s.Require().NoError(err)
	s.Require().Equal(testGroup, group)
}

func (s RotationSuite) TestDeleteGroup() {
	testID := faker.UUIDHyphenated()
	s.mockRepo.EXPECT().DeleteGroup(s.ctx, testID).Return(nil)

	err := s.rotationService.DeleteGroup(s.ctx, testID)
	s.Require().NoError(err)
}

func (s RotationSuite) TestAddBannerToSlot() {
	testSlotID := faker.UUIDHyphenated()
	testBannerID := faker.UUIDHyphenated()

	s.mockRepo.EXPECT().AddBannerToSlot(s.ctx, testSlotID, testBannerID).Return(nil)

	err := s.rotationService.AddBannerToSlot(s.ctx, testSlotID, testBannerID)
	s.Require().NoError(err)
}

func (s RotationSuite) TestDeleteBannerFromSlot() {
	testSlotID := faker.UUIDHyphenated()
	testBannerID := faker.UUIDHyphenated()

	s.mockRepo.EXPECT().DeleteBannerFromSlot(s.ctx, testSlotID, testBannerID).Return(nil)

	err := s.rotationService.DeleteBannerFromSlot(s.ctx, testSlotID, testBannerID)
	s.Require().NoError(err)
}

func (s RotationSuite) TestPersistClick() {
	testSlotID := faker.UUIDHyphenated()
	testBannerID := faker.UUIDHyphenated()
	testGroupID := faker.UUIDHyphenated()

	s.mockRepo.EXPECT().PersistClick(s.ctx, testSlotID, testGroupID, testBannerID).Return(nil)

	err := s.rotationService.PersistClick(s.ctx, testSlotID, testGroupID, testBannerID)
	s.Require().NoError(err)
}

// TestNextBannerIDBasic - test purpose is just to check that NextBannerID correctly works with a storage.
func (s RotationSuite) TestNextBannerIDBasic() {
	testStats := fakeStatsSlice()
	testTotalShowsAmount := countTotalShowsAmount(testStats)

	testSlotID := faker.UUIDHyphenated()
	testGroupID := faker.UUIDHyphenated()

	s.mockRepo.EXPECT().FindSlotBannerStats(s.ctx, testSlotID, testGroupID).Return(testStats, nil)
	s.mockRepo.EXPECT().CountTotalShowsAmount(s.ctx, testSlotID, testGroupID).MaxTimes(1).Return(testTotalShowsAmount, nil)
	s.mockRepo.EXPECT().PersistShow(s.ctx, testSlotID, testGroupID, gomock.Any()).Return(nil)

	_, err := s.rotationService.NextBannerID(s.ctx, testSlotID, testGroupID)
	s.Require().NoError(err)
}

// TestNextBannerIdAllShownAtLeastOnce - test purpose is to check all banners should be shown at least once.
func (s RotationSuite) TestAllBannersShownAtLeastOnce() {
	testStats := fakeStatsSliceWithAllNotShown(100)

	testSlotID := faker.UUIDHyphenated()
	testGroupID := faker.UUIDHyphenated()

	s.mockRepo.EXPECT().FindSlotBannerStats(s.ctx, testSlotID, testGroupID).Times(100).Return(testStats, nil)
	s.mockRepo.EXPECT().CountTotalShowsAmount(s.ctx, testSlotID, testGroupID).Times(0).Return(int64(-1), nil)
	s.mockRepo.EXPECT().PersistShow(
		s.ctx,
		testSlotID,
		testGroupID,
		gomock.Any(),
	).DoAndReturn(func(_ context.Context, _ string, _ string, resBannerId string) error {
		increaseShows(testStats, resBannerId)
		return nil
	}).Times(100)

	for _, v := range testStats {
		id, err := s.rotationService.NextBannerID(s.ctx, testSlotID, testGroupID)
		s.Require().NoError(err)
		s.Require().Equal(v.BannerID, id)
	}
}

func (s RotationSuite) TestMorePopularBannersShownMoreOften() {
	numOfShows := 500
	numOfBanners := 100
	// lets imagine that first 10 banners is a most popular one
	numOfPopularBanners := 10
	testStats := fakeStatsSliceWithAllNotShown(numOfBanners)

	testSlotID := faker.UUIDHyphenated()
	testGroupID := faker.UUIDHyphenated()

	s.mockRepo.EXPECT().FindSlotBannerStats(
		s.ctx,
		testSlotID,
		testGroupID,
	).Times(numOfShows).Return(testStats, nil)
	s.mockRepo.EXPECT().CountTotalShowsAmount(
		s.ctx,
		testSlotID,
		testGroupID,
	).Times(numOfShows - numOfBanners).DoAndReturn(func(_ context.Context, _, _ string) (int64, error) {
		return countTotalShowsAmount(testStats), nil
	})
	s.mockRepo.EXPECT().PersistShow(
		s.ctx,
		testSlotID,
		testGroupID,
		gomock.Any(),
	).DoAndReturn(func(_ context.Context, _ string, _ string, resBannerId string) error {
		increaseShows(testStats, resBannerId)
		return nil
	}).Times(numOfShows)

	s.mockRepo.EXPECT().PersistClick(
		s.ctx,
		testSlotID,
		testGroupID,
		gomock.Any(),
	).DoAndReturn(func(_ context.Context, _ string, _ string, resBannerId string) error {
		increaseClicks(testStats, resBannerId)
		return nil
	}).MaxTimes(numOfShows)

	popularBannersShows := 0
	unpopularBannersShows := 0
	for i := 0; i < numOfShows; i++ {
		id, err := s.rotationService.NextBannerID(s.ctx, testSlotID, testGroupID)
		s.Require().NoError(err)
		if isPopularBanner(testStats, id) {
			popularBannersShows++
		} else {
			unpopularBannersShows++
		}

		//nolint:gosec
		nextClickBanner := rand.Intn(numOfPopularBanners)
		err = s.rotationService.PersistClick(s.ctx, testSlotID, testGroupID, testStats[nextClickBanner].BannerID)
		s.Require().NoError(err)
	}

	s.Require().True(popularBannersShows > unpopularBannersShows)
}

func fakeBanner() (storage.Banner, error) {
	var banner storage.Banner
	err := faker.FakeData(&banner)
	return banner, err
}

func fakeSlot() (storage.Slot, error) {
	var slot storage.Slot
	err := faker.FakeData(&slot)
	return slot, err
}

func fakeGroup() (storage.SocialGroup, error) {
	var group storage.SocialGroup
	err := faker.FakeData(&group)
	return group, err
}

func countTotalShowsAmount(stats []storage.SlotBannerStat) int64 {
	var sum int64
	for _, v := range stats {
		sum += v.GetShows()
	}
	return sum
}

func increaseShows(stats []storage.SlotBannerStat, bannerID string) {
	for ind, v := range stats {
		if v.BannerID == bannerID {
			stats[ind].ShowAmount.Int64 = v.GetShows() + 1
			break
		}
	}
}

func increaseClicks(stats []storage.SlotBannerStat, bannerID string) {
	for ind, v := range stats {
		if v.BannerID == bannerID {
			stats[ind].ClickAmount.Int64 = v.GetClicks() + 1
		}
	}
}

func fakeStatsSlice() []storage.SlotBannerStat {
	res := make([]storage.SlotBannerStat, 0, 100)
	for i := 0; i < 100; i++ {
		res = append(res, storage.SlotBannerStat{
			BannerID: faker.UUIDHyphenated(),
			ClickAmount: sql.NullInt64{
				//nolint:gosec
				Int64: rand.Int63n(10000),
				Valid: true,
			},
			ShowAmount: sql.NullInt64{
				//nolint:gosec
				Int64: rand.Int63n(10000),
				Valid: true,
			},
		})
	}
	return res
}

func fakeStatsSliceWithAllNotShown(numOfBanners int) []storage.SlotBannerStat {
	res := make([]storage.SlotBannerStat, 0, numOfBanners)
	for i := 0; i < numOfBanners; i++ {
		res = append(res, storage.SlotBannerStat{
			BannerID: faker.UUIDHyphenated(),
			ClickAmount: sql.NullInt64{
				Valid: true,
			},
			ShowAmount: sql.NullInt64{
				Valid: true,
			},
		})
	}
	return res
}

func isPopularBanner(stats []storage.SlotBannerStat, bannerID string) bool {
	for ind, v := range stats {
		if v.BannerID == bannerID && ind < 10 {
			return true
		}
	}
	return false
}
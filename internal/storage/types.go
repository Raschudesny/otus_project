package storage

import (
	"database/sql"
	"errors"
)

var (
	ErrSlotNotFound                 = errors.New("slot not found")
	ErrBannerNotFound               = errors.New("banner not found")
	ErrGroupNotFound                = errors.New("social group not found")
	ErrSlotToBannerRelationNotFound = errors.New("slot to banner mapping not found")
	ErrFailedStatsInit              = errors.New("failed to init banner stats")
	ErrNoBannersFoundForSlot        = errors.New("no banners found for provided slot")
)

type Banner struct {
	ID          string `db:"banner_id"`
	Description string `db:"banner_description"`
}

type Slot struct {
	ID          string `db:"slot_id"`
	Description string `db:"slot_description"`
}

type SocialGroup struct {
	ID          string `db:"group_id"`
	Description string `db:"group_description"`
}

type SlotBannerStat struct {
	BannerID    string        `db:"banner_id"`
	ClickAmount sql.NullInt64 `db:"clicks_amount"`
	ShowAmount  sql.NullInt64 `db:"shows_amount"`
}

func (s SlotBannerStat) GetClicks() int64 {
	if !s.ClickAmount.Valid {
		return 0
	}
	return s.ClickAmount.Int64
}

func (s SlotBannerStat) GetShows() int64 {
	if !s.ShowAmount.Valid {
		return 0
	}
	return s.ShowAmount.Int64
}

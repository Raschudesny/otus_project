package storage

import "errors"

var (
	ErrSlotNotFound                 = errors.New("slot not found")
	ErrBannerNotFound               = errors.New("banner not found")
	ErrSlotToBannerRelationNotFound = errors.New("slot to banner mapping not found")
	ErrNoStatsFound                 = errors.New("no stats found")
)

type Banner struct {
	Id          string `db:"banner_id"`
	Description string `db:"banner_description"`
}

type Slot struct {
	Id          string `db:"slot_id"`
	Description string `db:"slot_description"`
}

type SocialGroup struct {
	Id          string `db:"group_id"`
	Description string `db:"group_description"`
}

type SlotBannerRelation struct {
	SlotID   string `db:"slot_id"`
	BannerID string `db:"banner_id"`
}

type BannerStatItem struct {
	BannerID      string `db:"banner_id"`
	SlotID        string `db:"slot_id"`
	SocialGroupID string `db:"group_id"`
	ClickAmount   uint64 `db:"clicks_amount"`
	ShowAmount    uint64 `db:"shows_amount"`
}

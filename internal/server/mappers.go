package server

import (
	pb2 "github.com/Raschudesny/otus_project/v1/internal/server/pb"
	storage2 "github.com/Raschudesny/otus_project/v1/internal/storage"
)

func MapBannerToPb(banner storage2.Banner) *pb2.Banner {
	return &pb2.Banner{
		Id:          banner.ID,
		Description: banner.Description,
	}
}

func MapSlotToPb(slot storage2.Slot) *pb2.Slot {
	return &pb2.Slot{
		Id:          slot.ID,
		Description: slot.Description,
	}
}

func MapGroupToPb(group storage2.SocialGroup) *pb2.Group {
	return &pb2.Group{
		Id:          group.ID,
		Description: group.Description,
	}
}

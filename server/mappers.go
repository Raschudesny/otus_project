package server

import (
	"github.com/Raschudesny/otus_project/v1/server/pb"
	"github.com/Raschudesny/otus_project/v1/storage"
)

func MapBannerToPb(banner storage.Banner) *pb.Banner {
	return &pb.Banner{
		Id:          banner.Id,
		Description: banner.Description,
	}
}

func MapSlotToPb(slot storage.Slot) *pb.Slot {
	return &pb.Slot{
		Id:          slot.Id,
		Description: slot.Description,
	}
}

/*func MapPbBannerToStorage(banner *pb.Banner) storage.Banner {
	return storage.Banner{
		Id:          banner.GetId(),
		Description: banner.GetDescription(),
	}
}*/

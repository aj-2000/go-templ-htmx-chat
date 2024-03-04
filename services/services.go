package services

type AppService struct {
	RoomService *RoomService
}

func NewAppService() *AppService {
	return &AppService{
		RoomService: NewRoomService(),
	}
}

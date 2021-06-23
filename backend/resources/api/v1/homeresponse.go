package v1resources

type HomeResponse struct {
	JoinedRoomList  RoomList            `json:"roomList"`
	RoomRequestList JoinRoomRequestList `json:"roomRequestList"`
}

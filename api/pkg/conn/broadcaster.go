package conn

import (
	"encoding/json"

	"github.com/metaclips/LetsTalk/api/pkg/rooms"
	"github.com/metaclips/LetsTalk/api/pkg/users"
	"github.com/metaclips/LetsTalk/api/values"
	log "github.com/sirupsen/logrus"
)

// broadcastOnlineStatusToAllUserRoom broadcasts users availability status to all users joined rooms.
func broadcastOnlineStatusToAllUserRoom(userHandler *users.UserHandler, roomHandler *rooms.RoomHandler, userEmail string, online bool) {
	user := users.User{Email: userEmail}
	associates, err := userHandler.GetAllUsersAssociates(&user, roomHandler)

	if err != nil {
		log.Errorln("could not get users associate", err)
		return
	}

	values.MapEmailToName.Mutex.RLock()
	for _, assassociateEmail := range associates {
		msg := map[string]interface{}{
			"msgType": values.OnlineStatusMsgType,
			"userID":  userEmail,
			"status":  online,
		}

		if data, err := json.Marshal(msg); err == nil {
			HubConstruct.sendMessage(data, assassociateEmail)
		}
	}

	values.MapEmailToName.Mutex.RUnlock()
}

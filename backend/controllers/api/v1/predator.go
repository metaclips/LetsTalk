package v1controller

import (
	"context"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	apihelpers "github.com/metaclips/LetsTalk/backend/apihelper"
	app_context "github.com/metaclips/LetsTalk/backend/middlewares/context"
	"github.com/metaclips/LetsTalk/backend/models"
	log "github.com/sirupsen/logrus"
)

type PredatorController struct {
	user         models.User
	message      models.Messages
	roomRequest  models.RoomRequest
	ctx          context.Context
	cancelFunc   context.CancelFunc
	checkChannel chan struct{}
}

func NewPredatorController(checkChannel chan struct{}, user models.User, message models.Messages, roomRequest models.RoomRequest) *PredatorController {
	ctx, cancelFunc := app_context.CreateChildCancelContext()

	go func() {
		ticker := time.NewTicker(time.Hour * 24)

		callBack := func() {
			users, err := user.GetRegisteredMinorsPruned(ctx)
			if err != nil {
				log.Errorln("error getting registered minors:", err)
				return
			}

			roomIDs, err := roomRequest.GetUsersRoom(ctx, users.Data)
			if err != nil {
				log.Errorln("error getting registered minors roomIDs:", err)
				return
			}

			for _, roomID := range roomIDs {
				_, err := message.RetrieveMessagesForRoom(ctx, roomID, 0, 0)
				if err != nil {
					log.Errorln("error getting registered minors messages:", err)
					continue
				}
			}
		}

		for {
			select {
			case <-ctx.Done():
				return

			case tick := <-ticker.C:
				callBack()
				log.Infoln("called predator detector at", tick)
			// We triggered a check.
			case <-checkChannel:
				callBack()
				log.Infoln("admin called predator detector at", time.Now())
			}
		}
	}()

	return &PredatorController{
		user, message, roomRequest, ctx, cancelFunc, checkChannel,
	}
}

func (p *PredatorController) StartPredatorSearch(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var failed bool

	select {
	case p.checkChannel <- struct{}{}:
	default:
		failed = true
	}

	if failed {
		apihelpers.Respond(w, http.StatusOK, nil, "Predator Search Already Started", "")
	} else {
		apihelpers.Respond(w, http.StatusOK, nil, "Predator Search Started", "")
	}
}

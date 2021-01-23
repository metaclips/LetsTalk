package users

import (
	"context"

	"github.com/metaclips/LetsTalk/api/pkg/database"
	"github.com/metaclips/LetsTalk/api/pkg/rooms"
	"go.mongodb.org/mongo-driver/bson"
)

type UserHandler struct {
	db  database.Database
	ctx context.Context
}

func NewUsersStore(ctx context.Context, db database.Database) *UserHandler {
	return &UserHandler{
		db:  db,
		ctx: ctx,
	}
}

func (b *UserHandler) GetUser(user *User) error {
	err := b.db.FindOne(b.ctx, bson.M{
		"_id": user.Email,
	}, nil, user)

	return err
}

func (b *UserHandler) GetAllUsersAssociates(user *User, roomHandler *rooms.RoomHandler) ([]string, error) {
	if err := b.GetUser(user); err != nil {
		return nil, err
	}

	usersChannel := make(chan []string)
	done := make(chan struct{})
	users := make([]string, 0)
	registeredUser := make(map[string]bool)

	go func() {
		for {
			data, ok := <-usersChannel
			if ok {
				for _, existingUser := range data {
					if _, exist := registeredUser[existingUser]; !exist && existingUser != user.Email {
						users = append(users, existingUser)
						registeredUser[existingUser] = true
					}
				}

				continue
			}

			close(done)
			break
		}
	}()

	for _, roomJoined := range user.RoomsJoined {
		room, err := roomHandler.FindRoom(roomJoined.RoomID)

		if err != nil {
			close(usersChannel)
			<-done
			return nil, err
		}

		usersChannel <- room.RegisteredUsers
	}

	close(usersChannel)
	<-done

	return users, nil
}

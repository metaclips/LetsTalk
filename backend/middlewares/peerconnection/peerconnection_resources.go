package peerconnection

import (
	"github.com/pion/webrtc/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SDPConstruct struct {
	MsgType        string                     `json:"msgType"`
	ClassSessionID string                     `json:"sessionID"`
	SDP            *webrtc.SessionDescription `json:"sdp"`
	RoomID         primitive.ObjectID         `json:"roomID"`
	UserID         primitive.ObjectID         `json:"userID"`
}

type PeerInterface interface {
	NewPublisherPeerConnection(userID primitive.ObjectID, data []byte, sendMessage func(data interface{}, roomID *primitive.ObjectID, userID *primitive.ObjectID) error)
	NewSubscriberPeerConnection(userID primitive.ObjectID, data []byte, sendMessage func(data interface{}, roomID *primitive.ObjectID, userID *primitive.ObjectID) error)
}

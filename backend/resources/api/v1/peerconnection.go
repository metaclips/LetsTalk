package v1resources

import (
	"sync"

	"github.com/pion/webrtc/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PeerConnection stores all data for all webrtc connections.
type PeerConnection struct {
	VideoTracks         map[string][]*webrtc.TrackLocalStaticRTP
	VideoTrackMutex     *sync.RWMutex
	AudioTracks         map[string][]*webrtc.TrackLocalStaticRTP
	AudioTrackMutex     *sync.RWMutex
	ConnectedUsers      map[string][]primitive.ObjectID
	ConnectedUsersMutex *sync.RWMutex

	PeerConnections      map[primitive.ObjectID]*webrtc.PeerConnection
	PeerConnectionsMutex *sync.RWMutex
}

type SDPConstruct struct {
	ClassSessionID string `json:"sessionID"`
	RoomID         string `json:"roomID"`
	SDP            string `json:"sdp"`
}

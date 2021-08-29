package test

import (
	"testing"
	"time"

	"github.com/metaclips/LetsTalk/backend/helpers/constants"
	"github.com/metaclips/LetsTalk/backend/middlewares/peerconnection"
	v1resources "github.com/metaclips/LetsTalk/backend/resources/api/v1"
	"github.com/pion/webrtc/v3"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// TODO: Cleanup this test. Prune duplicate code.
func TestPublisher(t *testing.T) {
	email := "peerConnection1@students.unilag.edu.ng"

	userData, err := testHandler.seedOneUser(email)
	assert.NoError(t, err)

	ws, err := createWSConn(email)
	assert.NoError(t, err)
	if ws == nil {
		t.FailNow()
	}

	err = ws.SetReadDeadline(time.Now().Add(time.Minute * 30))
	assert.NoError(t, err)

	msg := struct {
		MsgType string `json:"msgType"`
		v1resources.Room
	}{
		MsgType: constants.CreateRoomMsgType,
		Room: v1resources.Room{
			CreateByUserID: userData.ID,
			RoomName:       "My New Room",
		},
	}

	err = ws.WriteJSON(msg)
	assert.NoError(t, err)

	result := v1resources.Room{}
	err = ws.ReadJSON(&result)
	assert.NoError(t, err)

	joinedMsg := v1resources.Message{}
	err = ws.ReadJSON(&joinedMsg)
	assert.NoError(t, err)

	// Create a peerconnection to connects to the publisher
	pc, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	assert.NoError(t, err)

	videoTrack, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeVP8}, "video", "pion")
	assert.NoError(t, err)

	videoTrackSender, err := pc.AddTrack(videoTrack)
	assert.NoError(t, err)

	go func() {
		rtcpBuf := make([]byte, 1500)
		for {
			if _, _, rtcpErr := videoTrackSender.Read(rtcpBuf); rtcpErr != nil {
				return
			}
		}
	}()

	pc.OnICECandidate(func(i *webrtc.ICECandidate) {
		if i == nil {
			offer := peerconnection.SDPConstruct{
				MsgType: constants.StartClassSession,
				SDP:     pc.LocalDescription(),
				RoomID:  joinedMsg.RoomID,
				UserID:  joinedMsg.UserID,
			}

			err = ws.WriteJSON(offer)
			assert.NoError(t, err)
		}
	})

	offer, err := pc.CreateOffer(nil)
	assert.NoError(t, err)

	err = pc.SetLocalDescription(offer)
	assert.NoError(t, err)

	description := peerconnection.SDPConstruct{}
	err = ws.ReadJSON(&description)
	assert.NoError(t, err)

	assert.NotEmpty(t, description.ClassSessionID)

	// Set the remote SessionDescription
	err = pc.SetRemoteDescription(*description.SDP)
	assert.NoError(t, err)

	log.Infoln("gathering complete and webrtc connected")

	classSessionBroadcast := v1resources.Message{}
	err = ws.ReadJSON(&classSessionBroadcast)
	assert.NoError(t, err)

	assert.Equal(t, description.ClassSessionID, classSessionBroadcast.Message)
	// Start sending video packets
}

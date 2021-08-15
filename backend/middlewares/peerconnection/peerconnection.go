package peerconnection

import (
	"encoding/json"
	"errors"
	"io"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/metaclips/LetsTalk/backend/apihelper"
	v1resources "github.com/metaclips/LetsTalk/backend/resources/api/v1"
	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v3"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	rtcpPLIInterval = time.Second * 3
)

type PeerConnection struct {
	*v1resources.PeerConnection
}

func NewPeerConnection() *PeerConnection {
	return &PeerConnection{
		&v1resources.PeerConnection{
			VideoTracks:          make(map[string][]*webrtc.TrackLocalStaticRTP),
			VideoTrackMutex:      &sync.RWMutex{},
			AudioTracks:          make(map[string][]*webrtc.TrackLocalStaticRTP),
			AudioTrackMutex:      &sync.RWMutex{},
			ConnectedUsers:       make(map[string][]primitive.ObjectID),
			ConnectedUsersMutex:  &sync.RWMutex{},
			PeerConnections:      make(map[primitive.ObjectID]*webrtc.PeerConnection),
			PeerConnectionsMutex: &sync.RWMutex{},
		},
	}
}

func (p *PeerConnection) NewPublisherPeerConnection(userID primitive.ObjectID, data []byte, sendMessage func([]byte) error) {
	sessionID := uuid.New().String()
	sdp := v1resources.SDPConstruct{}

	if err := json.Unmarshal(data, &sdp); err != nil {
		// Send back a CreateSessionError indicating user already in a session.
		// TODO: Handle error here
		return
	}

	p.PeerConnectionsMutex.Lock()
	if p.PeerConnections[userID] != nil {
		p.PeerConnectionsMutex.Unlock()
		// TODO: Handle error here
		return
	}

	peerConnection, err := webrtc.NewPeerConnection(apihelper.Config.GetPeerConfig())
	if err != nil {
		return
	}

	p.PeerConnections[userID] = peerConnection
	p.PeerConnectionsMutex.Unlock()

	// Add publisher as a connected user
	p.ConnectedUsersMutex.Lock()
	p.ConnectedUsers[sessionID] = append(p.ConnectedUsers[sessionID], userID)
	p.ConnectedUsersMutex.Unlock()

	// Allow us to receive a video and audio track
	if _, err = peerConnection.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo); err != nil {
		log.Errorln("error creating video transceiver:", err)
		// TODO: Handle error here
		return
	}
	if _, err = peerConnection.AddTransceiverFromKind(webrtc.RTPCodecTypeAudio); err != nil {
		log.Errorln("error creating audio transceiver:", err)
		// TODO: Handle error here
		return
	}

	peerConnection.OnTrack(func(remoteTrack *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		// Send a PLI on an interval so that the publisher is pushing a keyframe every rtcpPLIInterval
		// This can be less wasteful by processing incoming RTCP events, then we would emit a NACK/PLI when a viewer requests it
		if remoteTrack.PayloadType() == webrtc.PayloadType(webrtc.RTPCodecTypeVideo) {
			go func() {
				ticker := time.NewTicker(rtcpPLIInterval)
				for range ticker.C {
					if rtcpSendErr := peerConnection.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{MediaSSRC: uint32(remoteTrack.SSRC())}}); rtcpSendErr != nil {
						log.Errorln("error writing PLI:", rtcpSendErr)
					}
				}
			}()

			defer peerConnection.Close()

			// Create a local track, all our SFU clients will be fed via this track
			localTrack, err := webrtc.NewTrackLocalStaticRTP(remoteTrack.Codec().RTPCodecCapability, "video", "letstalk")
			if err != nil {
				log.Errorln("error getting local track for publisher")
			}

			p.VideoTrackMutex.Lock()
			p.VideoTracks[sessionID] = append(p.VideoTracks[sessionID], localTrack)
			p.VideoTrackMutex.Unlock()

			rtpBuf := make([]byte, 1400)
			for {
				i, _, readErr := remoteTrack.Read(rtpBuf)
				if readErr != nil {
					log.Errorln("error reading remote track for publisher:", err)
					return
				}

				// ErrClosedPipe means we don't have any subscribers, this is ok if no peers have connected yet
				if _, err = localTrack.Write(rtpBuf[:i]); err != nil && !errors.Is(err, io.ErrClosedPipe) {
					log.Errorln("error writing local track for publisher:", err)
					return
				}
			}

		} else if remoteTrack.PayloadType() == webrtc.PayloadType(webrtc.RTPCodecTypeAudio) {
			defer peerConnection.Close()

			// Create a local track, all our SFU clients will be fed via this track
			localTrack, err := webrtc.NewTrackLocalStaticRTP(remoteTrack.Codec().RTPCodecCapability, "audio", "letstalk")
			if err != nil {
				log.Errorln("error getting local track for publisher")
			}

			p.AudioTrackMutex.Lock()
			p.VideoTracks[sessionID] = append(p.AudioTracks[sessionID], localTrack)
			p.AudioTrackMutex.Unlock()

			rtpBuf := make([]byte, 1400)
			for {
				i, _, readErr := remoteTrack.Read(rtpBuf)
				if readErr != nil {
					log.Errorln("error reading remote track for publisher:", err)
					return
				}

				// ErrClosedPipe means we don't have any subscribers, this is ok if no peers have connected yet
				if _, err = localTrack.Write(rtpBuf[:i]); err != nil && !errors.Is(err, io.ErrClosedPipe) {
					log.Errorln("error writing local track for publisher:", err)
					return
				}
			}
		}
	})

	peerConnection.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate == nil {
			// Send peerconnection to client.
			data, err := json.Marshal(peerConnection.LocalDescription())
			if err != nil {
				log.Errorln("error marshalling local description:", err)
				peerConnection.Close()
				return
			}

			if err := sendMessage(data); err != nil {
				log.Errorln("error sending local description:", err)
				peerConnection.Close()
				return
			}

			// Now: Try to construct ICE candidate connection to other peers
			return
		}
	})

	// Set the remote SessionDescription
	if err := p.createNegotiation(peerConnection, sdp); err != nil {
		log.Errorln("error creating negotiation:", err)
		return
	}

}

func (p *PeerConnection) NewSubscriberPeerConnection() {

}

func (p *PeerConnection) createNegotiation(connection *webrtc.PeerConnection, sdp v1resources.SDPConstruct) error {
	err := connection.SetRemoteDescription(
		webrtc.SessionDescription{
			SDP:  sdp.SDP,
			Type: webrtc.SDPTypeOffer,
		})

	if err != nil {
		return err
	}

	answer, err := connection.CreateAnswer(nil)
	if err != nil {
		return err
	}

	if err := connection.SetLocalDescription(answer); err != nil {
		return err
	}

	return nil
}

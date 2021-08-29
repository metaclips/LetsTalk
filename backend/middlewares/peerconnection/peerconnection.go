package peerconnection

import (
	"encoding/json"
	"errors"
	"io"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/metaclips/LetsTalk/backend/apihelper"
	"github.com/metaclips/LetsTalk/backend/helpers/constants"
	v1resources "github.com/metaclips/LetsTalk/backend/resources/api/v1"
	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v3"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	rtcpPLIInterval = time.Second * 3
)

// PeerConnection stores all data for all webrtc connections.
type peerConnection struct {
	videoTracks         map[string][]*webrtc.TrackLocalStaticRTP
	videoTrackMutex     *sync.RWMutex
	audioTracks         map[string][]*webrtc.TrackLocalStaticRTP
	audioTrackMutex     *sync.RWMutex
	connectedUsers      map[string][]primitive.ObjectID
	connectedUsersMutex *sync.RWMutex

	peerConnections      map[primitive.ObjectID]*webrtc.PeerConnection
	peerConnectionsMutex *sync.RWMutex
}

var _ PeerInterface = &peerConnection{}

func NewPionPeerConnection() PeerInterface {
	return &peerConnection{
		videoTracks:          make(map[string][]*webrtc.TrackLocalStaticRTP),
		videoTrackMutex:      &sync.RWMutex{},
		audioTracks:          make(map[string][]*webrtc.TrackLocalStaticRTP),
		audioTrackMutex:      &sync.RWMutex{},
		connectedUsers:       make(map[string][]primitive.ObjectID),
		connectedUsersMutex:  &sync.RWMutex{},
		peerConnections:      make(map[primitive.ObjectID]*webrtc.PeerConnection),
		peerConnectionsMutex: &sync.RWMutex{},
	}
}

func (p *peerConnection) NewPublisherPeerConnection(userID primitive.ObjectID, data []byte, sendMessage func(data interface{}, roomID *primitive.ObjectID, userID *primitive.ObjectID) error) {
	sessionID := uuid.New().String()
	sdp := SDPConstruct{}

	if err := json.Unmarshal(data, &sdp); err != nil {
		// Send back a CreateSessionError indicating user already in a session.
		// TODO: Handle error here
		return
	}

	p.peerConnectionsMutex.Lock()
	if p.peerConnections[userID] != nil {
		p.peerConnectionsMutex.Unlock()
		// TODO: Handle error here
		return
	}

	peerConnection, err := webrtc.NewPeerConnection(apihelper.Config.GetPeerConfig())
	if err != nil {
		return
	}

	p.peerConnections[userID] = peerConnection
	p.peerConnectionsMutex.Unlock()

	// Add publisher as a connected user
	p.connectedUsersMutex.Lock()
	p.connectedUsers[sessionID] = append(p.connectedUsers[sessionID], userID)
	p.connectedUsersMutex.Unlock()

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
				return
			}

			p.videoTrackMutex.Lock()
			p.videoTracks[sessionID] = append(p.videoTracks[sessionID], localTrack)
			p.videoTrackMutex.Unlock()

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
				return
			}

			p.audioTrackMutex.Lock()
			p.audioTracks[sessionID] = append(p.audioTracks[sessionID], localTrack)
			p.audioTrackMutex.Unlock()

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
			sdpConstruct := SDPConstruct{
				MsgType:        constants.OnICECandidate,
				ClassSessionID: sessionID,
				UserID:         userID,
				SDP:            peerConnection.LocalDescription(),
				RoomID:         sdp.RoomID,
			}

			if err := sendMessage(sdpConstruct, nil, &userID); err != nil {
				log.Errorln("error sending local description:", err)
				peerConnection.Close()
				return
			}

			// Also send a PC publisher join message to everyone in the room.
			publisherBroadcast := v1resources.Message{
				Message: sessionID,
				Type:    constants.MessageTypeClassSession,
				RoomID:  sdp.RoomID,
				UserID:  userID,
			}

			if err := sendMessage(publisherBroadcast, &sdp.RoomID, &userID); err != nil {
				log.Errorln("error sending local description:", err)
				peerConnection.Close()
				return
			}

			return
		}

		// Now: Try to construct ICE candidate connection to other peers
		p.connectedUsersMutex.RLock()
		connectedUsers, ok := p.connectedUsers[sessionID]
		if !ok {
			log.Errorln("connected user for", sessionID, "closed")
			return
		}
		p.connectedUsersMutex.RUnlock()

		p.peerConnectionsMutex.RLock()
		for _, user := range connectedUsers {
			if user == userID {
				continue
			}

			peerConnection, ok := p.peerConnections[user]
			if !ok {
				log.Errorln("peer connection for user with session ID", sessionID, "not found:", err)
				continue
			}

			if err := peerConnection.AddICECandidate(candidate.ToJSON()); err != nil {
				log.Errorln("error adding ice candidate:", err)
			}
		}
		p.peerConnectionsMutex.RUnlock()
	})

	peerConnection.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
		if state == webrtc.PeerConnectionStateClosed {
			log.Warning("publisher for sessionID", sessionID, "in closed state")

			// Close all other peer connections on the network.
			p.connectedUsersMutex.Lock()
			users, ok := p.connectedUsers[sessionID]
			delete(p.connectedUsers, sessionID)
			p.connectedUsersMutex.Unlock()
			if !ok {
				log.Errorln("connected user for", sessionID, "closed")
				return
			}

			p.peerConnectionsMutex.Lock()
			for _, user := range users {
				pc, ok := p.peerConnections[user]
				if !ok {
					log.Errorln("error retrieving users peeerconnection:", user)
					continue
				}
				if err := pc.Close(); err != nil {
					log.Errorln("error closing peer connection:", err)
				}
				delete(p.peerConnections, user)
			}
			p.peerConnectionsMutex.Unlock()
		} else if state == webrtc.PeerConnectionStateDisconnected {
			// Try to connect the network again.
			log.Warning("publisher for sessionID", sessionID, "in disconnected state")
		} else if state == webrtc.PeerConnectionStateConnected {
			// Send the link to others on the chat indicating that there's a connection.
			// Save message.
		} else {
			log.Warning("publisher for sessionID", sessionID, "in", state, "state")
		}
	})

	// Set the remote SessionDescription
	if err := p.createNegotiation(peerConnection, sdp); err != nil {
		log.Errorln("error creating negotiation:", err)
		return
	}

}

func (p *peerConnection) NewSubscriberPeerConnection(userID primitive.ObjectID, data []byte, sendMessage func(data interface{}, roomID *primitive.ObjectID, userID *primitive.ObjectID) error) {
	sdp := SDPConstruct{}

	if err := json.Unmarshal(data, &sdp); err != nil {
		// Send back a CreateSessionError indicating user already in a session.
		// TODO: Handle error here
		return
	}

	sessionID := sdp.ClassSessionID

	// COPIED!!!!!
	p.peerConnectionsMutex.Lock()
	if p.peerConnections[userID] != nil {
		p.peerConnectionsMutex.Unlock()
		// TODO: Handle error here
		return
	}

	peerConnection, err := webrtc.NewPeerConnection(apihelper.Config.GetPeerConfig())
	if err != nil {
		return
	}

	p.peerConnections[userID] = peerConnection
	p.peerConnectionsMutex.Unlock()

	// Add publisher as a connected user
	p.connectedUsersMutex.Lock()
	p.connectedUsers[sdp.ClassSessionID] = append(p.connectedUsers[sdp.ClassSessionID], userID)
	p.connectedUsersMutex.Unlock()

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
		if remoteTrack.PayloadType() == webrtc.PayloadType(webrtc.RTPCodecTypeAudio) {
			defer peerConnection.Close()

			// Create a local track, all our SFU clients will be fed via this track
			localTrack, err := webrtc.NewTrackLocalStaticRTP(remoteTrack.Codec().RTPCodecCapability, "audio", "letstalk")
			if err != nil {
				log.Errorln("error getting local track for publisher")
				return
			}

			p.connectedUsersMutex.RLock()
			users, ok := p.connectedUsers[sessionID]
			if !ok {
				log.Errorln("publisher session exited")
				return
			}
			p.connectedUsersMutex.RUnlock()

			p.peerConnectionsMutex.RLock()
			for _, user := range users {
				if user == userID {
					continue
				}

				pc, ok := p.peerConnections[user]
				if !ok {
					log.Errorln("peerconnection specified invalid")
					continue
				}

				if _, err := pc.AddTrack(localTrack); err != nil {
					log.Errorln("error adding user to track:", err)
					continue
				}
			}
			p.peerConnectionsMutex.RUnlock()

			p.audioTrackMutex.Lock()
			p.audioTracks[sdp.ClassSessionID] = append(p.audioTracks[sdp.ClassSessionID], localTrack)
			p.audioTrackMutex.Unlock()

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
			sdpConstruct := SDPConstruct{
				MsgType:        constants.OnICECandidate,
				ClassSessionID: sessionID,
				UserID:         userID,
				SDP:            peerConnection.LocalDescription(),
				RoomID:         sdp.RoomID,
			}

			d, err := json.Marshal(sdpConstruct)
			if err != nil {
				log.Errorln("error marshalling local description:", err)
				peerConnection.Close()
				return
			}

			if err := sendMessage(d, nil, &userID); err != nil {
				log.Errorln("error sending local description:", err)
				peerConnection.Close()
				return
			}

			return
		}

		// Now: Try to construct ICE candidate connection to other peers
		p.connectedUsersMutex.RLock()
		connectedUsers, ok := p.connectedUsers[sessionID]
		if !ok {
			log.Errorln("connected user for", sessionID, "closed")
			return
		}
		p.connectedUsersMutex.RUnlock()

		p.peerConnectionsMutex.RLock()
		for _, user := range connectedUsers {
			if user == userID {
				continue
			}

			peerConnection, ok := p.peerConnections[user]
			if !ok {
				log.Errorln("peer connection for user with session ID", sessionID, "not found:", err)
				continue
			}

			if err := peerConnection.AddICECandidate(candidate.ToJSON()); err != nil {
				log.Errorln("error adding ice candidate:", err)
			}
		}
		p.peerConnectionsMutex.RUnlock()
	})

	peerConnection.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
		if state == webrtc.PeerConnectionStateClosed {
			log.Warning("publisher for sessionID", sessionID, "in closed state")

			// Close all other peer connections on the network.
			p.connectedUsersMutex.Lock()
			users, ok := p.connectedUsers[sessionID]
			for i := range users {
				if users[i] == userID {
					users = append(users[:i], users[i+1:]...)
				}
			}
			p.connectedUsersMutex.Unlock()
			if !ok {
				log.Errorln("connected user for", sessionID, "closed")
				return
			}

			p.peerConnectionsMutex.Lock()
			delete(p.peerConnections, userID)
			p.peerConnectionsMutex.Unlock()

			if err := peerConnection.Close(); err != nil {
				log.Errorln("error closing peer connection:", err)
			}

		} else if state == webrtc.PeerConnectionStateDisconnected {
			// Try to connect the network again.
			log.Warning("publisher for sessionID", sessionID, "in disconnected state")
		} else if state == webrtc.PeerConnectionStateConnected {
			log.Infoln("user with sessionID", sessionID, "connected")
		} else {
			log.Warning("publisher for sessionID", sessionID, "in", state, "state")
		}
	})

	// Set the remote SessionDescription
	if err := p.createNegotiation(peerConnection, sdp); err != nil {
		log.Errorln("error creating negotiation:", err)
		return
	}
}

func (p *peerConnection) createNegotiation(connection *webrtc.PeerConnection, sdp SDPConstruct) error {
	err := connection.SetRemoteDescription(
		webrtc.SessionDescription{
			SDP:  sdp.SDP.SDP,
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

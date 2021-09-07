package test

// // TODO: Cleanup this test. Prune duplicate code.
// func TestPublisher(t *testing.T) {
// 	email := "peerConnection1@students.unilag.edu.ng"

// 	ws, joinedMsg := createRoom(t, email)

// 	// Create a peerconnection to connects to the publisher
// 	pc, err := webrtc.NewPeerConnection(webrtc.Configuration{})
// 	assert.NoError(t, err)

// 	videoTrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeVP8}, "video", "pion")
// 	assert.NoError(t, err)

// 	videoTrackSender, err := pc.AddTrack(videoTrack)
// 	assert.NoError(t, err)

// 	go func() {
// 		rtcpBuf := make([]byte, 1500)
// 		for {
// 			if _, _, rtcpErr := videoTrackSender.Read(rtcpBuf); rtcpErr != nil {
// 				return
// 			}
// 		}
// 	}()

// 	offer, err := pc.CreateOffer(nil)
// 	assert.NoError(t, err)

// 	err = pc.SetLocalDescription(offer)
// 	assert.NoError(t, err)

// 	pc.OnICECandidate(func(i *webrtc.ICECandidate) {
// 		log.Infoln(i)
// 	})

// 	gatherer := webrtc.GatheringCompletePromise(pc)
// 	<-gatherer

// 	err = ws.WriteJSON(peerconnection.SDPConstruct{
// 		MsgType: constants.StartClassSession,
// 		SDP:     pc.LocalDescription(),
// 		RoomID:  joinedMsg.RoomID,
// 		UserID:  joinedMsg.UserID,
// 	})
// 	assert.NoError(t, err)

// 	// Wait for OnICECandidate to fire.
// 	description := peerconnection.SDPConstruct{}
// 	err = ws.ReadJSON(&description)
// 	assert.NoError(t, err)

// 	assert.NotEmpty(t, description.ClassSessionID)

// 	// Set the remote SessionDescription
// 	err = pc.SetRemoteDescription(*description.SDP)
// 	assert.NoError(t, err)

// 	log.Infoln("gathering complete and webrtc connected")

// 	classSessionBroadcast := v1resources.Message{}
// 	err = ws.ReadJSON(&classSessionBroadcast)
// 	assert.NoError(t, err)

// 	assert.Equal(t, description.ClassSessionID, classSessionBroadcast.Message)

// 	conn := make(chan struct{})
// 	pc.OnConnectionStateChange(func(pcs webrtc.PeerConnectionState) {
// 		log.Errorln("pcs is", pcs)
// 		if !assert.NotEqual(t, webrtc.PeerConnectionStateFailed, pcs) {
// 			pc.SignalingState()
// 			panic("peerconnection failed")
// 		}

// 		if pcs == webrtc.PeerConnectionStateConnected {
// 			log.Infoln("Ok connected")
// 			conn <- struct{}{}
// 		}
// 	})

// 	<-conn

// 	go func() {
// 		close := time.Tick(time.Second * 10)

// 		for {
// 			select {
// 			case <-time.After(20 * time.Millisecond):
// 				assert.NoError(t, videoTrack.WriteSample(media.Sample{Data: []byte{0x00}, Duration: time.Second}))
// 			case <-close:
// 				return
// 			}
// 		}
// 	}()

// 	assert.NoError(t, err)
// }

// // func TestSubscriber(t *testing.T) {
// // 	email := "peerConnection2@students.unilag.edu.ng"

// // 	ws, joinedMsg := createRoom(t, email)
// // }

// func createRoom(t *testing.T, email string) (*websocket.Conn, *v1resources.Message) {
// 	userData, err := testHandler.seedOneUser(email)
// 	assert.NoError(t, err)

// 	ws, err := createWSConn(email)
// 	assert.NoError(t, err)
// 	if ws == nil {
// 		t.FailNow()
// 	}

// 	err = ws.SetReadDeadline(time.Now().Add(time.Minute * 30))
// 	assert.NoError(t, err)

// 	msg := struct {
// 		MsgType string `json:"msgType"`
// 		v1resources.Room
// 	}{
// 		MsgType: constants.CreateRoomMsgType,
// 		Room: v1resources.Room{
// 			CreateByUserID: userData.ID,
// 			RoomName:       "My New Room",
// 		},
// 	}

// 	err = ws.WriteJSON(msg)
// 	assert.NoError(t, err)

// 	result := v1resources.Room{}
// 	err = ws.ReadJSON(&result)
// 	assert.NoError(t, err)

// 	joinedMsg := v1resources.Message{}
// 	err = ws.ReadJSON(&joinedMsg)
// 	assert.NoError(t, err)

// 	return ws, &joinedMsg
// }

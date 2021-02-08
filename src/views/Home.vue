<template>
  <v-row no-gutters class="fill-height">
    <v-col cols="12" v-if="callUI">
      <CallUI
        :endCallSession="endCallSession"
        :changeAudio="changeAudio"
        :changeVideo="changeVideo"
        :isCallPublisher="isPublisher"
        :startDesktopSharing="startDesktopSharing"
        :stopDesktopSharing="stopDesktopSharing"
        :connectedUsers="classSessionConnectedUsers"
      />
    </v-col>

    <v-col cols="12" v-else>
      <v-row no-gutters class="fill-height">
        <v-col :cols="$vuetify.breakpoint.smAndDown ? '2' : '1'">
          <SideBar
            :activateAddRoomDialog="activateAddRoomDialog"
            :activateNotificationDialog="activateNotificationDialog"
            :deactivateAllDialogs="deactivateAllDialogs"
            :unreadNotifications="unreadNotificationsCount"
            :unreadRoomMessages="unreadRoomMessageCount"
            :showAvailableRooms="showAvailableRooms"
          />
        </v-col>

        <v-col :cols="$vuetify.breakpoint.smAndDown ? '10' : '11'">
          <v-row
            style="height: 100%"
            align="center"
            justify="center"
            v-if="showNotificationDialog || showAddRoomDialog"
          >
            <v-card
              v-if="showNotificationDialog"
              class="mx-5"
              height="min-content"
              outlined
              shaped
            >
              <v-row class="mx-5">
                <v-col cols="12">
                  <v-card-title>
                    <b>Notifications</b>
                  </v-card-title>
                </v-col>

                <v-col cols="12">
                  <v-divider></v-divider>
                </v-col>

                <v-col cols="12">
                  <v-card-text>
                    <v-container
                      fluid
                      style="max-height: 72vh"
                      class="overflow-y-auto"
                    >
                      <v-card-subtitle v-if="joinRequests.length == 0"
                        >No Notifications</v-card-subtitle
                      >
                      <v-card
                        v-for="(joinRequest, index) in joinRequests"
                        :key="index"
                        flat
                      >
                        <v-card-text
                          >{{ joinRequest.requestingUserName }} [{{
                            joinRequest.requestingUserID
                          }}] wants you to join room
                          {{ joinRequest.roomName }}</v-card-text
                        >
                        <v-card-actions>
                          <v-btn
                            @click="
                              joinRoom(
                                joinRequest.requestingUserID,
                                joinRequest.roomID,
                                joinRequest.roomName,
                                true,
                                index
                              )
                            "
                            text
                            color="green"
                            >Join</v-btn
                          >

                          <v-btn
                            @click="
                              joinRoom(
                                joinRequest.requestingUserID,
                                joinRequest.roomID,
                                joinRequest.roomName,
                                false,
                                index
                              )
                            "
                            text
                            color="red"
                            >Reject</v-btn
                          >
                        </v-card-actions>
                      </v-card>
                    </v-container>
                  </v-card-text>
                </v-col>
              </v-row>
            </v-card>

            <v-card
              class="mx-5"
              flat
              outlined
              shaped
              v-else-if="showAddRoomDialog"
            >
              <v-row class="mx-5">
                <v-col cols="12">
                  <v-card-text>
                    <b>Create New Room</b>
                  </v-card-text>
                </v-col>

                <v-col cols="12">
                  <v-text-field
                    @keyup.enter.exact="createRoom()"
                    placeholder="Room Name"
                    rounded
                    outlined
                    v-model="newRoomName"
                  />
                  <v-btn @click="createRoom()" text color="green"
                    >Create Room</v-btn
                  >
                </v-col>
              </v-row>
            </v-card>
          </v-row>

          <v-row v-else-if="$vuetify.breakpoint.smAndDown" no-gutters>
            <RoomsPage
              v-if="showAvailableRoomsPage"
              :indexOfCurrentViewedRoom="indexOfCurrentViewedRoom"
              :recentChatPreview="recentChatPreview"
              :changeViewedRoomIndex="changeViewedRoomIndex"
              :joinedRooms="joinedRooms"
              :sendWSMessage="sendWSMessage"
              :unreadRoomMessages="unreadRoomMessages"
            />

            <ChatPage
              v-else
              :userID="userID"
              :names="names"
              :initiateFile="initiateFile"
              :fileUploadDownload="fileUploadDownload"
              :associates="usersOnline"
              :clearFetchedUsers="clearFetchedUsers"
              :fetchedUsers="fetchedUsers"
              :sendWSMessage="sendWSMessage"
              :currentViewedRoom="currentViewedRoom"
              :changeDownloadStatus="changeDownloadStatus"
              :startCallSession="startCallSession"
              :joinCallSession="joinCallSession"
            />
          </v-row>

          <v-row v-else no-gutters>
            <v-col cols="3">
              <RoomsPage
                :indexOfCurrentViewedRoom="indexOfCurrentViewedRoom"
                :recentChatPreview="recentChatPreview"
                :changeViewedRoomIndex="changeViewedRoomIndex"
                :joinedRooms="joinedRooms"
                :sendWSMessage="sendWSMessage"
                :unreadRoomMessages="unreadRoomMessages"
              />
            </v-col>

            <v-col cols="9">
              <ChatPage
                v-if="showChatPage"
                :userID="userID"
                :names="names"
                :initiateFile="initiateFile"
                :fileUploadDownload="fileUploadDownload"
                :associates="usersOnline"
                :clearFetchedUsers="clearFetchedUsers"
                :fetchedUsers="fetchedUsers"
                :sendWSMessage="sendWSMessage"
                :currentViewedRoom="currentViewedRoom"
                :changeDownloadStatus="changeDownloadStatus"
                :startCallSession="startCallSession"
                :joinCallSession="joinCallSession"
              />
            </v-col>
          </v-row>
        </v-col>
      </v-row>
    </v-col>

    <v-dialog persistent width="400px" v-model="dialog">
      <v-card>
        <v-card-text>
          <div align="center" justify="center">
            <v-row>
              <v-col cols="12">
                <span>{{ dialogError }}</span>
              </v-col>
              <v-col>
                <v-row justify="center">
                  <v-btn text depressed color="green" @click="dialog = false">
                    Close
                  </v-btn>
                  <v-btn
                    depressed
                    text
                    color="red"
                    v-if="socketClosed"
                    @click="connectWebsocket"
                  >
                    Reconnect
                  </v-btn>
                </v-row>
              </v-col>
            </v-row>
          </div>
        </v-card-text>
      </v-card>
    </v-dialog>
  </v-row>
</template>

<script lang="ts">
"use strict";

let socket: WebSocket;
let videoTransceiver: RTCRtpTransceiver;
let audioTransceiver: RTCRtpTransceiver;

let videoTrack: MediaStreamTrack;
let audioTrack: MediaStreamTrack;
let stream: MediaStream;

let peerConnection: RTCPeerConnection;

// @ is an alias to /src
import Vue from "vue";
import CryptoJS from "crypto-js";

import { WSMessageType, MessageType, DefaultChunkSize } from "./Constants";

import SideBar from "../components/SideBar.vue";
import RoomsPage from "../components/RoomsPage.vue";
import ChatPage from "../components/ChatPage.vue";
import CallUI from "../components/CallUI.vue";

import {
  JoinedRoom,
  RoomPageDetails,
  Message,
  JoinRequest,
  SentRoomRequest,
  FetchedUsers,
  UnreadRooms,
  RecentChatPreview,
  UsersOnline,
  ExitRoomDetails,
  FileUploadDownloadDetails,
  FileUploadComplete,
  FileDownload,
} from "./Types";

export default Vue.extend({
  name: "Home",
  components: {
    SideBar,
    RoomsPage,
    ChatPage,
    CallUI,
  },

  data: () => ({
    joinedRooms: [] as JoinedRoom[],
    currentViewedRoom: {} as RoomPageDetails,
    joinRequests: [] as JoinRequest[],
    fetchedUsers: [] as FetchedUsers[],
    unreadRoomMessages: {} as UnreadRooms,
    recentChatPreview: {} as RecentChatPreview,
    usersOnline: {} as UsersOnline,
    classSessionConnectedUsers: {} as UsersOnline,
    fileUploadDownload: {} as FileDownload,

    userID: "",
    names: "",
    newRoomName: "",
    currentDownloadRoomID: "",
    dialogError: "",

    showAddRoomDialog: false,
    showNotificationDialog: false,
    showChatPage: false,
    isFile: false,
    callUI: false,
    isPublisher: false,
    socketClosed: false,
    dialog: false,
    showAvailableRoomsPage: true,

    indexOfCurrentViewedRoom: 0,
    unreadRoomMessageCount: 0,
    unreadNotificationsCount: 0,
  }),

  methods: {
    onUnAuthorizedAccess: function () {
      this.redirectToLogin();
    },

    redirectToLogin: function () {
      window.location.replace("/login");
    },

    activateAddRoomDialog: function () {
      this.showAddRoomDialog = true;
      this.showNotificationDialog = false;
    },

    activateNotificationDialog: function () {
      this.showNotificationDialog = true;
      this.showAddRoomDialog = false;
    },

    deactivateAllDialogs: function () {
      this.showAddRoomDialog = false;
      this.showNotificationDialog = false;
    },

    // Switch between available room page and chat page if user is on mobile.
    showAvailableRooms: function () {
      this.showAvailableRoomsPage = true;
    },

    onWebsocketOpen: function (
      joinedRooms: JoinedRoom[],
      joinRequests: JoinRequest[]
    ) {
      if (joinedRooms) {
        this.joinedRooms = joinedRooms;
      }

      if (joinRequests) {
        this.joinRequests = joinRequests;
        this.unreadNotificationsCount = joinRequests.length;
        console.log(this.unreadNotificationsCount);
      }
    },

    onRequestMessages: function (roomDetails: RoomPageDetails) {
      if (roomDetails.firstLoad) {
        console.log(roomDetails.registeredUsers);
        this.currentViewedRoom = roomDetails;
        if (!this.currentViewedRoom.messages) {
          this.currentViewedRoom.messages = [];
        }

        this.joinedRooms[
          this.indexOfCurrentViewedRoom
        ].roomIcon = this.currentViewedRoom.roomIcon;
      } else {
        for (let i = roomDetails.messages.length - 1; i >= 0; i--) {
          this.currentViewedRoom.messages.unshift(roomDetails.messages[i]);
        }
      }

      this.scrollToBottomOfChatPage();
    },

    onJoinRoom: function (joinedRoom: JoinedRoom) {
      if (joinedRoom.userID === this.userID) {
        this.joinedRooms.unshift(joinedRoom);
        this.indexOfCurrentViewedRoom++;
        return;
      } else if (this.currentViewedRoom.roomID === joinedRoom.roomID) {
        let message: string = joinedRoom.userID + " Rejected join request.";
        if (joinedRoom.joined) {
          message = joinedRoom.userID + " Accepted join request.";
        }

        this.currentViewedRoom.messages.push({
          time: "",
          type: MessageType.Info,
          name: "",
          userID: "",
          roomID: "",
          size: "",
          hash: "",
          message: message,
          index: this.currentViewedRoom.messages.length + 1,
        });
      }

      this.usersOnline[joinedRoom.userID] = {
        name: joinedRoom.name,
        isOnline: true,
      };
    },

    onNewMessage: function (message: Message) {
      this.updateRecentMessagePreview(message.roomID, message.message);
      if (this.currentViewedRoom.roomID === message.roomID) {
        this.currentViewedRoom.messages.push(message);
        this.scrollToBottomOfChatPage();
        return;
      }

      this.changeNumberOfUnreadNotification(message.roomID, true);
      // ToDo: add notification sound
    },

    onSentRoomRequest: function (sentRequest: SentRoomRequest) {
      const message =
        sentRequest.userRequested +
        " Was requested to join the room by " +
        sentRequest.requesterID;

      if (this.userID === sentRequest.userRequested) {
        this.joinRequests.unshift({
          requestingUserName: sentRequest.requesterName,
          requestingUserID: sentRequest.requesterID,
          roomID: sentRequest.roomID,
          roomName: sentRequest.roomName,
        });

        this.unreadNotificationsCount = this.joinRequests.length;
        return;
      } else if (this.currentViewedRoom.roomID == sentRequest.roomID) {
        this.currentViewedRoom.messages.push({
          type: MessageType.Info,
          message: message,
          size: "",
          hash: "",
          time: "",
          name: "",
          userID: "",
          roomID: "",
          index: this.currentViewedRoom.messages.length + 1,
        });

        console.log("sent notification");
      }

      this.updateRecentMessagePreview(sentRequest.roomID, message);
    },

    onExitRoom: function (exitRoom: ExitRoomDetails) {
      console.log(exitRoom);
      const message: Message = {
        time: "",
        name: "",
        size: "",
        hash: "",
        type: MessageType.Info,
        userID: "",
        roomID: "",
        message: "You left the room.",
        index: this.currentViewedRoom.messages.length + 1,
      };

      if (exitRoom.roomID == this.currentViewedRoom.roomID) {
        if (exitRoom.userID != this.userID)
          message.message = exitRoom.userID + " Left the room.";

        this.currentViewedRoom.messages.push(message);
      }
    },

    updateRecentMessagePreview: function (roomID: string, message: string) {
      this.recentChatPreview[roomID] = message;
    },

    updateChatRoomPage: function () {
      setTimeout(() => {
        this.$children[2].$forceUpdate();
      }, 0);
    },

    updateCallPage: function () {
      setTimeout(() => {
        this.$children[3].$forceUpdate();
      }, 0);
    },

    scrollToBottomOfChatPage: function () {
      this.$nextTick(() => {
        const scrollHeight = document.querySelector("#messages");
        if (scrollHeight) {
          scrollHeight.scrollTop = scrollHeight.scrollHeight;
        }
      });

      console.log("Scrolling to bottom");
    },

    changeNumberOfUnreadNotification: async function (
      roomID: string,
      addToUnreadNotifs: boolean
    ) {
      this.unreadRoomMessages[roomID] = addToUnreadNotifs;

      let messageCount = 0;
      for (const room in this.unreadRoomMessages) {
        if (
          this.unreadRoomMessages[room] &&
          this.unreadRoomMessages[room] == true
        )
          messageCount++;
      }

      this.unreadRoomMessageCount = messageCount;
    },

    changeViewedRoomIndex: function (index: number) {
      this.indexOfCurrentViewedRoom = index;
      this.changeNumberOfUnreadNotification(
        this.joinedRooms[index].roomID,
        false
      );
    },

    joinRoom: function (
      requestingUserID: string,
      roomID: string,
      roomName: string,
      joined: boolean,
      index: number
    ) {
      const message = {
        msgType: WSMessageType.JoinRoom,
        userID: this.userID,
        requesterID: requestingUserID,
        joined: joined,
        roomID: roomID,
        roomName: roomName,
      };

      socket.send(JSON.stringify(message));
      this.joinRequests.splice(index, 1);
      this.unreadNotificationsCount = this.joinRequests.length;
    },

    createRoom: function () {
      const message = {
        userID: this.userID,
        msgType: WSMessageType.CreateRoom,
        roomName: this.newRoomName,
      };

      socket.send(JSON.stringify(message));
      this.newRoomName = "";
    },

    clearFetchedUsers: function () {
      this.fetchedUsers = [];
    },

    sendWSMessage: function (message: string) {
      socket.send(message);
    },

    onUploadFileChunks: function (file: FileUploadDownloadDetails) {
      if (
        this.fileUploadDownload[file.fileHash] &&
        !this.fileUploadDownload[file.fileHash].isDownloader
      ) {
        const fileDetails = this.fileUploadDownload[file.fileHash];
        fileDetails.chunk = file.nextChunk;

        // If file is successfully downloaded, send a FileUploadSuccess message to server
        // so as to be broadcasted to other users.
        if (fileDetails.chunk == fileDetails.chunks) {
          fileDetails.downloading = false;
          fileDetails.progress = 100;
          console.log("success");

          const message = {
            msgType: WSMessageType.FileUploadSuccess,
            userID: this.userID,
            name: this.userID,
            fileName: fileDetails.fileName,
            roomID: fileDetails.roomID,
            fileSize: fileDetails.fileSize.toString() + "MB",
            fileHash: fileDetails.fileHash,
          };
          socket.send(JSON.stringify(message));

          const infoMessage: Message = {
            time: "",
            size: "",
            hash: "",
            type: MessageType.Info,
            name: "",
            userID: "",
            roomID: fileDetails.roomID,
            message: "File successfully uploaded.",
            index: 0,
          };

          delete this.fileUploadDownload[fileDetails.fileHash];

          if (this.currentViewedRoom.roomID === fileDetails.roomID)
            this.currentViewedRoom.messages.push(infoMessage);

          return;
        }

        fileDetails.progress = (fileDetails.chunk / fileDetails.chunks) * 100;
        this.updateChatRoomPage();

        if (fileDetails.downloading === false) {
          return;
        }

        const offset = fileDetails.chunk * DefaultChunkSize;

        const reader = new FileReader();
        reader.onload = (event: ProgressEvent<FileReader>) => {
          if (!event.target) {
            return;
          }

          const data = event.target.result;
          if (!data) {
            return;
          }

          const fileUniqueHash = CryptoJS.SHA256(data.toString()).toString();

          const message = {
            msgType: WSMessageType.UploadFileChunk,
            userID: this.userID,
            file: data,
            fileName: fileDetails.fileName,
            fileHash: fileDetails.fileHash,
            newChunkHash: fileUniqueHash,
            recentChunkHash: "",
            chunkIndex: fileDetails.chunk,
          };

          console.log("sending", fileDetails.chunk, "now");

          // Send new file upload information to server, chunk == 0 indicates a new file.
          if (fileDetails.chunk === 0) {
            socket.send(JSON.stringify(message));
            return;
          }

          const recentOffset = (fileDetails.chunk - 1) * DefaultChunkSize;

          // Get recent file chunk hash.
          const recentFileChunkReader = new FileReader();
          recentFileChunkReader.onload = (e: ProgressEvent<FileReader>) => {
            if (!e.target) {
              return;
            }

            const recentData = e.target.result;
            if (!recentData) {
              return;
            }

            const recentFileUniqueHash = CryptoJS.SHA256(
              recentData.toString()
            ).toString();

            message.recentChunkHash = recentFileUniqueHash;
            socket.send(JSON.stringify(message));
          };

          if (fileDetails.fileContent) {
            const fileContent = fileDetails.fileContent;

            recentFileChunkReader.readAsDataURL(
              fileContent.slice(recentOffset, recentOffset + DefaultChunkSize)
            );
          }
        };

        if (fileDetails.fileContent) {
          const fileContent = fileDetails.fileContent;

          reader.readAsDataURL(
            fileContent.slice(offset, offset + DefaultChunkSize)
          );
        }
      }
    },

    onDownloadFileChunk: function (fileDetails: FileUploadDownloadDetails) {
      if (
        !fileDetails.compressedFileHash ||
        !this.fileUploadDownload[fileDetails.compressedFileHash]
      ) {
        console.log(
          "Server sent an invalid hash on file download",
          fileDetails
        );
        this.dialogError = "Server sent an invalid file content.";
        this.dialog = true;
        return;
      }

      if (!fileDetails.fileChunk) {
        console.log("Server sent a null file chunk.");
        return;
      }

      const file = this.fileUploadDownload[fileDetails.compressedFileHash];

      if (
        fileDetails.fileChunk &&
        fileDetails.fileHash !==
          CryptoJS.SHA256(fileDetails.fileChunk).toString()
      ) {
        this.dialogError = "Incorrect file hash.";
        this.dialog = true;

        file.downloading = false;
        this.updateChatRoomPage();

        return;
      }

      fileDetails.fileChunk = fileDetails.fileChunk.substr(
        fileDetails.fileChunk.indexOf(",") + 1
      );

      const byteCharacters = atob(fileDetails.fileChunk);
      const byteArrays = [];

      for (
        let offset = 0;
        offset < byteCharacters.length;
        offset += DefaultChunkSize
      ) {
        const slice = byteCharacters.slice(offset, offset + DefaultChunkSize);
        const byteNumbers = new Array(slice.length);

        for (let i = 0; i < slice.length; i++) {
          byteNumbers[i] = slice.charCodeAt(i);
        }

        const byteArray = new Uint8Array(byteNumbers);
        byteArrays.push(byteArray);
      }

      const blob = new Blob(byteArrays, { type: "" });
      file.fileContent.push(blob);

      if (file.chunk == file.chunks) {
        this.downloadBlob(fileDetails.compressedFileHash);

        delete this.fileUploadDownload[fileDetails.compressedFileHash];
        this.updateChatRoomPage();
        return;
      }

      if (file.downloading === false) {
        console.log("file not downloading");
        return;
      }

      file.chunk++;
      file.progress = (file.chunk / file.chunks) * 100;

      this.updateChatRoomPage();
      const message = {
        userID: this.userID,
        msgType: WSMessageType.DownloadFileChunk,
        fileName: file.fileName,
        compressedFileHash: fileDetails.compressedFileHash,
        chunkIndex: file.chunk,
      };

      socket.send(JSON.stringify(message));
    },

    onRequestDownload: function (fileDetails: FileUploadDownloadDetails) {
      if (!fileDetails.fileHash) {
        this.dialogError = "Server sent an invalid file download detail.";
        this.dialog = true;

        console.log("File hash empty on request download", fileDetails);
        return;
      }

      const file = this.fileUploadDownload[fileDetails.fileHash];
      file.downloading = true;
      file.chunks = fileDetails.chunks;
      this.updateChatRoomPage();

      const message = {
        userID: this.userID,
        msgType: WSMessageType.DownloadFileChunk,
        fileName: file.fileName,
        // File hash of compressed file.
        compressedFileHash: file.fileHash,
        chunkIndex: file.chunk,
      };
      socket.send(JSON.stringify(message));
    },

    downloadBlob: function (fileHash: string) {
      const file = this.fileUploadDownload[fileHash];
      file.downloading = false;

      const blobB = new Blob(file.fileContent, {
        type: file.fileType,
      });

      const a = document.createElement("a");
      a.href = URL.createObjectURL(blobB);
      a.download = file.fileName;
      a.target = "_blank";
      a.click();
    },

    changeDownloadStatus: function (fileHash: string, isDownloading: boolean) {
      this.fileUploadDownload[fileHash].downloading = isDownloading;
      this.updateChatRoomPage();
    },

    broadcastFileToRoom: function (completeFileDetails: FileUploadComplete) {
      const message: Message = {
        time: "",
        size: completeFileDetails.fileSize,
        hash: completeFileDetails.fileHash,
        type: MessageType.File,
        name: completeFileDetails.name,
        userID: completeFileDetails.userID,
        roomID: completeFileDetails.roomID,
        message: completeFileDetails.fileName,
        index: this.currentViewedRoom.messages.length + 1,
      };

      if (this.currentViewedRoom.roomID === completeFileDetails.roomID) {
        this.currentViewedRoom.messages.push(message);
      }
    },

    initiateFile: function (file: FileUploadDownloadDetails) {
      this.fileUploadDownload[file.fileHash] = file;
    },

    startCallSession: function () {
      this.callUI = true;
      this.isPublisher = true;
      this.classSessionConnectedUsers = {};
      this.classSessionConnectedUsers[this.userID] = {
        name: this.userID,
        isOnline: true,
      };

      peerConnection = new RTCPeerConnection({
        iceServers: [
          {
            urls: ["stun:stun.l.google.com:19302"],
          },
        ],
      });

      peerConnection.oniceconnectionstatechange = () => {
        console.log("ICE connection state", peerConnection.iceConnectionState);

        if (peerConnection.iceConnectionState === "failed") {
          // Show error page so that users can choose to disconnect call.
          this.dialogError = "Peer connection failed.";
          this.dialog = true;
        }
      };

      peerConnection.onicecandidate = (event) => {
        console.log("New candidate added", event);

        if (
          event.candidate === null &&
          peerConnection.localDescription &&
          peerConnection.localDescription.sdp &&
          peerConnection.iceConnectionState === "new"
        ) {
          const message = {
            msgType: WSMessageType.StartClassSession,
            sdp: peerConnection.localDescription.sdp,
            roomID: this.currentViewedRoom.roomID,
            userID: this.userID,
          };

          socket.send(JSON.stringify(message));
        }
      };

      peerConnection.onnegotiationneeded = (event) => {
        console.log("Negotiation needed", event);
      };

      const MediaConstraints = {
        video: { width: 640, height: 480 },
        audio: true,
      };

      this.getUserMedia(MediaConstraints, (e: MediaStream) => {
        videoTrack = e.getVideoTracks()[0];
        audioTrack = e.getAudioTracks()[0];

        audioTransceiver = peerConnection.addTransceiver(audioTrack, {
          direction: "sendrecv",
        });

        videoTransceiver = peerConnection.addTransceiver(videoTrack, {
          direction: "sendrecv",
        });

        const el = document.getElementById("videoID") as HTMLVideoElement;

        if (!el) {
          this.dialogError = "Error getting video element.";
          this.dialog = true;
          console.log("Could not get video ID element");
          return;
        }

        const mediaStreamTrack = [videoTrack];
        stream = new MediaStream(mediaStreamTrack);

        el.srcObject = stream;
        el.autoplay = true;

        peerConnection
          .createOffer()
          .then((sdp) => {
            peerConnection.setLocalDescription(sdp);
          })
          .catch((log) =>
            console.log("Error creating peer connection offer, error:", log)
          );
      });

      peerConnection.ontrack = ({ transceiver, streams: [event] }) => {
        event.onaddtrack = (event) => {
          console.log("On add track called for start video session.");
          stream.addTrack(event.track);
          console.log(event.track.id);
        };

        event.onremovetrack = (event) => {
          console.log("On remove track called for start video session.");
          stream.removeTrack(event.track);
          console.log(event.track.id);
        };

        transceiver.receiver.track.onmute = () => {
          console.log("Track muted for start session.");
        };

        transceiver.receiver.track.onended = () => {
          console.log("Track ended for start session");

          const id = event.id;

          this.classSessionConnectedUsers[id] = {
            name: id,
            isOnline: false,
          };

          this.updateCallPage();
        };

        transceiver.receiver.track.onunmute = () => {
          console.log("Track started for start session.");
          console.log("track", transceiver, "\n event", event);

          const id = event.id;
          this.classSessionConnectedUsers[id] = {
            name: id,
            isOnline: true,
          };

          stream.addTrack(transceiver.receiver.track);
        };
      };
    },

    joinCallSession: function (sessionData: Message) {
      this.callUI = true;
      this.isPublisher = false;
      this.classSessionConnectedUsers = {};

      peerConnection = new RTCPeerConnection({
        iceServers: [
          {
            urls: ["stun:stun.l.google.com:19302"],
          },
        ],
      });

      peerConnection.oniceconnectionstatechange = () => {
        console.log(peerConnection.iceConnectionState);

        if (peerConnection.iceConnectionState === "failed") {
          this.dialogError = "Connection state failed";
          this.dialog = true;
        }
      };

      peerConnection.onnegotiationneeded = (event) => {
        console.log("Negotiation needed", event);
      };

      peerConnection.onicecandidate = (event) => {
        console.log("A new user added", event);

        if (
          event.candidate == null &&
          peerConnection.localDescription &&
          peerConnection.localDescription.sdp &&
          peerConnection.iceConnectionState === "new"
        ) {
          console.log("Calling join class session");

          const message = {
            msgType: WSMessageType.JoinClassSession,
            sdp: peerConnection.localDescription.sdp,
            roomID: this.currentViewedRoom.roomID,
            userID: this.userID,
            sessionID: sessionData.hash,
            publisherID: sessionData.userID,
          };
          socket.send(JSON.stringify(message));
        }
      };

      const MediaConstraints = {
        audio: true,
      };

      this.getUserMedia(MediaConstraints, (e: MediaStream) => {
        console.log("Adding audio track");
        audioTrack = e.getAudioTracks()[0];
        audioTrack.enabled = true;
        audioTransceiver = peerConnection.addTransceiver(audioTrack, {
          direction: "sendrecv",
        });

        const el = document.getElementById("videoID") as HTMLVideoElement;

        if (!el) {
          this.dialogError = "Error getting video element.";
          this.dialog = true;

          console.log("Could not get video ID element");
          return;
        }

        stream = new MediaStream();

        el.srcObject = stream;
        el.autoplay = true;

        console.log("Sending to server to create offer.");
        peerConnection
          .createOffer()
          .then((sdp) => {
            peerConnection.setLocalDescription(sdp);
          })
          .catch((log) => {
            console.log("Error creating offer on join session, error:", log);
            this.dialogError = "Error creating peer offer";
            this.dialog = true;
          });
      });

      peerConnection.ontrack = ({ transceiver, streams: [event] }) => {
        event.onaddtrack = (event) => {
          console.log("On add track called for join video session");
          stream.addTrack(event.track);
          console.log("track ID", event.track.id);
        };

        event.onremovetrack = (track) => {
          console.log("On remove track called for join video session.");
          stream.removeTrack(track.track);
          console.log(track.track.id);

          const id = track.track.id;
          delete this.classSessionConnectedUsers[id];

          this.updateCallPage();
        };

        transceiver.receiver.track.onmute = () =>
          console.log("Track muted for join session");

        transceiver.receiver.track.onended = () => {
          console.log("Track ended for Join session");
        };

        transceiver.receiver.track.onunmute = () => {
          if (transceiver.receiver.track.kind === "audio") {
            const id = transceiver.receiver.track.id;
            this.classSessionConnectedUsers[id] = {
              name: event.id,
              isOnline: true,
            };
          }

          console.log("Track unmuted");

          console.log("Track started for Join session");

          console.log("track", transceiver, "\n event", event);
          stream.addTrack(transceiver.receiver.track);
        };
      };
    },

    endCallSession: function () {
      this.callUI = false;

      if (audioTrack) {
        audioTrack.enabled = false;
        audioTrack.stop();
        audioTransceiver.stop();
      }

      if (videoTrack) {
        videoTrack.enabled = false;
        videoTrack.stop();
        videoTransceiver.stop();
      }

      if (peerConnection) peerConnection.close();
      // Send websocket close message to server.
      const message = {
        msgType: WSMessageType.EndClassSession,
        userID: this.userID,
      };

      this.isPublisher = false;
      console.log("sending websocket peer close message");

      socket.send(JSON.stringify(message));
    },

    getUserMedia: function (
      mediaConstraints: MediaStreamConstraints,
      onMedia: Function
    ) {
      navigator.mediaDevices
        .getUserMedia(mediaConstraints)
        .then(function (e) {
          onMedia(e);
        })
        .catch((error) => {
          console.log("Error getting media device, error:", error);
          this.dialogError = "Error starting media devices";
          this.dialog = true;
        });
    },

    changeVideo: function () {
      videoTrack.enabled = !videoTrack.enabled;

      videoTrack.enabled
        ? videoTransceiver.sender.replaceTrack(videoTrack)
        : videoTransceiver.sender.replaceTrack(null);
    },

    changeAudio: function () {
      audioTrack.enabled = !audioTrack.enabled;
      audioTrack.enabled
        ? audioTransceiver.sender.replaceTrack(audioTrack)
        : audioTransceiver.sender.replaceTrack(null);
    },

    startDesktopSharing: function () {
      // 480p video constraints.
      const mediaConstraints = { video: { width: 640, height: 480 } };

      navigator.mediaDevices
        .getDisplayMedia(mediaConstraints)
        .then((event) => {
          // Disable webcam and remove video from stream.
          videoTrack.enabled = false;
          stream.removeTrack(videoTrack);

          const tracks = event.getVideoTracks();
          this.addTracks(tracks);

          if (tracks.length > 0) {
            videoTransceiver.sender.replaceTrack(tracks[0]);
          } else {
            console.log("Could not find any track for desktop sharing");

            this.dialogError = "Could not get track for desktop sharing.";
            this.dialog = true;
          }
        })
        .catch((error) => {
          console.log("Error creating desktop sharing", error);
          this.dialogError = "Error starting desktop sharing.";
          this.dialog = true;
        });
    },

    stopDesktopSharing: function (isVideoOn: boolean) {
      const tracks = stream.getVideoTracks();
      this.removeTracks(tracks);

      videoTrack.enabled = isVideoOn;
      stream.addTrack(videoTrack);

      videoTransceiver.sender.replaceTrack(videoTrack);
    },

    removeTracks: function (tracks: MediaStreamTrack[]) {
      for (let i = 0; i < tracks.length; i++) {
        tracks[i].enabled = false;
        stream.removeTrack(tracks[i]);
      }
    },

    addTracks: function (tracks: MediaStreamTrack[]) {
      for (let i = 0; i < tracks.length; i++) {
        stream.addTrack(tracks[i]);
      }
    },

    connectWebsocket: function () {
      socket = new WebSocket("wss://" + document.location.host + "/ws");

      socket.onmessage = (event: MessageEvent) => {
        const jsonContent = JSON.parse(event.data);

        switch (jsonContent.msgType) {
          case WSMessageType.UnauthorizedAccess:
            this.onUnAuthorizedAccess();
            break;

          case WSMessageType.SearchUser:
            this.fetchedUsers = jsonContent.fetchedUsers;
            console.log(this.fetchedUsers);
            break;

          case WSMessageType.WebsocketOpen:
            this.onWebsocketOpen(
              jsonContent.joinedRooms,
              jsonContent.joinRequests
            );
            this.usersOnline = jsonContent.associates;
            this.userID = jsonContent.userID;
            this.names = jsonContent.names;
            break;

          case WSMessageType.RequestMessages:
            if (!this.showChatPage) {
              this.showChatPage = true;
            }

            this.showAvailableRoomsPage = false;

            this.onRequestMessages(jsonContent.roomPageDetails);
            break;

          case WSMessageType.JoinRoom:
            this.onJoinRoom(jsonContent);
            break;

          case WSMessageType.NewMessage:
            this.onNewMessage(jsonContent);
            break;

          case WSMessageType.SentRoomRequest:
            this.onSentRoomRequest(jsonContent);
            break;

          case WSMessageType.OnlineStatus:
            if (this.usersOnline[jsonContent.userID]) {
              this.usersOnline[jsonContent.userID].isOnline =
                jsonContent.status;
            }
            break;

          case WSMessageType.ExitRoom:
            this.onExitRoom(jsonContent);
            break;

          case WSMessageType.UploadFileChunk:
            this.onUploadFileChunks(jsonContent);
            break;

          case WSMessageType.FileUploadSuccess:
            this.broadcastFileToRoom(jsonContent);
            break;

          case WSMessageType.UploadFileError:
            if (this.$children.length > 2) {
              this.$children[2].$emit("onUploadError");
            }
            this.dialogError = "File upload error.";
            this.dialog = true;
            break;

          case WSMessageType.FileRequestDownload:
            this.onRequestDownload(jsonContent);
            break;

          case WSMessageType.DownloadFileChunk:
            console.log("Downloading");
            this.onDownloadFileChunk(jsonContent);
            break;

          case WSMessageType.ClassSessionError:
            this.dialogError = jsonContent.errorDetails;
            this.dialog = true;
            break;

          case WSMessageType.ClassSession:
            if (this.currentViewedRoom.roomID === jsonContent.roomID) {
              this.currentViewedRoom.messages.push({
                type: MessageType.ClassSession,
                name: jsonContent.name,
                userID: jsonContent.userID,
                hash: jsonContent.sessionID,
                index: this.currentViewedRoom.messages.length + 1,
                roomID: jsonContent.roomID,
                message: "",
              });
              this.scrollToBottomOfChatPage();
            }
            break;

          case WSMessageType.NegotiateSDP:
            if (peerConnection) {
              console.log("Negotiating peer connection");

              try {
                peerConnection.setRemoteDescription(
                  new RTCSessionDescription({
                    type: "answer",
                    sdp: jsonContent.sdp,
                  })
                );
              } catch (e) {
                console.log(
                  "Error setting remote description on negotiating",
                  e
                );
                this.dialogError =
                  "Error setting peer remote description on negotiating.";
                this.dialog = true;
              }
            }
            break;

          case WSMessageType.RenegotiateSDP:
            console.log("Renegotiating SDP");

            if (peerConnection && jsonContent.sessionID !== "") {
              peerConnection.setRemoteDescription(
                new RTCSessionDescription({
                  type: "offer",
                  sdp: jsonContent.sdp,
                })
              );

              peerConnection
                .createAnswer()
                .then((sdp) => {
                  peerConnection.setLocalDescription(sdp).then(() => {
                    const message = {
                      msgType: WSMessageType.RenegotiateSDP,
                      sdp: sdp.sdp,
                      userID: this.userID,
                    };

                    socket.send(JSON.stringify(message));
                  });
                })
                .catch((e) => {
                  console.log("Error renegotiating peerconnection, error: ", e);
                  this.dialogError = "Error renegotiating peer connection.";
                  this.dialog = true;
                });
            }
            break;
        }
      };

      socket.onerror = (event: Event) => {
        // Reconnect Websocket if not UnAuthorized.
        console.log("Websocket errored.", event);

        this.dialogError = "Websocket connection error.";
        this.dialog = true;
        this.socketClosed = true;
      };

      socket.onopen = () => {
        console.log("Websocket open.");
        this.socketClosed = false;
        this.dialog = false;

        // Fetch users content from API.
        // Contents that are to be fetched from API are, Registered rooms and room request.
        const message = {
          msgType: WSMessageType.WebsocketOpen,
          userID: this.userID,
        };

        socket.send(JSON.stringify(message));
      };

      socket.onclose = (closeEvent: CloseEvent) => {
        console.log("Websocket closed.", closeEvent);
        this.socketClosed = true;

        this.dialog = true;
        this.dialogError = "Websocket disconnected. Reconnecting.";
      };
    },
  },

  mounted() {
    this.connectWebsocket();
  },
});
</script>


<style  scoped>
.transparent_sheet {
  opacity: 0.2;
  border-color: transparent !important;
  z-index: -1;
  position: absolute;
}
</style>
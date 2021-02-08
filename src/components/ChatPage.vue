<template>
  <v-row no-gutters>
    <v-col cols="12" style="height: 10vh">
      <v-app-bar tile flat height="80vh">
        <v-app-bar-nav-icon>
          <v-avatar>
            <v-img
              class="mx-2"
              :src="
                currentViewedRoom.roomIcon
                  ? currentViewedRoom.roomIcon
                  : require('../assets/unilag.svg')
              "
            ></v-img>
          </v-avatar>
        </v-app-bar-nav-icon>

        <v-toolbar-title
          style="width: 28vw"
          class="d-inline-block text--secondary"
        >
          <b>{{ currentViewedRoom.roomName }}</b>
        </v-toolbar-title>

        <v-spacer></v-spacer>

        <v-menu
          v-model="showAddUsersDialog"
          left
          rounded
          :close-on-content-click="false"
          offset-y
        >
          <template v-slot:activator="{ on, attrs }">
            <v-btn
              v-bind="attrs"
              v-on="on"
              icon
              :x-large="!$vuetify.breakpoint.smAndDown"
            >
              <v-icon>mdi-account-multiple-plus-outline</v-icon>
            </v-btn>
          </template>

          <v-card align="center">
            <v-row style="max-width: 40vh" align="center" justify="center">
              <v-col cols="12"></v-col>
              <v-col cols="12" sm="8">
                <v-text-field
                  v-model="searchUserName"
                  append-icon="mdi-magnify"
                  @keyup.enter.exact="searchUsers"
                  @click:append="searchUsers"
                  rounded
                  filled
                  placeholder="Search for contact email"
                />
              </v-col>

              <v-col cols="12">
                <v-container
                  fluid
                  class="ml-auto overflow-y-auto scroll-behavior-smooth"
                  style="max-height: 40vh"
                >
                  <v-checkbox
                    class="text--truncate"
                    dense
                    multiple
                    v-for="(user, i) in fetchedUsers"
                    :key="i"
                    :label="user.name + ' [' + user.userID + ']'"
                    :value="user.userID"
                    v-model="selectedUsersToAddToRoom"
                  ></v-checkbox>
                </v-container>
              </v-col>
            </v-row>

            <v-card-actions>
              <v-btn @click="closeAddUserDialog" color="red" text>exit</v-btn>
              <v-btn
                @click="requestUserToJoinRoom"
                :disabled="selectedUsersToAddToRoom.length === 0"
                color="green"
                text
                >add users</v-btn
              >
            </v-card-actions>
          </v-card>
        </v-menu>

        <v-btn
          @click="startCallSession"
          icon
          :x-large="!$vuetify.breakpoint.smAndDown"
        >
          <v-icon>mdi-video</v-icon>
        </v-btn>

        <v-menu
          v-model="showOnlineUsersDialog"
          left
          :close-on-content-click="false"
          nudge-width="300"
          offset-y
        >
          <template v-slot:activator="{ on, attrs }">
            <v-btn
              v-bind="attrs"
              v-on="on"
              icon
              :x-large="!$vuetify.breakpoint.smAndDown"
            >
              <v-icon>mdi-information</v-icon>
            </v-btn>
          </template>

          <v-card max-width="50wh">
            <v-row class="mx-10">
              <v-col cols="12"></v-col>
              <v-col
                cols="12"
                v-for="(onlineUser, index) in currentViewedRoom.registeredUsers"
                :key="index"
              >
                <template v-if="onlineUser == userID">
                  <v-badge inline dot color="green"></v-badge>
                  <span class="mx-4">
                    <b>You</b>
                  </span>
                </template>

                <template v-else>
                  <v-badge
                    inline
                    dot
                    :color="
                      associates[onlineUser] &&
                      associates[onlineUser].isOnline === true
                        ? 'green'
                        : 'red'
                    "
                  ></v-badge>
                  <span class="mx-4">
                    <b class="d-inline-block text--secondary"
                      >{{
                        associates[onlineUser]
                          ? associates[onlineUser].name
                          : ""
                      }}
                      {{ onlineUser }}</b
                    >
                  </span>
                </template>
              </v-col>
            </v-row>

            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn @click="exitRoom" color="red" text>exit room</v-btn>
              <v-spacer></v-spacer>
            </v-card-actions>
          </v-card>
        </v-menu>
      </v-app-bar>
    </v-col>

    <v-col cols="12" style="height: 78vh">
      <v-container
        id="messages"
        class="overflow-y-auto scroll-behavior-smooth"
        style="height: 75vh"
        fluid
      >
        <v-row v-scroll:#messages="onScroll" dense>
          <v-col
            class="my-2"
            cols="12"
            v-for="(message, index) in currentViewedRoom.messages"
            :key="index"
          >
            <template v-if="message.type === messageType.Message">
              <v-row
                dense
                align="end"
                justify="end"
                v-if="message.userID === userID"
              >
                <v-card
                  flat
                  :color="vuetify.framework.theme.dark ? '#778899' : '#bbeaff'"
                  max-width="65%"
                >
                  <v-card-text class="text--wrap">
                    <b>{{ message.message }}</b>
                  </v-card-text>

                  <v-card-subtitle align="right"
                    >{{ message.userID }} {{ message.time }}</v-card-subtitle
                  >
                </v-card>
              </v-row>

              <v-row dense v-else>
                <v-avatar class="mx-2" style="align-self: flex-end" rounded>
                  <v-img
                    height="50px"
                    width="50px"
                    contain
                    :src="require('../assets/unilag.svg')"
                  ></v-img>
                </v-avatar>
                <v-card
                  flat
                  :color="vuetify.framework.theme.dark ? '' : '#DCDCDC'"
                  max-width="65%"
                >
                  <v-card-text>
                    <b>{{ message.message }}</b>
                  </v-card-text>
                  <v-card-subtitle align="right" class="ml-auto">
                    <b>{{ message.userID }} {{ message.time }}</b>
                  </v-card-subtitle>
                </v-card>
              </v-row>
            </template>

            <template v-else>
              <v-col v-if="message.type === messageType.File" cols="12">
                <div align="center">
                  <v-chip
                    @click="
                      downloadFile(
                        message.hash,
                        message.message,
                        message.size,
                        message.type
                      )
                    "
                  >
                    <b
                      style="max-width: 60vw"
                      class="d-inline-block text--secondary"
                      >{{ message.message }} sent by {{ message.name }} [{{
                        message.userID
                      }}]. Click to download.</b
                    >
                  </v-chip>
                </div>
              </v-col>

              <v-col v-if="message.type === messageType.Info" cols="12">
                <div align="center">
                  <v-chip>
                    <b
                      style="max-width: 60vw"
                      class="d-inline-block text--secondary"
                      >{{ message.message }}.</b
                    >
                  </v-chip>
                </div>
              </v-col>

              <v-col v-if="message.type === messageType.ClassSession" cols="12">
                <div align="center">
                  <v-chip @click="joinCallSession(message)">
                    <b
                      style="max-width: 60vw"
                      class="d-inline-block text--secondary"
                      >Class session started by {{ message.name }}. Click to
                      join.</b
                    >
                  </v-chip>
                </div>
              </v-col>

              <v-col
                v-if="message.type === messageType.ClassSessionLink"
                cols="12"
              >
                <div align="center">
                  <v-chip @click="downloadSession(chat.message)">
                    <b
                      style="max-width: 60vw"
                      class="d-inline-block text--secondary"
                      >Class session recording by {{ message.name }}. Click to
                      download.</b
                    >
                  </v-chip>
                </div>
              </v-col>
            </template>
          </v-col>

          <v-col
            cols="12"
            v-if="Object.keys(this.fileUploadDownload).length > 0"
          >
            <v-row v-for="(file, index) in fileUploadDownload" :key="index">
              <v-col
                v-if="file.roomID == currentViewedRoom.roomID"
                :align="file.isDownloader ? 'right' : 'left'"
              >
                <v-card align="left" shaped max-width="40%">
                  <v-card-subtitle>
                    <h3>
                      {{ file.userID }}
                    </h3>
                  </v-card-subtitle>

                  <v-card-text>
                    <v-row align="center">
                      <v-col cols="mx-auto"
                        >{{ file.fileName }} ({{ file.fileSize }} MB)</v-col
                      >
                      <v-col cols="auto">
                        <v-progress-circular
                          :rotate="360"
                          :size="50"
                          :width="5"
                          :value="file.progress"
                          color="teal"
                        >
                          <v-btn
                            icon
                            @click="
                              file.downloading
                                ? stopFileProgress(file.fileHash)
                                : startFileProgress(file.fileHash)
                            "
                          >
                            <v-icon v-if="file.downloading">mdi-close</v-icon>
                            <v-icon v-else-if="file.isDownloader"
                              >mdi-cloud-download</v-icon
                            >
                            <v-icon v-else>mdi-cloud-upload</v-icon>
                          </v-btn>
                        </v-progress-circular>
                      </v-col>
                    </v-row>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>
          </v-col>
        </v-row>
      </v-container>
    </v-col>

    <v-col cols="12" style="height: 10vh">
      <v-expand-transition>
        <v-form v-model="fileInputValid" v-show="showFileInput">
          <v-file-input
            :rules="rules"
            v-show="showFileInput"
            class="mx-5"
            v-model="file"
            counter
            label="File input"
            placeholder="Select your file"
            outlined
            append-outer-icon="mdi-send"
            :show-size="1000"
            @click:append-outer="sendFile"
          ></v-file-input>
        </v-form>
      </v-expand-transition>

      <v-expand-transition>
        <v-textarea
          class="mx-5"
          v-show="showTextField"
          v-model="messageContent"
          prepend-inner-icon="mdi-emoticon"
          prepend-icon="mdi-paperclip"
          append-outer-icon="mdi-send"
          solo
          hide-details="auto"
          clearable
          no-resize
          rows="3"
          rounded
          @click:prepend="hideTextField"
          @click:append-outer="sendMessage"
          @keyup.enter.exact="sendMessage"
        ></v-textarea>
      </v-expand-transition>
    </v-col>
  </v-row>
</template>

<script lang="ts">
"use strict";

import Vue from "vue";
import vuetify from "@/plugins/vuetify";
import { Prop } from "vue/types/options";
import CryptoJS from "crypto-js";

import {
  RoomPageDetails,
  FetchedUsers,
  UsersOnline,
  FileDownload,
  FileUploadDownloadDetails,
} from "../views/Types";

import {
  WSMessageType,
  MessageType,
  DefaultChunkSize,
} from "../views/Constants";

export default Vue.extend({
  name: "ChatPage",

  props: {
    currentViewedRoom: {} as Prop<RoomPageDetails>,
    fetchedUsers: Array as Prop<FetchedUsers[]>,
    associates: {} as Prop<UsersOnline>,
    fileUploadDownload: {} as Prop<FileDownload>,

    sendWSMessage: Function,
    clearFetchedUsers: Function,
    initiateFile: Function,
    changeDownloadStatus: Function,
    startCallSession: Function,
    joinCallSession: Function,
    userID: String,
    names: String,
  },

  data: () => ({
    vuetify: vuetify,
    messageContent: "",
    newRoomName: "",
    searchUserName: "",

    showTextField: true,
    showFileInput: false,
    showAddUsersDialog: false,
    showOnlineUsersDialog: false,
    fileInputValid: false,
    isFile: false, // Indicate file download or upload.

    selectedUsersToAddToRoom: [] as string[],
    file: {} as File,

    messageType: MessageType,

    rules: [
      (value: File) => !!value || "File cannot be empty",
      (value: File) =>
        !value || value.size < 256 * 1024 * 1024 || "File size limit is 256MB!",
    ],
  }),

  methods: {
    onScroll: function (e: { target: Element }) {
      const messages = this.currentViewedRoom.messages;
      if (
        e.target &&
        e.target.scrollTop < 100 &&
        messages.length > 0 &&
        messages[0].index > 1
      ) {
        this.loadMoreMessages();
      }
    },

    sendMessage: function () {
      if (!this.messageContent.match(/\S/)) return;

      const message = {
        msgType: WSMessageType.NewMessage,
        userID: this.userID,
        roomID: this.currentViewedRoom.roomID,
        message: this.messageContent,
        type: this.messageType.Message,
      };

      this.sendWSMessage(JSON.stringify(message));
      this.messageContent = "";
    },

    sendFile: function () {
      if (this.fileInputValid == false) return;

      const reader = new FileReader();
      reader.onloadend = (event: ProgressEvent<FileReader>) => {
        if (event.target) {
          const result = event.target.result;
          if (result) {
            const uniqueFileHash = CryptoJS.SHA256(
              result.toString()
            ).toString();

            const chunks = Math.ceil(this.file.size / DefaultChunkSize);

            const file: FileUploadDownloadDetails = {
              roomID: this.currentViewedRoom.roomID,
              userID: this.userID,
              fileName: this.file.name,
              fileHash: uniqueFileHash,
              downloading: true,
              isDownloader: false,
              fileSize: (this.file.size / (1024 * 1024)).toFixed(),
              progress: 0,
              chunks: chunks,
              chunk: 0,
              nextChunk: 0,
              fileType: this.file.type,
              fileContent: this.file,
            };

            this.initiateFile(file);

            this.showTextField = true;
            this.showFileInput = false;
            this.file = new File([], "");

            setTimeout(() => {
              this.$forceUpdate();
              this.scrollToBottomOfChatPage();
            }, 0);

            this.requestFileUpload(uniqueFileHash);
          }
        }
      };

      reader.readAsDataURL(this.file);
    },

    downloadFile: function (
      fileHash: string,
      fileName: string,
      fileSize: string,
      fileType?: string
    ) {
      const file: FileUploadDownloadDetails = {
        roomID: this.currentViewedRoom.roomID,
        userID: this.userID,
        fileName: fileName,
        fileHash: fileHash,
        downloading: true,
        fileContent: [],
        isDownloader: true,
        fileSize: fileSize,
        progress: 0,
        chunks: 0,
        chunk: 0,
        nextChunk: 0,
        fileType: fileType,
      };

      this.initiateFile(file);

      setTimeout(() => {
        this.$forceUpdate();
        this.scrollToBottomOfChatPage();
      }, 0);

      this.requestFileDownload(fileName, fileHash);
    },

    startFileProgress: function (fileHash: string) {
      const file = this.fileUploadDownload[fileHash];

      if (file.isDownloader) {
        this.requestFileDownload(file.fileName, file.fileHash);
      } else {
        this.requestFileUpload(fileHash);
      }

      this.changeDownloadStatus(fileHash, true);
    },

    stopFileProgress: function (fileHash: string) {
      this.changeDownloadStatus(fileHash, false);
    },

    requestFileUpload: function (fileHash: string) {
      const message = {
        msgType: WSMessageType.NewFileUpload,
        userID: this.userID,
        fileName: this.fileUploadDownload[fileHash].fileName,
        fileHash: this.fileUploadDownload[fileHash].fileHash,
        fileSize: this.fileUploadDownload[fileHash].fileSize,
        fileType: this.fileUploadDownload[fileHash].fileType,
      };
      this.sendWSMessage(JSON.stringify(message));
    },

    requestFileDownload: function (fileName: string, fileHash: string) {
      const message = {
        msgType: WSMessageType.FileRequestDownload,
        userID: this.userID,
        fileName: fileName,
        fileHash: fileHash,
      };
      this.sendWSMessage(JSON.stringify(message));
    },

    scrollToBottomOfChatPage: function () {
      setTimeout(() => {
        const scrollHeight = this.$el.querySelector("#messages");
        if (scrollHeight) {
          scrollHeight.scrollTop = scrollHeight.scrollHeight;
        }
      }, 0);
    },

    loadMoreMessages: function () {
      const message = {
        msgType: WSMessageType.RequestMessages,
        userID: this.userID,
        roomID: this.currentViewedRoom.roomID,

        messageCount: this.currentViewedRoom.messages[0].index - 1,
      };

      this.sendWSMessage(JSON.stringify(message));
    },

    searchUsers: function () {
      const message = {
        msgType: WSMessageType.SearchUser,
        userID: this.userID,
        searchText: this.searchUserName,
      };

      this.sendWSMessage(JSON.stringify(message));
    },

    requestUserToJoinRoom: function () {
      const message = {
        msgType: WSMessageType.RequestUsersToJoinRoom,
        userID: this.userID,
        roomID: this.currentViewedRoom.roomID,
        roomName: this.currentViewedRoom.roomName,
        requestingUserName: this.names,
        requestingUserID: this.userID,
        users: this.selectedUsersToAddToRoom,
      };

      this.sendWSMessage(JSON.stringify(message));
      this.closeAddUserDialog();
    },

    closeAddUserDialog: function () {
      this.selectedUsersToAddToRoom = [];
      this.showAddUsersDialog = false;

      this.clearFetchedUsers();
    },

    exitRoom: function () {
      const message = {
        msgType: WSMessageType.ExitRoom,
        roomID: this.currentViewedRoom.roomID,
        userID: this.userID,
      };

      this.sendWSMessage(JSON.stringify(message));
      this.showOnlineUsersDialog = false;
    },

    hideTextField: function () {
      if (this.fileUploadDownload.downloading) return;
      this.showFileInput = true;
      this.showTextField = false;
    },

    revealTextField: function () {
      this.showFileInput = false;
      this.showTextField = true;
    },

    downloadSession: function (link: string) {
      window.open(link, "_blank");
    },
  },

  watch: {
    file: function () {
      if (!this.file) {
        this.revealTextField();
      }
    },
  },
});
</script>

<template>
  <div>
    <v-card flat height="100vh" align="center" outlined width="max-content">
      <v-row class="fill-height">
        <v-col class="my-2" cols="12">
          <v-avatar>
            <img src="../assets/unilag.svg" alt="Unilag Logo" />
          </v-avatar>
        </v-col>

        <v-col cols="12">
          <v-row>
            <v-col class="my-2" cols="12">
              <v-badge
                bordered
                overlap
                :content="unreadRoomMessages"
                :value="unreadRoomMessages"
              >
                <v-btn
                  @click="onRoom"
                  :x-large="!$vuetify.breakpoint.smAndDown"
                  icon
                >
                  <v-icon
                    :x-large="!$vuetify.breakpoint.smAndDown"
                    :color="highlightedSidebarOption == 'Rooms' ? 'blue' : ''"
                    >mdi-message-outline</v-icon
                  >
                </v-btn>
              </v-badge>
            </v-col>

            <v-col class="my-2" cols="12">
              <v-btn
                :x-large="!$vuetify.breakpoint.smAndDown"
                @click="onAddRoom"
                icon
              >
                <v-icon
                  :color="highlightedSidebarOption == 'AddRoom' ? 'blue' : ''"
                  :x-large="!$vuetify.breakpoint.smAndDown"
                  >mdi-message-plus</v-icon
                >
              </v-btn>
            </v-col>

            <v-col class="my-2" cols="12">
              <v-badge
                bordered
                overlap
                :content="unreadNotifications"
                :value="unreadNotifications"
              >
                <v-btn
                  :x-large="!$vuetify.breakpoint.smAndDown"
                  @click="onNotification"
                  icon
                >
                  <v-icon
                    :color="
                      highlightedSidebarOption == 'Notifications' ? 'blue' : ''
                    "
                    :x-large="!$vuetify.breakpoint.smAndDown"
                    >mdi-bell-outline</v-icon
                  >
                </v-btn>
              </v-badge>
            </v-col>
          </v-row>
        </v-col>

        <v-col class="mt-auto" cols="12">
          <v-row>
            <v-col class="my-2" cols="12">
              <v-btn
                :x-large="!$vuetify.breakpoint.smAndDown"
                @click="changeTheme"
                icon
              >
                <v-icon
                  :color="themeColor"
                  :x-large="!$vuetify.breakpoint.smAndDown"
                  >mdi-moon-waxing-crescent</v-icon
                >
              </v-btn>
            </v-col>

            <v-col class="my-2" cols="12">
              <v-btn :x-large="!$vuetify.breakpoint.smAndDown" icon>
                <v-icon :x-large="!$vuetify.breakpoint.smAndDown"
                  >mdi-cog</v-icon
                >
              </v-btn>
            </v-col>

            <v-col class="my-2" cols="12">
              <v-btn
                :x-large="!$vuetify.breakpoint.smAndDown"
                @click="logOff"
                icon
              >
                <v-icon :x-large="!$vuetify.breakpoint.smAndDown"
                  >mdi-power-standby</v-icon
                >
              </v-btn>
            </v-col>
          </v-row>
        </v-col>
      </v-row>
    </v-card>
  </div>
</template>

<script lang="ts">
"use strict";

import Vue from "vue";
import vuetify from "@/plugins/vuetify";

import { SideBarOptions } from "../views/Constants";

export default Vue.extend({
  name: "SideBar",

  props: {
    activateAddRoomDialog: Function,
    activateNotificationDialog: Function,
    deactivateAllDialogs: Function,
    showAvailableRooms: Function,

    unreadRoomMessages: Number,
    unreadNotifications: Number,
  },

  data: () => ({
    themeColor: "",
    highlightedSidebarOption: "Rooms",
    newRoomName: "",

    showAddRoomDialog: false,
  }),

  methods: {
    changeTheme: function () {
      vuetify.framework.theme.dark = !vuetify.framework.theme.dark;
      if (vuetify.framework.theme.dark) {
        this.themeColor = "blue";
      } else {
        this.themeColor = "";
      }
    },

    onNotification: function () {
      this.highlightedSidebarOption = SideBarOptions.Notifications;
      this.activateNotificationDialog();
    },

    onRoom: function () {
      this.highlightedSidebarOption = SideBarOptions.Rooms;
      this.deactivateAllDialogs();
      this.showAvailableRooms();
    },

    onAddRoom: function () {
      this.highlightedSidebarOption = SideBarOptions.AddRoom;
      this.activateAddRoomDialog();
    },

    logOff: function () {
      // this.$store.commit("setToken", "");
      // this.$router.push("/login");
    },
  },
});
</script>

<style>
</style>
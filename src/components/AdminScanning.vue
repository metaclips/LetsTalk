<template>
  <v-container fluid>
    <v-row class="mx-6 my-4" no-gutters>
      <v-col class="mx-auto" cols="12" sm="3">
        <v-select
          rounded
          filled
          outlined
          :items="scanningOptions"
          v-model="scanningOption"
          label="Scanning Options"
        >
        </v-select>
      </v-col>
      <v-col class="mx-auto" align="center" cols="12"
        ><v-btn depressed dark color="green">Start</v-btn></v-col
      >
    </v-row>

    <v-row class="mt-6 ml-6">
      <template v-if="showUserDetails">
        <v-col cols="6">
          <v-row>
            <v-col cols="4">
              <p class="font-weight-black">NAME</p>
            </v-col>
            <v-col class="ml-auto" cols="8">
              <p class="font-weight-regular">
                {{ userName.toUpperCase() }}
              </p></v-col
            >

            <v-col cols="12"><v-divider></v-divider></v-col>

            <v-col cols="4">
              <p class="font-weight-black">LEVEL</p>
            </v-col>
            <v-col class="ml-auto" cols="8">
              <p class="font-weight-regular">{{ level }}</p></v-col
            >

            <v-col cols="12"><v-divider></v-divider></v-col>

            <v-col cols="4">
              <p class="font-weight-black">DEPARTMENT</p>
            </v-col>
            <v-col class="ml-auto" cols="8">
              <p class="font-weight-regular">
                {{ department.toUpperCase() }}
              </p></v-col
            >

            <v-col cols="12"><v-divider></v-divider></v-col>

            <v-col cols="4">
              <p class="font-weight-black">NUMBER OF ROOMS JOINED</p>
            </v-col>
            <v-col cols="8">
              <v-row no-gutters
                ><v-col cols="4">
                  <p class="font-weight-regular">{{ roomsJoined }}</p>
                </v-col>
                <v-col cols="8"
                  ><v-btn depressed dark color="red"> Scan Rooms</v-btn></v-col
                >
              </v-row>
            </v-col>

            <v-col cols="12"><v-divider></v-divider></v-col>

            <v-col cols="4">
              <p class="font-weight-black">LEVEL</p>
            </v-col>
            <v-col cols="8">
              <p class="font-weight-regular">
                {{ level }}
              </p></v-col
            >
            <v-col cols="12"><v-divider></v-divider></v-col>

            <v-col cols="4">
              <p class="font-weight-black">Banned</p>
            </v-col>
            <v-col cols="8">
              <p class="font-weight-regular">
                {{ isBanned ? "Yes" : "No" }}
              </p></v-col
            >
          </v-row>
        </v-col>
        <v-col cols="6"></v-col>
      </template>

      <template v-else>
        <v-col cols="12"> </v-col>
      </template>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import Vue from "vue";
export default Vue.extend({
  name: "AdminScanning",
  data: () => ({
    scanningOptions: ["Scan All Rooms", "Scan Single User"],
    scanningOption: "Scan All Rooms",

    userName: "My Names",
    level: 1,
    department: "Computer Engineering",
    roomsJoined: 0,
    isBanned: false,
  }),

  computed: {
    showUserDetails(): boolean {
      return this.scanningOption == "Scan Single User";
    },
  },
});
</script>
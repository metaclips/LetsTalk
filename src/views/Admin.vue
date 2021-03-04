<template>
  <v-container fluid>
    <v-row align="center" class="my-auto">
      <v-col
        style="display: grid; justify-content: center"
        cols="12"
        class="mx-auto"
      >
        <v-img width="200" height="200" src="../assets/unilag.svg"></v-img>
      </v-col>
      <v-col
        style="display: grid; justify-content: center"
        cols="12"
        class="mx-auto"
      >
        <h2>Admin Portal</h2>
      </v-col>
      <v-col cols="12" class="mx-auto">
        <v-tabs v-model="tab" fixed-tabs>
          <v-tab>Add User</v-tab>
          <v-tab>Update Users Details</v-tab>
          <v-tab>Room Messages Scanning</v-tab>
        </v-tabs>

        <v-tabs-items v-model="tab">
          <v-tab-item>
            <Register />
          </v-tab-item>

          <v-tab-item>
            <UpdateUser />
          </v-tab-item>

          <v-tab-item>
            <AdminScanner />
          </v-tab-item>
        </v-tabs-items>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
"use strict";

import Vue from "vue";

import Register from "../components/Register.vue";
import UpdateUser from "../components/UpdateUser.vue";
import AdminScanner from "../components/AdminScanning.vue";

export default Vue.extend({
  name: "Admin",
  components: {
    Register,
    UpdateUser,
    AdminScanner,
  },

  data: () => ({
    tab: 0,
  }),

  methods: {},
});
</script>
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
          <v-tab>User Action</v-tab>
          <v-tab>Room Message Scanning</v-tab>
        </v-tabs>

        <v-tabs-items v-model="tab">
          <v-tab-item>
            <v-row class="my-10" justify="center" align="center">
              <v-col class="text-center" cols="12" sm="4">
                <v-col><h2>Register User To Platform</h2></v-col>
                <v-form action="/register" method="POST" v-model="dataIsValid">
                  <v-text-field
                    v-model="email"
                    outlined
                    filled
                    type="email"
                    name="email"
                    required
                    placeholder="Email Address"
                    :rules="emailRules"
                    :disabled="detailsDisabled"
                  ></v-text-field>

                  <v-text-field
                    type="password"
                    v-model="password"
                    name="password"
                    outlined
                    filled
                    required
                    :rules="fillInputRules"
                    placeholder="Password"
                    :disabled="detailsDisabled"
                  ></v-text-field>

                  <v-text-field
                    v-model="name"
                    outlined
                    filled
                    name="name"
                    type="text"
                    required
                    :rules="fillInputRules"
                    placeholder="Full Name. Surname first."
                    :disabled="detailsDisabled"
                  ></v-text-field>

                  <v-menu
                    v-model="datePickerMenu"
                    :close-on-content-click="false"
                    transition="scale-transition"
                    offset-y
                    max-width="290px"
                    min-width="290px"
                  >
                    <template v-slot:activator="{ on, attrs }">
                      <v-text-field
                        outlined
                        filled
                        v-model="computedDateFormatted"
                        label="Date (read only text field)"
                        hint="MM/DD/YYYY format"
                        persistent-hint
                        readonly
                        v-bind="attrs"
                        name="DOB"
                        v-on="on"
                      ></v-text-field>
                    </template>
                    <v-date-picker
                      v-model="DOB"
                      no-title
                      @input="datePickerMenu = false"
                    ></v-date-picker>
                  </v-menu>

                  <v-text-field
                    outlined
                    filled
                    label="Parents Phone Number"
                  ></v-text-field>

                  <v-btn
                    type="submit"
                    :disabled="detailsDisabled || !dataIsValid"
                    tile
                    >Register</v-btn
                  >
                </v-form>
              </v-col>
            </v-row>
          </v-tab-item>

          <v-tab-item>
            <h2>Hello</h2>
          </v-tab-item>
        </v-tabs-items>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
"use strict";

import Vue from "vue";

export default Vue.extend({
  name: "Admin",

  data: () => ({
    tab: 0,
  }),

  methods: {},
});
</script>
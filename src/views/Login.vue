<template>
  <div
    align="center"
    style="
      display: grid;
      align-items: center;
      height: 100%;
      justify-content: center;
    "
  >
    <v-row class="fill-height" justify="center" align="center">
      <v-col cols="12">
        <v-img
          src="../assets/unilag.svg"
          contain
          height="300"
          width="300"
        ></v-img>
      </v-col>

      <v-col cols="12">
        <span v-if="isAdmin">Welcome Admin</span>
        <span v-else>Welcome to the University of Lagos chatting platform</span>
      </v-col>

      <v-col cols="12" md="5">
        <v-form action="/login" method="POST" v-model="detailsValid">
          <v-text-field
            v-model="email"
            outlined
            filled
            type="text"
            :rules="emailRules"
            name="email"
            required
            placeholder="Email Address"
            :disabled="detailsDisabled"
          ></v-text-field>

          <v-text-field
            type="password"
            v-model="password"
            outlined
            filled
            required
            :rules="passwordRules"
            name="password"
            placeholder="Password"
            :disabled="detailsDisabled"
          ></v-text-field>

          <v-alert v-if="signinErrorDetails" dense text type="error">{{
            signinErrorDetails
          }}</v-alert>

          <v-btn href="/register" color="green" class="mr-4" tile dark
            >Register</v-btn
          >
          <v-btn type="submit" :disabled="detailsDisabled || !detailsValid" tile
            >Log in</v-btn
          >
        </v-form>
      </v-col>
    </v-row>
  </div>
</template>

<script lang="ts">
"use strict";

import Vue from "vue";

export default Vue.extend({
  name: "Login",

  props: {
    isAdmin: Boolean,
  },

  data: () => ({
    detailsValid: false,
    detailsDisabled: false,

    signinErrorDetails: "",
    password: "",
    email: "",

    emailRules: [
      (v: string) => !!v || "E-mail is required",
      (v: string) => /.+@.+\..+/.test(v) || "E-mail must be valid",
    ],
    passwordRules: [(v: string) => !!v || "Password is required"],
  }),

  methods: {},
});
</script>
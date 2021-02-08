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

      <v-col class="text-center" cols="12">
        <span class="text-center mx-auto"
          >Welcome to the University of Lagos chatting platform</span
        >
      </v-col>

      <v-col cols="12" sm="6">
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

          <v-btn class="mr-4" dark href="/login" color="green" tile
            >Log In</v-btn
          >
          <v-btn type="submit" :disabled="detailsDisabled || !dataIsValid" tile
            >Register</v-btn
          >
        </v-form>
      </v-col>
    </v-row>
  </div>
</template>

<script lang="ts">
import Vue from "vue";

export default Vue.extend({
  name: "Register",
  data: () => ({
    detailsDisabled: false,
    dataIsValid: false,
    datePickerMenu: false,
    email: "",
    password: "",
    name: "",

    DOB: new Date().toISOString().substr(0, 10),

    emailRules: [
      (v: string) => !!v || "E-mail is required",
      (v: string) => /.+@.+\..+/.test(v) || "E-mail must be valid",
    ],
    fillInputRules: [(v: string) => !!v || "Field is required"],
  }),

  computed: {
    computedDateFormatted(): string {
      return this.formatDate(this.DOB);
    },
  },

  methods: {
    formatDate(date: string): string {
      if (!date) return "";

      const [year, month, day] = date.split("-");
      return `${month}/${day}/${year}`;
    },

    parseDate(date: string): string {
      if (!date) return "";

      const [month, day, year] = date.split("/");
      return `${year}-${month.padStart(2, "0")}-${day.padStart(2, "0")}`;
    },
  },
});
</script>

<style>
</style>
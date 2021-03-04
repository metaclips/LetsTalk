<template>
  <v-row class="my-10" justify="center" align="center">
    <v-col class="text-center" cols="12" sm="4">
      <v-col><h2>STUDENT REGISTRATION</h2></v-col>
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

        <v-text-field
          outlined
          disabled
          filled
          name="department"
          type="text"
          value="Computer Engineering"
          required
        ></v-text-field>

        <v-text-field
          outlined
          filled
          placeholder="Entry Year"
          name="year"
          type="number"
          required
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

        <v-btn type="submit" rounded :disabled="detailsDisabled || !dataIsValid"
          >Register</v-btn
        >
      </v-form>
    </v-col>
  </v-row>
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
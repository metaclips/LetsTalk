<template>
  <v-contianer fluid>
    <v-row class="mx-12 my-12">
      <v-col cols="12" sm="4">
        <v-text-field
          v-model="email"
          outlined
          filled
          type="email"
          name="email"
          required
          placeholder="Enter Students Email Address"
          :rules="emailRules"
          :disabled="detailsDisabled"
        ></v-text-field>

        <v-btn depressed class="text-center" color="green" dark>Fetch</v-btn>
      </v-col>

      <v-col class="text-center" cols="12">
        <v-row
          class="my-10"
          v-if="usersDetails"
          justify="center"
          align="center"
        >
          <v-col cols="12" sm="4">
            <v-form action="/register" method="POST" v-model="dataIsValid">
              <v-text-field
                rounded
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
                rounded
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
                rounded
                outlined
                filled
                name="department"
                type="text"
                value="Computer Engineering"
                required
              ></v-text-field>

              <v-text-field
                rounded
                outlined
                filled
                placeholder="Entry Year"
                name="year"
                type="number"
                required
              ></v-text-field>

              <v-select
                rounded
                filled
                :items="userBanned"
                label="Banned"
              ></v-select>

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
                    rounded
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
                rounded
                outlined
                filled
                label="Parents Phone Number"
              ></v-text-field>

              <v-btn
                type="submit"
                rounded
                :disabled="detailsDisabled || !dataIsValid"
                >Update Details</v-btn
              >
            </v-form>
          </v-col>
        </v-row>
      </v-col>
    </v-row>
  </v-contianer>
</template>

<script lang="ts">
import Vue from "vue";

import { UsersDetails } from "../views/Types";

export default Vue.extend({
  name: "UpdateUser",
  data: () => ({
    email: "",
    userBanned: ["True", "False"],
    usersDetails: {} as UsersDetails,
  }),
});
</script>
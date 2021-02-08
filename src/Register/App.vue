<template>
  <div id="app">
    <v-app>
      <v-main>
        <Home />
      </v-main>
      <Footer />
    </v-app>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Home from "../views/Register.vue";

export default Vue.extend({
  name: "App",
  components: {
    Home: Home,
  },

  data: () => ({
    //
  }),
});
</script>
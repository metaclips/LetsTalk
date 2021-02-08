<template>
  <div id="app">
    <v-app>
      <v-main>
        <Home />
      </v-main>
    </v-app>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Home from "../views/Home.vue";

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
<script>
import VFooter from './components/layout/VFooter.vue'
import VHeader from './components/layout/VHeader.vue'

export default {
  name: 'App',

  components: {
    VHeader,
    VFooter,
  },

  props: {},

  data() {
    return {}
  },

  computed: {},

  methods: {},

  watch: {},

  created() {},
  mounted() {},
  updated() {},
  beforeUnmount() {},
  unmounted() {},
}
</script>

<template>
  <div id="app">
    <div
      v-if="!$route.meta.authPage"
      class="flex flex-col items-center justify-between min-h-screen w-full bg-gray-50 pt-20"
    >
      <VHeader />
      <div class="container">
        <router-view />
      </div>
      <div class="pt-5 w-full">
        <VFooter />
      </div>
    </div>
    <div v-else>
      <router-view />
    </div>
  </div>
</template>

<style scoped></style>

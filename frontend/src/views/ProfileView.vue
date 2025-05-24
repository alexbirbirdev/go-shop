<script>
import VButton from '@/components/forms/VButton.vue'

export default {
  name: 'ProfileView',

  components: { VButton },

  props: {},

  data() {
    return {}
  },

  computed: {},

  methods: {
    logout() {
      localStorage.removeItem('token')
      this.$store.commit('SET_AUTH', false)
      this.$router.push('/auth/signin')
    }
  },

  watch: {},

  created() {},
  mounted() {},
  updated() {},
  beforeUnmount() {},
  unmounted() {},
}
</script>

<template>
  <div>
    <VButton variant="danger" @click="logout">Выйти</VButton>
  </div>
</template>

<style scoped></style>

<script>
export default {
  name: 'VHeader',

  components: {},

  props: {},

  data() {
    return {
      headerRouters: [
        {name: "Каталог", link: "/catalog"},
        {name: "Корзина", link: "/cart"},
        {name: "Избранное", link: "/favorite"},
        {name: "Заказы", link: "/orders"},
      ],
    }
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
  <header class="bg-white shadow-sm py-4 fixed top-0 left-0 w-full flex justify-center">
    <div class="container px-4 flex justify-between items-center">
      <h1 class="text-xl font-bold">GO Shop</h1>
      <nav class="flex items-center gap-4">
        <router-link :to="r.link" v-for="r in headerRouters" :key="r.id">{{ r.name }}</router-link>
      </nav>
    </div>
  </header>
</template>

<style scoped></style>

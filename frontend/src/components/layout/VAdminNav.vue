<script>
export default {
  name: 'VAdminNav',

  components: {},

  props: {},

  data() {
    return {
      links: [
        { name: 'Пользователи', path: '/admin/users' },
        { name: 'Товары', path: '/admin/products' },
        { name: 'Заказы', path: '/admin/orders' },
        { name: 'Категории', path: '/admin/categories' },
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
  <div class="grid gap-4" :class="'grid-cols-' + links.length">
    <router-link
      v-for="link in links"
      :key="link.path"
      :to="link.path"
      exact-active-class="'bg-sky-200 '"
      :active-class="'bg-sky-200'"
      class="flex items-center justify-center aspect-[16/4] bg-sky-100 rounded-2xl hover:bg-sky-300 duration-200 text-2xl font-bold text-sky-950"
      >{{ link.name }}</router-link
    >
  </div>
</template>

<style scoped></style>

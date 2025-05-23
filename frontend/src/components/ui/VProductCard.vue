<script>
export default {
  name: 'VProductCard',

  components: {},

  props: {
    product: {
      type: Object,
      required: true
    }
  },

  data() {
    return {};
  },

  computed: {},

  methods: {},

  watch: {},

  created() {},
  mounted() {},
  updated() {},
  beforeUnmount() {},
  unmounted() {},
};
</script>

<template>
  <div class="bg-white rounded-lg shadow overflow-hidden">
    <img :src="product.image" :alt="product.name" class="w-full h-48 object-cover" />
    <div class="p-4">
      <h3 class="font-semibold text-lg">{{ product.name }}</h3>
      <p class="text-gray-600 mt-1">{{ product.price }} ₽</p>
      <button
        @click="$emit('add-to-cart', product)"
        class="mt-4 w-full bg-blue-500 text-white py-2 rounded hover:bg-blue-600"
      >
        В корзину
      </button>
    </div>
  </div>
</template>

<style scoped></style>

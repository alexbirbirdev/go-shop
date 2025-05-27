<script>
export default {
  name: 'VProductCard',

  components: {},

  props: {
    product: {
      type: Object,
      required: true,
    },
  },

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
  <div class="bg-white rounded-lg shadow overflow-hidden flex flex-col justify-start">
    <router-link :to="'/catalog/' + product.id" class="relative before:w-full before:h-full before:absolute before:duration-200 hover:before:bg-sky-50/50">
      <img
        v-if="!product.images.length"
        src="https://placehold.co/600x400"
        :alt="product.name"
        class="w-full h-48 object-cover"
      />
      <img
        v-else
        :src="product.images[0].image_url"
        :alt="product.name"
        class="w-full h-48 object-cover"
      />
    </router-link>

    <div class="px-4 py-2">
      <router-link :to="'/catalog/' + product.id" class="font-semibold duration-200 hover:text-blue-500 text-md leading-[120%] line-clamp-2">{{ product.name }}</router-link>
      <p class="text-gray-600 mt-1">{{ product.price }} ₽</p>
    </div>
    <!-- <div class="px-4 py-2 mt-auto">
      <button
        @click="$emit('add-to-cart', product)"
        class="w-full bg-blue-500 text-white py-2 rounded cursor-pointer hover:bg-blue-600"
      >
        В корзину
      </button>
    </div> -->
  </div>
</template>

<style scoped></style>

<script>
import VBlockLoader from '@/components/loaders/VBlockLoader.vue'
import VButton from '@/components/forms/VButton.vue'
import axios from 'axios'
export default {
  name: 'OrdersView',

  components: { VBlockLoader, VButton },

  props: {},

  data() {
    return {
      isLoading: true,
      orders: [],
    }
  },

  computed: {},

  methods: {
    async getOrders() {
      try {
        this.isLoading = true
        const response = await axios.get('/api/orders/', {
          headers: {
            Authorization: localStorage.getItem('token'),
          },
        })
        this.orders = response.data.orders
      } catch (error) {
        console.log(error)
      } finally {
        this.isLoading = false
      }
    },
  },

  watch: {},

  created() {
    this.getOrders()
  },
  mounted() {},
  updated() {},
  beforeUnmount() {},
  unmounted() {},
}
</script>

<template>
  <div class="max-w-3xl mx-auto">
    <div class="text-4xl font-bold mb-10">Заказы</div>
    <div class="grid gap-6" v-if="!isLoading && orders.length">
      <router-link
        :to="'/orders/' + String(order.ID)"
        class="shadow rounded-xl duration-200 hover:shadow-xl"
        v-for="order in orders"
        :key="order"
      >
        <div class="w-full p-4 bg-neutral-100 flex items-center justify-between">
          <div class="">
            Заказ №<span>{{ order.ID }}</span> от <span>{{ order.CreatedAt }}</span>
          </div>
          <div class="">
            {{ order.status }}
          </div>
        </div>
        <div class="p-4 grid gap-2 grid-cols-6">
          <div
            class="flex flex-col items-start gap-1"
            v-for="product in order.order_items"
            :key="product"
          >
            <div
              class="aspect-square flex items-center select-none justify-center border-2 border-neutral-300 rounded-lg overflow-hidden"
            >
              <img
                v-if="!product.product_variant.Product.images.length"
                src="https://placehold.co/600x400"
                :alt="product.product_variant.Product.name"
                class="max-h-full max-w-full"
              />
              <img
                v-else
                :src="product.product_variant.Product.images[0].image_url"
                :alt="product.product_variant.Product.name"
                class="max-h-full max-w-full"
              />
            </div>
            <div class="text-xs leading-[120%]">{{ product.product_variant.Product.name }}</div>
          </div>
        </div>
      </router-link>
    </div>
    <div class="grid gap-6" v-if="isLoading">
      <div class="shadow rounded-xl">
        <div class="w-full p-4 bg-neutral-100 flex items-center justify-between">
          <VBlockLoader class="w-60 h-6"></VBlockLoader>
          <VBlockLoader class="w-20 h-6"></VBlockLoader>
        </div>
        <div class="p-4 grid gap-2 grid-cols-6">
          <div class="flex flex-col items-start gap-1" v-for="product in 5" :key="product">
            <VBlockLoader class="w-full aspect-square"></VBlockLoader>
            <VBlockLoader class="w-full h-4"></VBlockLoader>
          </div>
        </div>
      </div>
      <div class="shadow rounded-xl duration-200 hover:shadow-xl">
        <div class="w-full p-4 bg-neutral-100 flex items-center justify-between">
          <VBlockLoader class="w-60 h-6"></VBlockLoader>
          <VBlockLoader class="w-20 h-6"></VBlockLoader>
        </div>
        <div class="p-4 grid gap-2 grid-cols-6">
          <div class="flex flex-col items-start gap-1" v-for="product in 3" :key="product">
            <VBlockLoader class="w-full aspect-square"></VBlockLoader>
            <VBlockLoader class="w-full h-4"></VBlockLoader>
          </div>
        </div>
      </div>
      <div class="shadow rounded-xl duration-200 hover:shadow-xl">
        <div class="w-full p-4 bg-neutral-100 flex items-center justify-between">
          <VBlockLoader class="w-60 h-6"></VBlockLoader>
          <VBlockLoader class="w-20 h-6"></VBlockLoader>
        </div>
        <div class="p-4 grid gap-2 grid-cols-6">
          <div class="flex flex-col items-start gap-1" v-for="product in 2" :key="product">
            <VBlockLoader class="w-full aspect-square"></VBlockLoader>
            <VBlockLoader class="w-full h-4"></VBlockLoader>
          </div>
        </div>
      </div>
    </div>

    <div
      v-if="!orders.length && !isLoading"
      class="w-full p-5 text-2xl flex items-center gap-4 flex-col justify-center bg-neutral-200 rounded-2xl"
    >
      Вы не сделали ни одного заказа
      <div>
        <router-link to="/catalog/">
          <VButton class="text-xs">В каталог</VButton>
        </router-link>
      </div>
    </div>
  </div>
</template>

<style scoped></style>

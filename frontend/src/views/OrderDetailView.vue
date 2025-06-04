<script>
import axios from 'axios'
export default {
  name: 'OrderDetailView',

  components: {},

  props: {},

  data() {
    return {
      order: null,
      isLoading: true,
    }
  },

  computed: {},

  methods: {
    async getOrders() {
      try {
        this.isLoading = true
        const response = await axios.get('/api/orders/' + this.$route.params.id, {
          headers: {
            Authorization: localStorage.getItem('token'),
          },
        })
        this.order = response.data.order
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
  <div class="max-w-3xl mx-auto" v-if="!isLoading">
    <div class="text-4xl font-bold mb-10 items-center flex justify-start gap-5">
      Заказ #{{ order.ID }}
      <span class="bg-yellow-100 rounded-2xl px-5 py-2 text-sm">{{ order.status }}</span>
    </div>
    <div class="mb-2 text-sm text-gray-500">
      Заказ от <span class="font-semibold">{{ order.CreatedAt }}</span>
    </div>
    <div class="grid gap-5">
      <div class="grid gap-4 bg-neutral-50 rounded-2xl shadow px-4 py-6">
        <div class="flex justify-between gap-4" v-for="product in order.order_items" :key="product">
          <div
            class="aspect-square w-35 flex items-center select-none justify-center border-2 border-neutral-300 rounded-lg overflow-hidden"
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
          <div class="flex flex-1 flex-col justify-between">
            <div class="">
              <div class="text-lg">
                {{ product.product_variant.Product.name }}
              </div>
              <div class="text-sm mt-1 text-gray-500">
                {{ product.product_variant.Product.description }}
              </div>
            </div>

            <div class="flex items-center justify-start gap-4">
              <div class="text-xs">
                Кол-во: <span class="text-sm font-bold">{{ product.stock }}</span>
              </div>
              <div class="text-xs">
                Стоимость: <span class="text-sm font-bold">{{ product.price }}₽</span> / шт
              </div>
              <div class="text-xs">
                Вариант:
                <span class="py-1 px-3 text-xs rounded-2xl bg-neutral-200 duration-200 select-none">
                  {{ product.product_variant.name }}</span
                >
              </div>
            </div>
          </div>
        </div>
      </div>
      <div class="grid grid-cols-3 gap-10">
        <div class="flex flex-col gap-2">
          <div class="text-xl font-medium">Получатель</div>
          <div class="grid gap-2">
            <div class="grid">
              <div class="text-xs text-gray-500">Телефон</div>
              <div class="text-xs">{{ order.phone }}</div>
            </div>
            <div class="grid">
              <div class="text-xs text-gray-500">Имя, фамилия</div>
              <div class="text-xs">{{ order.first_name }} {{ order.last_name }}</div>
            </div>
          </div>
        </div>
        <div class="flex flex-col gap-2">
          <div class="text-xl font-medium">Адрес доставки</div>
          <div class="grid gap-2">
            <div class="grid">
              <div class="text-xs text-gray-500">Адрес</div>
              <div class="text-xs">{{ order.address }}</div>
            </div>
            <div class="grid">
              <div class="text-xs text-gray-500">Комментарий</div>
              <div class="text-xs">{{ order.address_com }}</div>
            </div>
          </div>
        </div>
        <div class="flex flex-col gap-2">
          <div class="text-xl font-medium flex items-center justify-start gap-4">
            Итого
            <div class="text-xl font-bold uppercase text-nowrap">{{ order.total_price }} ₽</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped></style>

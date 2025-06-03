<script>
import axios from 'axios'
export default {
  name: 'AdminOrders',

  components: {},

  props: {},

  data() {
    return {
      orders: [],
      isLoading: true,
      statusOptions: ['pending', 'processing', 'shipped', 'completed', 'cancelled'],
      isUpdating: false,
    }
  },

  computed: {},

  methods: {
    async updateOrderStatus(orderId, newStatus) {
      if (this.isUpdating) return
      try {
        this.isUpdating = true
        await axios.put(
          `http://localhost:8080/admin/orders/${orderId}/status`,
          { status: newStatus },
          {
            headers: {
              Authorization: localStorage.getItem('token'),
            },
          },
        )
      } catch (error) {
        console.error(`Ошибка при обновлении статуса заказа ${orderId}:`, error)
      } finally {
        this.isUpdating = false
      }
    },
    async getOrders() {
      try {
        this.isLoading = true
        const response = await axios.get('http://localhost:8080/admin/orders/', {
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
  <div>
    <div class="">
      <div class="grid gap-6" v-if="!isLoading && orders.length">
        <div
          class="shadow rounded-xl duration-200 hover:shadow-lg"
          v-for="order in orders"
          :key="order"
        >
          <div class="w-full p-4 bg-neutral-100">
            <div class="flex items-center justify-between">
              <div class="">
                Заказ №<span class="font-bold">{{ order.ID }}</span> от
                <span class="font-bold">{{ order.CreatedAt }}</span>
              </div>
              <div class="">
                <select
                  class="px-2 py-1 rounded border border-neutral-300 text-sm"
                  v-model="order.status"
                  @change="updateOrderStatus(order.ID, order.status)"
                >
                  <option v-for="status in statusOptions" :key="status" :value="status">
                    {{ status }}
                  </option>
                </select>
              </div>
            </div>
            <div class="grid grid-cols-2 gap-10 mt-2">
              <div class="flex flex-col gap-2">
                <div class="text-base font-medium">Получатель</div>
                <div class="grid gap-2 grid-cols-2">
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
                <div class="text-base font-medium">Адрес доставки</div>
                <div class="grid gap-2 grid-cols-2">
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
              <div class="flex items-center justify-between w-full *:text-[10px]">
                <span
                  class="py-1 px-2 rounded-2xl duration-200 select-none bg-neutral-200 text-neutral-900"
                >
                  {{ product.product_variant.name }}</span
                >
                <div class="leading-[120%]">
                  <strong>{{ product.stock }}</strong> шт
                </div>
                <div class="leading-[120%]">
                  <strong>{{ product.price }}</strong> ₽
                </div>
              </div>

              <div class="text-xs leading-[120%]">{{ product.product_variant.Product.name }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped></style>

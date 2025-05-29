<script>
import { mapMutations, mapGetters } from 'vuex'
import axios from 'axios'
import VButton from '@/components/forms/VButton.vue'
import VBlockLoader from '@/components/loaders/VBlockLoader.vue'

export default {
  name: 'CartView',

  components: { VButton, VBlockLoader },

  props: {},

  data() {
    return {
      products: [],
      total: 0,
      isDecLoading: false,
      isIncLoading: false,
      isLoading: true,
    }
  },

  computed: {
    ...mapGetters('cart', {
      cartCount: 'count',
    }),
  },

  methods: {
    ...mapMutations('cart', ['INCREMENT', 'DECREMENT', 'SET_COUNT']),
    async getCart() {
      try {
        this.isLoading = true
        const response = await axios.get('http://localhost:8080/cart/', {
          headers: {
            Authorization: localStorage.getItem('token'),
          },
        })
        this.total = response.data.total
        this.products = response.data.cart_items
      } catch (error) {
        console.log(error)
      } finally {
        this.isLoading = false
      }
    },
    async deleteCartItem(id, price, quantity) {
      try {
        await axios.delete('http://localhost:8080/cart/' + id, {
          headers: {
            Authorization: localStorage.getItem('token'),
          },
        })
        this.products = this.products.filter((p) => p.product_variant.ID !== id)
        this.total -= price * quantity
        this.$store.commit('cart/SET_COUNT', this.cartCount - quantity)
      } catch (error) {
        console.log(error)
      }
    },
    async decCartItem(id, price) {
      if (this.isDecLoading) return
      try {
        this.isDecLoading = true
        const item = this.products.find((p) => p.ID === id)
        if (item && item.quantity > 1) {
          axios.post(
            'http://localhost:8080/cart/decrement/' + id,
            {},
            {
              headers: {
                Authorization: localStorage.getItem('token'),
              },
            },
          )

          item.quantity -= 1
          this.total -= price
          this.$store.commit('cart/DECREMENT')
        }
      } catch (error) {
        console.log(error)
      } finally {
        this.isDecLoading = false
      }
    },
    async incCartItem(id, price) {
      if (this.isIncLoading) return

      try {
        this.isIncLoading = true
        const item = this.products.find((p) => p.ID === id)
        if (item.quantity < item.product_variant.stock) {
          await axios.post(
            'http://localhost:8080/cart/increment/' + id,
            {},
            {
              headers: {
                Authorization: localStorage.getItem('token'),
              },
            },
          )
          item.quantity += 1
          this.total += price * 1
          this.$store.commit('cart/INCREMENT')
        }
      } catch (error) {
        console.log(error)
      } finally {
        this.isIncLoading = false
      }
    },
  },

  watch: {},

  created() {
    this.getCart()
  },
  mounted() {},
  updated() {},
  beforeUnmount() {},
  unmounted() {},
}
</script>

<template>
  <div class="max-w-5xl mx-auto">
    <h1 class="text-4xl font-bold">Корзина</h1>
    <div class="mt-10 grid grid-cols-3 gap-10" v-if="isLoading">
      <div class="col-span-2">
        <div
          v-for="cartItem in 5"
          :key="cartItem"
          class="border-b-1 border-neutral-200 flex gap-5 py-5"
        >
          <VBlockLoader class="aspect-square w-50" />
          <div class="flex flex-col justify-between items-start flex-1">
            <div class="flex items-start gap-4 justify-between w-full">
              <div class="grid gap-2 w-full">
                <VBlockLoader class="w-50 h-7" />
                <VBlockLoader class="w-12 h-6" />
                <VBlockLoader class="w-20 h-4" />
              </div>
              <VBlockLoader class="w-10 aspect-square" />
            </div>
            <div class="inline-flex gap-2">
              <VBlockLoader class="h-10 aspect-square" />
              <VBlockLoader class="h-10 w-15" />
              <VBlockLoader class="h-10 aspect-square" />
            </div>
          </div>
        </div>
      </div>
      <div class="">
        <div class="grid gap-5 bg-neutral-100 rounded-2xl p-5">
          <div class="flex items-center justify-between">
            <VBlockLoader class="w-17 h-7" />
            <VBlockLoader class="w-12 h-6" />
          </div>
          <VBlockLoader class="w-full h-10" />
        </div>
      </div>
    </div>
    <div class="mt-10 grid grid-cols-3 gap-10" v-if="!isLoading && products.length">
      <div class="col-span-2">
        <div
          v-for="cartItem in products"
          :key="cartItem"
          class="border-b-1 border-neutral-200 flex gap-5 py-5"
        >
          <router-link
            :to="'/catalog/' + cartItem.product_variant.Product.ID"
            class="aspect-square flex items-center select-none justify-center border-2 border-neutral-300 rounded-lg w-50 overflow-hidden relative before:w-full before:h-full before:absolute before:duration-200 hover:before:bg-sky-50/50"
          >
            <img
              v-if="!cartItem.product_variant.Product.images.length"
              src="https://placehold.co/600x400"
              :alt="cartItem.product_variant.Product.name"
              class="max-h-full max-w-full"
            />
            <img
              v-else
              :src="cartItem.product_variant.Product.images[0].image_url"
              :alt="cartItem.product_variant.Product.name"
              class="max-h-full max-w-full"
            />
          </router-link>
          <div class="flex flex-col justify-between items-start flex-1">
            <div class="flex items-start gap-4 justify-between w-full">
              <div class="grid gap-2 w-full">
                <router-link
                  :to="'/catalog/' + cartItem.product_variant.Product.ID"
                  class="text-lg font-bold duration-200 hover:text-blue-600 select-none"
                  >{{ cartItem.product_variant.Product.name }}</router-link
                >
                <div class="">
                  <span
                    class="py-1 px-3 text-xs rounded-2xl bg-neutral-200 duration-200 select-none"
                  >
                    {{ cartItem.product_variant.name }}</span
                  >
                </div>
                <div class="font-bold text-xs select-none">
                  {{ cartItem.product_variant.price }} ₽ / шт
                </div>
              </div>
              <div
                @click="
                  deleteCartItem(
                    cartItem.product_variant.ID,
                    cartItem.product_variant.price,
                    cartItem.quantity,
                  )
                "
                class="p-2 rounded-xl cursor-pointer duration-200 hover:bg-neutral-200"
              >
                <svg
                  width="24"
                  height="24"
                  viewBox="0 0 24 24"
                  fill="none"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    d="M17.4698 5.46984C17.7627 5.17695 18.2375 5.17695 18.5304 5.46984C18.8233 5.76274 18.8233 6.2375 18.5304 6.53039L13.0607 12.0001L18.5304 17.4698L18.5821 17.5265C18.8225 17.8211 18.805 18.2558 18.5304 18.5304C18.2558 18.805 17.8211 18.8225 17.5265 18.5821L17.4698 18.5304L12.0001 13.0607L6.53039 18.5304C6.2375 18.8233 5.76274 18.8233 5.46984 18.5304C5.17695 18.2375 5.17695 17.7627 5.46984 17.4698L10.9396 12.0001L5.46984 6.53039L5.41809 6.47375C5.17778 6.17917 5.19524 5.74445 5.46984 5.46984C5.74445 5.19524 6.17917 5.17778 6.47375 5.41809L6.53039 5.46984L12.0001 10.9396L17.4698 5.46984Z"
                    fill="black"
                  />
                </svg>
              </div>
            </div>
            <div class="inline-flex gap-2">
              <div
                :disabled="cartItem.quantity <= 1"
                @click="decCartItem(cartItem.ID, cartItem.product_variant.price)"
                :class="
                  cartItem.quantity <= 1
                    ? 'opacity-50 scale-90 cursor-not-allowed'
                    : ' cursor-pointer hover:bg-neutral-300'
                "
                class="p-2 bg-neutral-200 rounded-xl duration-200"
              >
                <svg
                  v-if="isDecLoading"
                  width="24"
                  height="24"
                  class="*:fill-white animate-spin"
                  viewBox="0 0 24 24"
                  fill="none"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    d="M12 2.25C17.3848 2.25 21.75 6.61522 21.75 12C21.75 17.3848 17.3848 21.75 12 21.75C6.61522 21.75 2.25 17.3848 2.25 12C2.25 10.4448 2.61447 8.97237 3.26367 7.66602C3.44804 7.29514 3.89863 7.1438 4.26953 7.32812C4.6404 7.51249 4.79174 7.96308 4.60742 8.33398C4.05901 9.43754 3.75 10.6816 3.75 12C3.75 16.5563 7.44365 20.25 12 20.25C16.5563 20.25 20.25 16.5563 20.25 12C20.25 7.44365 16.5563 3.75 12 3.75C11.5858 3.75 11.25 3.41421 11.25 3C11.25 2.58579 11.5858 2.25 12 2.25Z"
                    fill="black"
                  />
                </svg>
                <svg
                  v-else
                  width="24"
                  height="24"
                  viewBox="0 0 24 24"
                  fill="none"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    d="M18 11.25C18.4142 11.25 18.75 11.5858 18.75 12C18.75 12.4142 18.4142 12.75 18 12.75H6C5.58579 12.75 5.25 12.4142 5.25 12C5.25 11.5858 5.58579 11.25 6 11.25H18Z"
                    fill="black"
                  />
                </svg>
              </div>
              <div class="h-full select-none p-2 rounded-xl text-center w-15 bg-neutral-200">
                <span>{{ cartItem.quantity }}</span>
              </div>
              <div
                @click="incCartItem(cartItem.ID, cartItem.product_variant.price)"
                :class="
                  cartItem.quantity == cartItem.product_variant.stock
                    ? 'opacity-50 scale-90 cursor-not-allowed'
                    : ' cursor-pointer hover:bg-neutral-300'
                "
                class="p-2 bg-neutral-200 rounded-xl duration-200"
              >
                <svg
                  v-if="isIncLoading"
                  width="24"
                  height="24"
                  class="*:fill-white animate-spin"
                  viewBox="0 0 24 24"
                  fill="none"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    d="M12 2.25C17.3848 2.25 21.75 6.61522 21.75 12C21.75 17.3848 17.3848 21.75 12 21.75C6.61522 21.75 2.25 17.3848 2.25 12C2.25 10.4448 2.61447 8.97237 3.26367 7.66602C3.44804 7.29514 3.89863 7.1438 4.26953 7.32812C4.6404 7.51249 4.79174 7.96308 4.60742 8.33398C4.05901 9.43754 3.75 10.6816 3.75 12C3.75 16.5563 7.44365 20.25 12 20.25C16.5563 20.25 20.25 16.5563 20.25 12C20.25 7.44365 16.5563 3.75 12 3.75C11.5858 3.75 11.25 3.41421 11.25 3C11.25 2.58579 11.5858 2.25 12 2.25Z"
                    fill="black"
                  />
                </svg>
                <svg
                  v-else
                  width="24"
                  height="24"
                  viewBox="0 0 24 24"
                  fill="none"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    d="M12 5.25C12.4142 5.25 12.75 5.58579 12.75 6V11.25H18C18.4142 11.25 18.75 11.5858 18.75 12C18.75 12.4142 18.4142 12.75 18 12.75H12.75V18C12.75 18.4142 12.4142 18.75 12 18.75C11.5858 18.75 11.25 18.4142 11.25 18V12.75H6C5.58579 12.75 5.25 12.4142 5.25 12C5.25 11.5858 5.58579 11.25 6 11.25H11.25V6C11.25 5.58579 11.5858 5.25 12 5.25Z"
                    fill="black"
                  />
                </svg>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div class="">
        <div class="grid gap-5 bg-neutral-200 rounded-2xl p-5">
          <div class="flex items-center justify-between">
            <div class="text-xl font-bold">Итого:</div>
            <div class="">{{ total }} ₽</div>
          </div>
          <VButton>Оформить заказ</VButton>
        </div>
      </div>
    </div>

    <div
      v-if="!products.length && !isLoading"
      class="w-full p-5 text-2xl flex mt-10 items-center gap-4 flex-col justify-center bg-neutral-200 rounded-2xl"
    >
      В корзине пустно!
      <div>
        <router-link to="/catalog/">
          <VButton class="text-xs">В каталог</VButton>
        </router-link>
      </div>
    </div>
  </div>
</template>

<style scoped></style>

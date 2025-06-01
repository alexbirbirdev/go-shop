<script>
import { mapMutations } from 'vuex'
import axios from 'axios'
import VButton from '@/components/forms/VButton.vue'
import VBlockLoader from '@/components/loaders/VBlockLoader.vue'

export default {
  name: 'ProductView',

  components: {
    VButton,
    VBlockLoader,
  },

  props: {},

  data() {
    return {
      product: {},
      isLoading: true,
      activePrice: 0,
      activeVariant: null,
      cartLoading: false,
      favLoading: false,
    }
  },

  computed: {
    activeVariantData() {
      if (!this.product.product_variants) return null
      return this.product.product_variants.find((v) => v.id === this.activeVariant) || null
    },
  },

  methods: {
    ...mapMutations('cart', ['INCREMENT', 'DECREMENT']),
    makeActive(id, price) {
      this.activeVariant = id
      this.activePrice = price
    },
    ...mapMutations('favorites', ['INCREMENT', 'DECREMENT']),
    async getProduct() {
      try {
        this.isLoading = true
        const response = await axios.get(
          'http://localhost:8080/products/' + this.$route.params.id,

          {
            headers: {
              Authorization: localStorage.getItem('token'),
            },
          },
        )
        this.activePrice = response.data.product.price

        const match = response.data.product.product_variants.find(
          (v) => v.price === response.data.product.price,
        )
        this.activeVariant = match ? match.id : null
        this.product = response.data.product
      } catch (error) {
        console.log(error)
      } finally {
        this.isLoading = false
      }
    },

    async toggleFavorite() {
      if (!this.activeVariantData) return
      if (this.favLoading) return
      if (!localStorage.getItem("role")) {
        this.$router.push('/auth/signin?redirect=/catalog/' + this.product.id)
        return
      }

      const variantId = this.activeVariantData.id
      const isFav = this.activeVariantData.is_fav

      try {
        this.favLoading = true
        if (isFav) {
          await axios.delete(`http://localhost:8080/favorites/${variantId}`, {
            headers: {
              Authorization: localStorage.getItem('token'),
            },
          })
          this.DECREMENT()
          this.activeVariantData.is_fav = false
        } else {
          await axios.post(
            `http://localhost:8080/favorites/`,
            {
              product_variant_id: variantId,
            },
            {
              headers: {
                Authorization: localStorage.getItem('token'),
              },
            },
          )
          this.INCREMENT()
          this.activeVariantData.is_fav = true
        }
      } catch (error) {
        console.error('Ошибка при обновлении избранного:', error)
      } finally {
        this.favLoading = false
      }
    },
    async toggleCart() {
      if (!this.activeVariantData) return
      if (this.cartLoading) return
      if (!localStorage.getItem("role")) {
        this.$router.push('/auth/signin?redirect=/catalog/' + this.product.id)
        return
      }

      const variantId = this.activeVariantData.id
      const inCart = this.activeVariantData.in_cart

      try {
        this.cartLoading = true
        if (inCart) {
          await axios.delete(`http://localhost:8080/cart/${variantId}`, {
            headers: {
              Authorization: localStorage.getItem('token'),
            },
          })
          this.activeVariantData.in_cart = false
          this.$store.commit('cart/DECREMENT')
        } else {
          await axios.post(
            `http://localhost:8080/cart/`,
            {
              product_variant_id: variantId,
              quantity: 1,
            },
            {
              headers: {
                Authorization: localStorage.getItem('token'),
              },
            },
          )
          this.$store.commit('cart/INCREMENT')
          this.activeVariantData.in_cart = true
        }
      } catch (error) {
        console.error('Ошибка при обновлении избранного:', error)
      } finally {
        this.cartLoading = false
      }
    },
  },

  watch: {},

  created() {
    this.getProduct()
  },
  mounted() {},
  updated() {},
  beforeUnmount() {},
  unmounted() {},
}
</script>

<template>
  <div class="max-w-4xl mx-auto">
    <div class="grid grid-cols-2 gap-10" v-if="isLoading">
      <div class="">
        <VBlockLoader class="w-full aspect-square" />
        <div class="grid-cols-4 mt-4 grid gap-2">
          <VBlockLoader class="w-full aspect-square" v-for="i in 4" :key="i" />
        </div>
      </div>
      <div class="flex flex-col gap-4">
        <VBlockLoader class="w-full h-10" />
        <VBlockLoader class="w-40 h-8" />
        <div class="grid gap-2">
          <VBlockLoader class="w-full h-4" />
          <VBlockLoader class="w-full h-4" />
          <VBlockLoader class="w-1/2 h-4" />
        </div>
        <div class="flex flex-wrap gap-2">
          <VBlockLoader class="w-12 h-10 odd:w-16" v-for="variant in 5" :key="variant" />
        </div>
        <div class="flex gap-2">
          <VBlockLoader class="flex-1 h-10" />
          <VBlockLoader class="w-10 h-10" />
        </div>
      </div>
    </div>
    <div class="grid grid-cols-2 gap-10" v-if="!isLoading">
      <div class="">
        <div
          class="aspect-square overflow-hidden rounded-2xl border-4 border-neutral-100 flex items-center justify-center"
        >
          <img
            v-if="!product.images.length"
            src="https://placehold.co/600x400"
            :alt="product.name"
            class="max-h-full max-w-full"
          />
          <img
            v-else
            :src="product.images[0].image_url"
            :alt="product.name"
            class="max-h-full max-w-full"
          />
        </div>
        <div class="grid-cols-4 mt-4 grid gap-2" v-if="product.images.length > 1">
          <div
            v-for="image in product.images"
            :key="image"
            class="aspect-square overflow-hidden rounded-2xl border-1 border-sky-900 flex items-center justify-center"
          >
            <img class="max-w-full max-h-full" :src="image.image_url" alt="" />
          </div>
        </div>
      </div>
      <div class="flex flex-col gap-4">
        <div class="text-4xl font-bold">{{ product.name }}</div>
        <div class="text-2xl">{{ activePrice }} ₽</div>
        <div class="">
          {{ product.description }}
        </div>
        <div class="flex flex-wrap gap-2">
          <span
            v-for="variant in product.product_variants"
            :key="variant"
            class="py-2 px-4 rounded-2xl bg-neutral-200 duration-200 select-none"
            :class="
              activeVariant == variant.id
                ? 'bg-sky-200 text-sky-900 cursor-default'
                : 'cursor-pointer hover:bg-sky-100'
            "
            @click="makeActive(variant.id, variant.price)"
          >
            {{ variant.name }}</span
          >
        </div>
        <div class="flex gap-2 max-w-lg">
          <VButton
            :disabled="cartLoading"
            @click="toggleCart"
            v-if="!activeVariantData?.in_cart"
            class="flex-1"
            ><span v-if="cartLoading" class="animate-spin">
              <svg
                width="24"
                height="24"
                class="*:fill-white"
                viewBox="0 0 24 24"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  d="M12 2.25C17.3848 2.25 21.75 6.61522 21.75 12C21.75 17.3848 17.3848 21.75 12 21.75C6.61522 21.75 2.25 17.3848 2.25 12C2.25 10.4448 2.61447 8.97237 3.26367 7.66602C3.44804 7.29514 3.89863 7.1438 4.26953 7.32812C4.6404 7.51249 4.79174 7.96308 4.60742 8.33398C4.05901 9.43754 3.75 10.6816 3.75 12C3.75 16.5563 7.44365 20.25 12 20.25C16.5563 20.25 20.25 16.5563 20.25 12C20.25 7.44365 16.5563 3.75 12 3.75C11.5858 3.75 11.25 3.41421 11.25 3C11.25 2.58579 11.5858 2.25 12 2.25Z"
                  fill="black"
                />
              </svg> </span
            ><span v-else>В корзину</span></VButton
          >
          <VButton
            :disabled="cartLoading"
            @click="toggleCart"
            v-if="activeVariantData?.in_cart"
            :variant="'danger'"
            class="flex-1"
            ><span v-if="cartLoading" class="animate-spin">
              <svg
                width="24"
                height="24"
                class="*:fill-white"
                viewBox="0 0 24 24"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  d="M12 2.25C17.3848 2.25 21.75 6.61522 21.75 12C21.75 17.3848 17.3848 21.75 12 21.75C6.61522 21.75 2.25 17.3848 2.25 12C2.25 10.4448 2.61447 8.97237 3.26367 7.66602C3.44804 7.29514 3.89863 7.1438 4.26953 7.32812C4.6404 7.51249 4.79174 7.96308 4.60742 8.33398C4.05901 9.43754 3.75 10.6816 3.75 12C3.75 16.5563 7.44365 20.25 12 20.25C16.5563 20.25 20.25 16.5563 20.25 12C20.25 7.44365 16.5563 3.75 12 3.75C11.5858 3.75 11.25 3.41421 11.25 3C11.25 2.58579 11.5858 2.25 12 2.25Z"
                  fill="black"
                />
              </svg> </span
            ><span v-else>Убрать из корзины</span></VButton
          >
          <router-link :disabled="cartLoading" to="/cart/" v-if="activeVariantData?.in_cart"
            ><VButton class="">Перейти в корзину</VButton></router-link
          >

          <div
            :disabled="favLoading"
            :class="activeVariantData?.is_fav ? 'bg-red-200' : 'bg-neutral-200'"
            @click="toggleFavorite"
            class="w-10 h-10 p-2 bg-neutral-200 cursor-pointer rounded-lg flex items-center justify-center"
          >
            <svg
              v-if="favLoading"
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
              class="w-full *:fill-sky-950"
              viewBox="0 0 24 24"
              fill="none"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                d="M12.4891 5.92615C14.6343 3.73233 18.1651 3.73227 20.3103 5.92615C22.3882 8.05145 22.3881 11.4473 20.3103 13.5726L13.2107 20.8334C12.5877 21.4705 11.5878 21.5102 10.9178 20.9525L10.7889 20.8334L3.68925 13.5726C1.61143 11.4473 1.61136 8.05148 3.68925 5.92615C5.83453 3.73233 9.3653 3.73227 11.5105 5.92615L11.9998 6.42615L12.4891 5.92615ZM19.2381 6.97498C17.6812 5.38286 15.1192 5.38292 13.5623 6.97498L12.5359 8.02381C12.3949 8.16803 12.2015 8.24937 11.9998 8.24939C11.798 8.24939 11.6048 8.16806 11.4637 8.02381L10.4373 6.97498C8.88045 5.38286 6.3184 5.38292 4.76152 6.97498C3.25388 8.51732 3.25394 10.9814 4.76152 12.5238L11.8611 19.7845L11.8914 19.8099C11.9674 19.8613 12.0718 19.8528 12.1385 19.7845L19.2381 12.5238C20.7456 10.9814 20.7457 8.51728 19.2381 6.97498Z"
                fill="black"
              />
            </svg>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped></style>

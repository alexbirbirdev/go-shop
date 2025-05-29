<script>
import axios from 'axios'
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
    return {
      isDeleting: false,
    }
  },
  emits: ['deleted'],

  computed: {},

  methods: {
    
    async deleteFavorite(id) {
      if (this.isDeleting) return
      try {
        this.isDeleting = true
        await axios.delete('http://localhost:8080/favorites/' + id, {
          headers: {
            Authorization: localStorage.getItem('token'),
          },
        })

        this.$emit('deleted', id)
      } catch (error) {
        console.log(error)
      } finally {
        this.isDeleting = false
      }
    },
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
  <div
    class="bg-white rounded-lg shadow overflow-hidden flex flex-col justify-start relative"
  >
    <div
      @click="deleteFavorite(product.variant_id)"
      class="w-10 h-10 p-2 bg-red-200 rounded-lg flex items-center justify-center absolute top-2 right-2 z-10 duration-200 hover:bg-red-300 hover:cursor-pointer"
    >
      <svg
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
    <router-link
      :to="'/catalog/' + product.product_id"
      class="relative before:w-full before:h-full before:absolute before:duration-200 hover:before:bg-sky-50/50"
    >
      <img
        v-if="!product.images.length"
        src="https://placehold.co/600x400"
        :alt="product.name"
        class="w-full aspect-square object-cover"
      />
      <img
        v-else
        :src="product.images[0].image_url"
        :alt="product.name"
        class="w-full aspect-square object-cover"
      />
    </router-link>

    <div class="px-4 py-2">
      <div class="mb-2">
        <span class="py-1 px-3 text-xs rounded-2xl bg-neutral-200 duration-200 select-none">
          {{ product.variant_name }}</span
        >
      </div>
      <router-link
        :to="'/catalog/' + product.product_id"
        class="font-semibold duration-200 hover:text-blue-500 text-md leading-[120%] line-clamp-2"
        >{{ product.name }}</router-link
      >
      <p class="text-gray-600 mt-1">{{ product.price }} ₽</p>
    </div>
  </div>
</template>

<style scoped></style>

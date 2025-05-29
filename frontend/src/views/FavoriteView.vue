<script>
import axios from 'axios'
import { mapMutations } from 'vuex'
import VFavoriteCard from '@/components/ui/VFavoriteCard.vue'
import VButton from '@/components/forms/VButton.vue'
import VBlockLoader from '@/components/loaders/VBlockLoader.vue'
export default {
  name: 'FavoriteView',

  components: {
    VFavoriteCard,
    VButton,
    VBlockLoader,
  },

  props: {},

  data() {
    return {
      products: [],
      isLoading: false,
    }
  },

  computed: {},

  methods: {
    async deleteAll() {
      try {
        await axios.delete('http://localhost:8080/favorites/', {
          headers: {
            Authorization: localStorage.getItem('token'),
          },
        })
        this.products = []
      } catch (error) {
        console.log(error)
      }
    },
    async getFavorites() {
      try {
        this.isLoading = true
        const response = await axios.get('http://localhost:8080/favorites/', {
          headers: {
            Authorization: localStorage.getItem('token'),
          },
        })
        if (response.data.favorites == null) {
          this.products = []
        } else {
          this.products = response.data.favorites
        }
      } catch (error) {
        console.log(error)
      } finally {
        this.isLoading = false
      }
    },
    onCardDeleted(deletedVariantId) {
      this.DECREMENT()
      this.products = this.products.filter((p) => p.variant_id !== deletedVariantId)
    },
    ...mapMutations('favorites', ['INCREMENT', 'DECREMENT']),
  },

  watch: {},

  created() {
    this.getFavorites()
  },
  mounted() {},
  updated() {},
  beforeUnmount() {},
  unmounted() {},
}
</script>

<template>
  <div class="grid gap-4">
    <div v-if="products.length && !isLoading">
      <VButton :variant="'danger'" @click="deleteAll">Удалить все из избранного</VButton>
    </div>
    <div class="grid grid-cols-5 gap-4" v-if="products.length && !isLoading">
      <VFavoriteCard
        @deleted="onCardDeleted"
        v-for="product in products"
        :key="product.id"
        :product="product"
      />
    </div>
    <div class="grid grid-cols-5 gap-4" v-if="isLoading">
      <div
        v-for="i in 10"
        :key="i"
        class="bg-white rounded-lg shadow overflow-hidden flex flex-col justify-start"
      >
        <VBlockLoader class="w-full !rounded-none aspect-square" />

        <div class="px-4 py-2">
          <VBlockLoader class="w-full h-4" />
          <VBlockLoader class="w-2/5 h-4 mt-1" />
        </div>
      </div>
    </div>
    <div
      v-if="!products.length && !isLoading"
      class="w-full p-5 text-2xl flex items-center gap-4 flex-col justify-center bg-neutral-200 rounded-2xl"
    >
      Вы ничего не добавили в избранное
      <div>
        <router-link to="/catalog/">
          <VButton class="text-xs">В каталог</VButton>
        </router-link>
      </div>
    </div>
  </div>
</template>

<style scoped></style>

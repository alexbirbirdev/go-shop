<script>
import axios from 'axios'
import VButton from '@/components/forms/VButton.vue'
// import VSelect from '@/components/forms/VSelect.vue'
import ProductCard from '@/components/ui/VProductCard.vue'
import VBlockLoader from '@/components/loaders/VBlockLoader.vue'
export default {
  name: 'CatalogView',

  components: {
    ProductCard,
    // VSelect,
    VButton,
    VBlockLoader,
  },

  props: {},

  data() {
    return {
      categoryLoading: false,
      categoryCurrent: '',
      categoryBack: '',
      subCategories: [],

      products: [
        { id: 1, name: 'Футболка', price: 1200, image: 'https://placehold.co/400x400' },
        { id: 2, name: 'Кроссовки', price: 5600, image: 'https://placehold.co/400x400' },
      ],
      cart: [],
      selectedSort: {
        value: 'name',
        direction: 'asc',
      },
      sortOptions: [
        { value: 'name_asc', label: 'По имени ↑', direction: 'asc' },
        { value: 'name_desc', label: 'По имени ↓', direction: 'desc' },
        { value: 'price_asc', label: 'По цене ↑', direction: 'asc' },
        { value: 'price_desc', label: 'По цене ↓', direction: 'desc' },
        { value: 'date_asc', label: 'По дате ↑', direction: 'asc' },
        { value: 'date_desc', label: 'По дате ↓', direction: 'desc' },
      ],

      isLoading: true,
    }
  },

  methods: {
    addToCart() {
      
    },
    sortProducts(sortOption) {
      console.log('Сортировка по:', sortOption)
    },

    async getCategories(id) {
      this.categoryLoading = true
      try {
        let response
        if (!id) {
          response = await axios.get('http://localhost:8080/category/')
          this.subCategories = response.data.categories
        } else {
          response = await axios.get('http://localhost:8080/category/?parent_id=' + id)
          this.subCategories = response.data.categories[0].Children
        }
      } catch (error) {
        console.log('Ошибка при запросе категорий: ', error)
      } finally {
        this.categoryLoading = false
      }
    },
    async getParentCategory(id) {
      this.categoryLoading = true
      if (id == null) {
        this.categoryCurrent = ''
        return
      }
      try {
        const response = await axios.get('http://localhost:8080/category/' + id)
        this.categoryCurrent = response.data.category.name
        this.categoryBack = response.data.category.parent_id
      } catch (error) {
        console.log('Ошибка при запросе категорий: ', error)
      } finally {
        this.categoryLoading = false
      }
    },

    async getProducts() {
      try {
        this.isLoading = true
        const response = await axios.get('http://localhost:8080/products/', {
          params: {},
          headers: {
            Authorization: localStorage.getItem('token'),
          },
        })
        this.products = response.data.products
      } catch (error) {
        console.log(error)
      } finally {
        this.isLoading = false
      }
    },
  },
  computed: {
    cartItemCount() {
      return this.cart.length
    },
  },

  watch: {
    selectedSort(newVal) {
      this.sortProducts(newVal)
    },
    '$route.query.category'(newVal) {
      this.categoryBack = ''
      this.getParentCategory(newVal)
      this.getCategories(newVal)
      this.categoryParentName = ''
    },
  },

  created() {
    this.getProducts()
    if (this.$route.query.category) {
      this.getParentCategory(this.$route.query.category)
      this.getCategories(this.$route.query.category)
    } else {
      this.categoryParentName = ''
      this.getCategories(this.$route.query.category)
    }
  },
  mounted() {},
  updated() {},
  beforeUnmount() {},
  unmounted() {},
}
</script>

<template>
  <div class="flex flex-col gap-10">
    <!-- Каталог  -->
    <div class="grid grid-cols-5 gap-5">
      <div class="bg-white rounded-lg py-10 px-5 flex flex-col gap-4">
        <div>
          <div class="text-xs text-blue-600 *:duration-200 *:hover:text-blue-400">
            <RouterLink
              class="mb-2 block"
              v-if="categoryBack"
              :to="'/catalog/?category=' + categoryBack"
              >Назад</RouterLink
            >
            <RouterLink
              class="mb-2 block"
              v-if="categoryBack == null && categoryCurrent != ''"
              to="/catalog/"
              >Назад</RouterLink
            >
          </div>

          <div class="text-2xl leading-[120%]" v-if="categoryCurrent != ''">
            {{ categoryCurrent }}
          </div>
          <div
            class="flex flex-col gap-2"
            :class="subCategories.length > 0 ? 'mt-2' : ''"
            v-if="subCategories.length > 0"
          >
            <!-- {{ subCategories }} -->
            <router-link
              :to="'/catalog/?category=' + sc.id"
              v-for="sc in subCategories"
              :key="sc.id"
              class="leading-[120%] hover:text-blue-600 duration-200"
              >{{ sc.name }}</router-link
            >
          </div>
        </div>
        <div class="">фильтры</div>
        <div class="">btn</div>
      </div>
      <div class="col-span-4 flex flex-col gap-10">
        <div class="">
          <!-- <VSelect v-model="selectedSort" :options="sortOptions" /> -->
        </div>
        <div class="grid grid-cols-4 gap-4" v-if="isLoading">
          <div v-for="i in 8" :key="i" class="bg-white rounded-lg shadow overflow-hidden flex flex-col justify-start">
            <VBlockLoader class="w-full !rounded-none aspect-square" />

            <div class="px-4 py-2">
              <VBlockLoader class="w-full h-4" />
              <VBlockLoader class="w-2/5 h-4 mt-1" />
            </div>
          </div>
        </div>
        <div class="grid grid-cols-4 gap-4" v-if="!isLoading">
          <ProductCard
            v-for="product in products"
            :key="product.id"
            :product="product"
            @add-to-cart="addToCart"
          />
        </div>
        <div class="flex justify-center">
          <VButton>Загрузить еще</VButton>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped></style>

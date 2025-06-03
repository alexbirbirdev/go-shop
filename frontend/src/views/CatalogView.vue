<script>
import axios from 'axios'
import VButton from '@/components/forms/VButton.vue'
import VSelect from '@/components/forms/VSelect.vue'
import ProductCard from '@/components/ui/VProductCard.vue'
import VBlockLoader from '@/components/loaders/VBlockLoader.vue'
export default {
  name: 'CatalogView',

  components: {
    ProductCard,
    VSelect,
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
        value: null,
        label: 'Сортировать по:',
      },
      sortOptions: [
        { value: 'name_asc', label: 'По имени ↑' },
        { value: 'name_desc', label: 'По имени ↓' },
        { value: 'price_asc', label: 'По цене ↑' },
        { value: 'price_desc', label: 'По цене ↓' },
        { value: 'created_asc', label: 'По дате ↑' },
        { value: 'created_desc', label: 'По дате ↓' },
      ],

      isLoading: false,

      filters: {
        price_min: null,
        price_max: null,
      },

      currentPage: null,
      itemsPerPage: 20,
      moreAvailable: true,
      prevAvailable: false,
      
      noMore: false,
    }
  },

  methods: {
    sortProducts(sortOption) {
      const query = {
        ...this.$route.query,
        sort: sortOption.value,
      }

      // Убираем null/пустые параметры из URL
      Object.keys(query).forEach((key) => {
        if (query[key] === null || query[key] === '') {
          delete query[key]
        }
      })

      this.$router.push({ path: '/catalog', query })
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
      if (this.isLoading) return
      try {
        this.isLoading = true
        if (!this.$route.query.page || this.$route.query.page < 2) {
          this.prevAvailable = false
        } else {
          this.prevAvailable = true
        }
        const response = await axios.get('http://localhost:8080/products/', {
          params: {
            price_min: this.$route.query.price_min,
            price_max: this.$route.query.price_max,
            category: this.$route.query.category || null,
            sort: this.$route.query.sort || null,
            page: this.$route.query.page || 1,
            limit: this.itemsPerPage,
          },
          headers: {
            Authorization: localStorage.getItem('token'),
          },
        })
        this.products = response.data.products
        this.moreAvailable = response.data.products.length === this.itemsPerPage

        if (this.products.length == 0) {
          this.noMore = true
          this.moreAvailable = false
          this.$router.replace({
            query: {
              ...this.$route.query,
              page: this.currentPage - 1,
            },
          })
        }
        if (this.moreAvailable) {
          if (this.currentPage !== parseInt(this.$route.query.page)) {
            this.$router.replace({
              query: {
                ...this.$route.query,
                page: this.currentPage,
              },
            })
          }
        }
      } catch (error) {
        console.log(error)
      } finally {
        this.isLoading = false
      }
    },
    loadMoreProducts() {
      if (this.moreAvailable && !this.isLoading) {
        const nextPage = this.currentPage + 1
        this.$router.push({ query: { ...this.$route.query, page: nextPage } })
      }
    },
    loadPreviousProducts() {
      if (this.currentPage > 1 && this.prevAvailable && !this.isLoading) {
        const prevPage = this.currentPage - 1
        this.noMore = false
        this.$router.push({ query: { ...this.$route.query, page: prevPage } })
      }
    },
    handleFilters() {
      this.noMore = false
      const query = {
        ...this.$route.query,
        price_min: this.filters.price_min,
        price_max: this.filters.price_max,
      }

      // Убираем null/пустые параметры из URL
      Object.keys(query).forEach((key) => {
        if (query[key] === null || query[key] === '') {
          delete query[key]
        }
      })

      this.$router.push({ path: '/catalog', query })
    },
    clearFilters() {
      this.filters.price_min = null
      this.filters.price_max = null

      this.noMore = false

      const query = { ...this.$route.query }
      delete query.price_min
      delete query.price_max

      this.$router.push({ path: '/catalog', query })
    },
  },
  computed: {},

  watch: {
    selectedSort(newVal) {
      this.sortProducts(newVal)
    },
    '$route.query.category'(newVal) {
      this.filters.price_max = null
      this.filters.price_min = null
      this.categoryBack = ''
      this.getParentCategory(newVal)
      this.getCategories(newVal)
      this.categoryParentName = ''
      this.noMore = false
      this.getProducts()
    },

    '$route.query.price_min': 'getProducts',
    '$route.query.price_max': 'getProducts',
    '$route.query.sort': {
      handler() {
        this.selectedSort = this.sortOptions.find(
          (opt) => opt.value === this.$route.query.sort,
        ) || { value: null, label: 'Сортировать по:' }

        this.getProducts()
      },
      immediate: true,
    },
    '$route.query.page': {
      handler(newVal) {
        const page = parseInt(newVal)
        this.currentPage = isNaN(page) || page < 1 ? 1 : page
        this.getProducts()
      },
      immediate: true,
    },
  },

  created() {
    const query = this.$route.query
    const page = parseInt(query.page)
    this.currentPage = isNaN(page) || page < 1 ? 1 : page
    this.filters.price_min = query.price_min || null
    this.filters.price_max = query.price_max || null

    // this.getProducts()

    if (query.category) {
      this.getParentCategory(query.category)
      this.getCategories(query.category)
    } else {
      this.categoryParentName = ''
      this.getCategories()
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
          <div class="text-lg font-bold mb-3">Категория</div>
          <div class="text-xs text-blue-600 *:duration-200 *:hover:text-blue-400">
            <RouterLink class="block" v-if="categoryBack" :to="'/catalog?category=' + categoryBack"
              >Назад</RouterLink
            >
            <RouterLink
              class="block"
              v-if="categoryBack == null && categoryCurrent != ''"
              to="/catalog"
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
              :to="'/catalog?category=' + sc.id"
              v-for="sc in subCategories"
              :key="sc.id"
              class="leading-[120%] hover:text-blue-600 duration-200"
              >{{ sc.name }}</router-link
            >
          </div>
        </div>
        <form @submit.prevent.stop="handleFilters" class="grid gap-3 mt-3">
          <div class="text-lg font-bold">Фильтры</div>
          <div>
            <label for="price_min" class="block text-sm font-medium text-gray-700">Цена от</label>
            <input
              id="price_min"
              v-model="filters.price_min"
              type="text"
              class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
            />
          </div>
          <div>
            <label for="price_max" class="block text-sm font-medium text-gray-700">Цена до</label>
            <input
              id="price_max"
              v-model="filters.price_max"
              type="text"
              class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
            />
          </div>
          <VButton class="w-full">Применить</VButton>
          <VButton @click="clearFilters" class="bg-neutral-200 w-full">Сбросить</VButton>
        </form>
      </div>
      <div class="col-span-4 flex flex-col gap-5">
        <div class="">
          <VSelect v-model="selectedSort" :options="sortOptions" />
        </div>

        <div class="grid grid-cols-4 gap-4" v-if="!isLoading && products.length">
          <ProductCard v-for="product in products" :key="product.id" :product="product" />
        </div>
        <div
          v-if="!products.length && !isLoading && $route.query.page == 1"
          class="w-full p-5 text-2xl flex items-center gap-4 flex-col justify-center bg-neutral-200 rounded-2xl text-center"
        >
          Товаров нет :( <br />
          <small>Попробуйте сбросить фильтры или выбрать другую категорию</small>
        </div>
        <div class="grid grid-cols-4 gap-4" v-if="isLoading">
          <div
            v-for="i in 8"
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
        <div class="flex justify-center items-center gap-4">
          <div v-if="prevAvailable && !isLoading" class="flex justify-center mt-10">
            <VButton :disabled="isLoading" @click="loadPreviousProducts">
              <span v-if="isLoading" class="animate-spin"
                ><svg
                  class="w-4 h-4 *:fill-white"
                  viewBox="0 0 24 24"
                  fill="none"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    d="M12 2.25C17.3848 2.25 21.75 6.61522 21.75 12C21.75 17.3848 17.3848 21.75 12 21.75C6.61522 21.75 2.25 17.3848 2.25 12C2.25 10.4448 2.61447 8.97237 3.26367 7.66602C3.44804 7.29514 3.89863 7.1438 4.26953 7.32812C4.6404 7.51249 4.79174 7.96308 4.60742 8.33398C4.05901 9.43754 3.75 10.6816 3.75 12C3.75 16.5563 7.44365 20.25 12 20.25C16.5563 20.25 20.25 16.5563 20.25 12C20.25 7.44365 16.5563 3.75 12 3.75C11.5858 3.75 11.25 3.41421 11.25 3C11.25 2.58579 11.5858 2.25 12 2.25Z"
                    fill="black"
                  />
                </svg> </span
              ><span v-else
                ><svg
                  class="w-4 h-4 *:fill-white"
                  viewBox="0 0 24 24"
                  fill="none"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    d="M13.4697 5.46967C13.7626 5.17678 14.2373 5.17678 14.5302 5.46967C14.8231 5.76256 14.8231 6.23732 14.5302 6.53022L9.06049 11.9999L14.5302 17.4697L14.582 17.5263C14.8223 17.8209 14.8048 18.2556 14.5302 18.5302C14.2556 18.8048 13.8209 18.8223 13.5263 18.582L13.4697 18.5302L7.46967 12.5302C7.17678 12.2373 7.17678 11.7626 7.46967 11.4697L13.4697 5.46967Z"
                    fill="black"
                  />
                </svg>
              </span>
            </VButton>
          </div>
          <div v-if="moreAvailable && !noMore" class="flex justify-center mt-10">
            <VButton :disabled="isLoading || !moreAvailable" @click="loadMoreProducts">
              <span v-if="isLoading" class="animate-spin"
                ><svg
                  class="w-4 h-4 *:fill-white"
                  viewBox="0 0 24 24"
                  fill="none"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    d="M12 2.25C17.3848 2.25 21.75 6.61522 21.75 12C21.75 17.3848 17.3848 21.75 12 21.75C6.61522 21.75 2.25 17.3848 2.25 12C2.25 10.4448 2.61447 8.97237 3.26367 7.66602C3.44804 7.29514 3.89863 7.1438 4.26953 7.32812C4.6404 7.51249 4.79174 7.96308 4.60742 8.33398C4.05901 9.43754 3.75 10.6816 3.75 12C3.75 16.5563 7.44365 20.25 12 20.25C16.5563 20.25 20.25 16.5563 20.25 12C20.25 7.44365 16.5563 3.75 12 3.75C11.5858 3.75 11.25 3.41421 11.25 3C11.25 2.58579 11.5858 2.25 12 2.25Z"
                    fill="black"
                  />
                </svg> </span
              ><span v-else>
                <svg
                  class="w-4 h-4 *:fill-white"
                  viewBox="0 0 24 24"
                  fill="none"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    d="M9.46984 5.46984C9.74445 5.19524 10.1792 5.17778 10.4737 5.41809L10.5304 5.46984L16.5304 11.4698L16.5821 11.5265C16.8065 11.8014 16.8065 12.1988 16.5821 12.4737L16.5304 12.5304L10.5304 18.5304C10.2375 18.8233 9.76274 18.8233 9.46984 18.5304C9.17695 18.2375 9.17695 17.7627 9.46984 17.4698L14.9396 12.0001L9.46984 6.53039L9.41809 6.47375C9.17778 6.17917 9.19524 5.74445 9.46984 5.46984Z"
                    fill="black"
                  />
                </svg>
              </span>
            </VButton>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped></style>

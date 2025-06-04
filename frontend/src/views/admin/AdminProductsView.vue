<script>
import VButton from '@/components/forms/VButton.vue'
import VBlockLoader from '@/components/loaders/VBlockLoader.vue'
import VPopupCreateProduct from '@/components/popups/VPopupCreateProduct.vue'
import axios from 'axios'
export default {
  name: 'AdminProductsView',

  components: {
    VButton,
    VPopupCreateProduct,
    VBlockLoader,
  },

  props: {},

  data() {
    return {
      products: [],
      tableHeaders: [
        'Действия',
        'ID',
        'Фото',
        'Название',
        'Активность',
        'Категория',
        'Цена',
        'Наличие',
      ],
      isLoading: true,
      currentPage: 1,
      itemsPerPage: 20,
      moreAvailable: false,
      prevAvailable: false,
      showCreateProduct: false,
      isDeleting: false,
    }
  },

  computed: {},

  methods: {
    openCreateProductPopup() {
      this.showCreateProduct = true
    },
    async getProducts(reset = false) {
      try {
        if (reset) {
          this.currentPage = 1
          this.products = []
        }
        if (this.currentPage < 2) {
          this.prevAvailable = false
        } else {
          this.prevAvailable = true
        }
        this.isLoading = true
        const response = await axios.get('/api/admin/products/', {
          params: {
            page: this.currentPage,
            limit: this.itemsPerPage,
          },
          headers: {
            Authorization: localStorage.getItem('token'),
          },
        })
        this.products = response.data.products
        this.$router.replace({
          query: {
            page: this.currentPage,
            limit: this.itemsPerPage,
          },
        })
        if (response.data.products.length === this.itemsPerPage) {
          this.moreAvailable = true
        } else {
          this.moreAvailable = false
          this.currentPage--
        }
      } catch (error) {
        console.log(error)
      } finally {
        this.isLoading = false
      }
    },
    loadPreviousProducts() {
      if (this.currentPage > 1) {
        this.currentPage--
        this.getProducts()
      }
    },
    loadMoreProducts() {
      if (this.moreAvailable && !this.isLoading) {
        this.currentPage++
        this.getProducts()
      }
    },
    async deleteProduct(id) {
      if (this.isDeleting) return
      try {
        this.isDeleting = true
        await axios.delete('/api/admin/products/' + id, {
          headers: {
            Authorization: localStorage.getItem('token'),
          },
        })
        this.products = this.products.filter((product) => product.id !== id)
      } catch (error) {
        console.log(error)
      } finally {
        this.isDeleting = false
      }
    },
  },

  watch: {},

  created() {
    this.getProducts()
  },
  mounted() {},
  updated() {},
  beforeUnmount() {},
  unmounted() {},
}
</script>

<template>
  <div class="">
    <VPopupCreateProduct v-model="showCreateProduct" />
    <div class="mb-10" v-if="!isLoading">
      <VButton @click="openCreateProductPopup">Добавить</VButton>
    </div>

    <div v-if="isLoading && !products.length">
      <table class="w-full">
        <thead class="">
          <tr class="text-left">
            <th v-for="i in tableHeaders" :key="i" class="">
              <div class="bg-blue-100 py-4 px-2">
                <VBlockLoader class="w-17 h-5" />
              </div>
            </th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="i in 3" :key="i">
            <td class="">
              <div class="flex gap-2">
                <VBlockLoader class="w-7 h-7" />
                <!-- <VBlockLoader class="w-7 h-7" /> -->
                <VBlockLoader class="w-7 h-7" />
              </div>
            </td>
            <td>
              <div>
                <VBlockLoader class="w-10 h-4" />
              </div>
            </td>
            <td>
              <div class="">
                <VBlockLoader class="w-20 h-20" />
              </div>
            </td>
            <td>
              <div>
                <VBlockLoader class="w-25 h-4" />
              </div>
            </td>
            <td>
              <div class="">
                <VBlockLoader class="w-20 h-6" />
              </div>
            </td>
            <td>
              <div>
                <VBlockLoader class="w-10 h-4" />
              </div>
            </td>
            <td>
              <div>
                <VBlockLoader class="w-15 h-4" />
              </div>
            </td>
            <td>
              <div>
                <VBlockLoader class="w-10 h-4" />
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <table class="w-full" v-if="products.length">
      <thead class="">
        <tr class="text-left">
          <th v-for="header in tableHeaders" :key="header" class="">
            <div class="bg-blue-100 py-4 px-2">{{ header }}</div>
          </th>
        </tr>
      </thead>
      <tbody v-if="!isLoading">
        <tr v-for="product in products" :key="product.ID">
          <td class="">
            <div class="flex gap-2">
              <VButton :variant="'warning'" class="!p-2">
                <svg
                  class="w-5 *:fill-yellow-200"
                  viewBox="0 0 24 24"
                  fill="none"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    d="M12 2.25034C12.4142 2.25034 12.75 2.58612 12.75 3.00034C12.7498 3.41441 12.4141 3.75034 12 3.75034H11C9.09324 3.75034 7.73859 3.75183 6.71094 3.88999C5.8306 4.00835 5.27638 4.21718 4.86816 4.55112L4.70215 4.70249C4.27903 5.12568 4.0249 5.70531 3.88965 6.71128C3.75154 7.73889 3.75 9.09376 3.75 11.0003V13.0003C3.75 14.9069 3.7515 16.2618 3.88965 17.2894C4.02492 18.2953 4.27898 18.875 4.70215 19.2982L4.86816 19.4496C5.27638 19.7834 5.83068 19.9923 6.71094 20.1107C7.73859 20.2488 9.09327 20.2503 11 20.2503H13C14.9067 20.2503 16.2614 20.2488 17.2891 20.1107C18.295 19.9754 18.8747 19.7213 19.2979 19.2982L19.4492 19.1322C19.7831 18.724 19.992 18.1696 20.1104 17.2894C20.2485 16.2618 20.25 14.9069 20.25 13.0003V12.0003C20.25 11.5861 20.5858 11.2503 21 11.2503C21.4142 11.2503 21.75 11.5861 21.75 12.0003V13.0003C21.75 14.8646 21.7514 16.3388 21.5967 17.4896C21.4489 18.5884 21.1477 19.4911 20.4941 20.2162L20.3584 20.3587C19.6101 21.1069 18.6615 21.4394 17.4893 21.597C16.3384 21.7517 14.8644 21.7503 13 21.7503H11C9.13561 21.7503 7.66158 21.7517 6.51074 21.597C5.41183 21.4492 4.50926 21.148 3.78418 20.4945L3.6416 20.3587C2.89336 19.6105 2.56098 18.6618 2.40332 17.4896C2.2486 16.3388 2.25 14.8646 2.25 13.0003V11.0003C2.25 9.13608 2.24863 7.66188 2.40332 6.51108C2.56095 5.33885 2.8934 4.39022 3.6416 3.64194L3.78418 3.5062C4.50928 2.85261 5.41178 2.55144 6.51074 2.40366C7.66159 2.24893 9.13559 2.25034 11 2.25034H12ZM17.3535 2.62631C18.3976 1.84438 19.8664 1.90052 20.8496 2.79428L21.0303 2.97006L21.2061 3.15073L21.374 3.35385C22.1036 4.32821 22.1037 5.67252 21.374 6.64682L21.2061 6.84995L21.0303 7.03061L17.0977 10.9642C16.441 11.6208 15.9989 12.0684 15.5088 12.391L15.2949 12.5218C14.7128 12.8513 14.0615 13.0077 13.0322 13.265L11.1816 13.7279C10.9262 13.7916 10.6559 13.7167 10.4697 13.5306C10.2836 13.3445 10.2088 13.0741 10.2725 12.8187L10.7354 10.9681C10.9927 9.93886 11.1491 9.28749 11.4785 8.70542L11.6094 8.49155C11.932 8.00141 12.3795 7.55931 13.0361 6.90268L16.9697 2.97006L17.1504 2.79428L17.3535 2.62631ZM14.0977 7.9642C13.3919 8.66998 13.0729 8.99481 12.8662 9.30893L12.7842 9.4437C12.5781 9.80773 12.4671 10.2257 12.1904 11.3324L12.0303 11.9691L12.668 11.8099L13.3789 11.6292C13.9802 11.4717 14.2836 11.3707 14.5566 11.2162L14.6914 11.1341C15.0055 10.9275 15.3304 10.6085 16.0361 9.90268L18.3926 7.54624C17.5069 7.20046 16.7987 6.49258 16.4531 5.60678L14.0977 7.9642ZM19.8408 3.90366C19.3938 3.49735 18.7255 3.47182 18.251 3.82749L18.1592 3.90366L18.0303 4.03061L17.7227 4.33725C17.6984 5.41363 18.5867 6.30021 19.6631 6.27573L19.9697 5.97006L20.0967 5.84116L20.1729 5.74936C20.5048 5.3064 20.5047 4.69433 20.1729 4.25131L20.0967 4.15952L19.9697 4.03061L19.8408 3.90366Z"
                    fill="black"
                  />
                </svg>
              </VButton>
              <!-- <VButton class="!p-2">
                <svg
                  class="w-5 *:fill-blue-100"
                  viewBox="0 0 24 24"
                  fill="none"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    d="M12 4.25C15.6923 4.25 18.373 5.57392 20.1338 7.19922C21.0095 8.00759 21.6562 8.88899 22.0869 9.71582C22.5106 10.5293 22.75 11.3384 22.75 12C22.75 12.6616 22.5106 13.4707 22.0869 14.2842C21.6562 15.111 21.0095 15.9924 20.1338 16.8008C18.373 18.4261 15.6923 19.75 12 19.75C8.30768 19.75 5.62695 18.4261 3.86621 16.8008C2.99052 15.9924 2.34376 15.111 1.91309 14.2842C1.48938 13.4707 1.25 12.6616 1.25 12C1.25 11.3384 1.48938 10.5293 1.91309 9.71582C2.34376 8.88899 2.99052 8.00759 3.86621 7.19922C5.62695 5.57392 8.30768 4.25 12 4.25ZM12 5.75C8.69247 5.75 6.37305 6.92616 4.88379 8.30078C4.13467 8.99228 3.5937 9.73614 3.24316 10.4092C2.88572 11.0955 2.75 11.6616 2.75 12C2.75 12.3384 2.88572 12.9045 3.24316 13.5908C3.5937 14.2639 4.13467 15.0077 4.88379 15.6992C6.37305 17.0738 8.69247 18.25 12 18.25C15.3075 18.25 17.6269 17.0738 19.1162 15.6992C19.8653 15.0077 20.4063 14.2639 20.7568 13.5908C21.1143 12.9045 21.25 12.3384 21.25 12C21.25 11.6616 21.1143 11.0955 20.7568 10.4092C20.4063 9.73614 19.8653 8.99228 19.1162 8.30078C17.6269 6.92616 15.3075 5.75 12 5.75ZM12 8.25C14.071 8.25005 15.75 9.92898 15.75 12C15.75 14.071 14.071 15.7499 12 15.75C9.92893 15.75 8.25 14.0711 8.25 12C8.25001 9.92894 9.92894 8.25 12 8.25ZM12 9.75C10.7574 9.75 9.75001 10.7574 9.75 12C9.75 13.2426 10.7574 14.25 12 14.25C13.2426 14.2499 14.25 13.2426 14.25 12C14.25 10.7574 13.2426 9.75005 12 9.75Z"
                    fill="black"
                  />
                </svg>
              </VButton> -->
              <VButton @click="deleteProduct(product.id)" :variant="'danger'" class="!p-2">
                <svg
                  class="w-5 *:fill-red-100"
                  viewBox="0 0 24 24"
                  fill="none"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    d="M13.5586 2.25C13.8891 2.25 14.2 2.24876 14.459 2.27832C14.7371 2.31008 15.0298 2.38336 15.3027 2.58008L15.3994 2.65625C15.6151 2.84023 15.7521 3.06291 15.8555 3.28516C15.9654 3.52151 16.0625 3.81724 16.167 4.13086L16.541 5.25H21C21.4142 5.25 21.75 5.58579 21.75 6C21.75 6.41421 21.4142 6.75 21 6.75H19.7012L19.1221 15.4492C19.0349 16.7562 18.9659 17.8097 18.8018 18.6338C18.6544 19.3733 18.4163 20.006 17.9629 20.5234L17.7549 20.7383C17.1634 21.2916 16.4403 21.5308 15.5859 21.6426C14.7527 21.7516 13.6967 21.75 12.3867 21.75H11.6133C10.3033 21.75 9.2473 21.7516 8.41406 21.6426C7.66641 21.5447 7.01906 21.3496 6.47266 20.9316L6.24512 20.7383C5.65357 20.1849 5.36665 19.479 5.19824 18.6338C5.07509 18.0156 5.00617 17.2682 4.94238 16.3818L4.87793 15.4492L4.29883 6.75H3C2.58579 6.75 2.25 6.41421 2.25 6C2.25 5.58579 2.58579 5.25 3 5.25H7.45898L7.83301 4.13086L7.98633 3.67871C8.03745 3.53613 8.08955 3.40333 8.14453 3.28516C8.26261 3.03133 8.42436 2.77679 8.69727 2.58008L8.80078 2.51172C9.04332 2.36556 9.29772 2.3061 9.54102 2.27832C9.79998 2.24876 10.1109 2.25 10.4414 2.25H13.5586ZM6.375 15.3496C6.46484 16.6973 6.52839 17.6353 6.66895 18.3408C6.80526 19.025 6.99626 19.3868 7.26953 19.6426L7.37695 19.7344C7.64122 19.9389 8.00362 20.076 8.60938 20.1553C9.32273 20.2486 10.2627 20.25 11.6133 20.25H12.3867C13.7373 20.25 14.6773 20.2486 15.3906 20.1553C16.0826 20.0648 16.4571 19.8984 16.7305 19.6426L16.8291 19.542C17.0506 19.2919 17.2117 18.9397 17.3311 18.3408C17.4716 17.6353 17.5352 16.6973 17.625 15.3496L18.1982 6.75H5.80176L6.375 15.3496ZM10 9.25C10.4142 9.25 10.75 9.58579 10.75 10V17C10.75 17.4142 10.4142 17.75 10 17.75C9.58579 17.75 9.25 17.4142 9.25 17V10C9.25 9.58579 9.58579 9.25 10 9.25ZM14 9.25C14.4142 9.25 14.75 9.58579 14.75 10V14C14.75 14.4142 14.4142 14.75 14 14.75C13.5858 14.75 13.25 14.4142 13.25 14V10C13.25 9.58579 13.5858 9.25 14 9.25ZM10.4414 3.75C10.0763 3.75 9.86526 3.75094 9.71094 3.76855C9.64144 3.7765 9.60514 3.78629 9.58789 3.79199C9.58004 3.7946 9.57596 3.79649 9.5752 3.79688H9.57422L9.57324 3.79785C9.57245 3.79865 9.57021 3.80219 9.56543 3.80859C9.55454 3.82321 9.53437 3.85458 9.50488 3.91797C9.47208 3.98849 9.43807 4.07352 9.39844 4.18457L9.25586 4.60449L9.04102 5.25H14.959L14.7441 4.60449C14.6287 4.25825 14.5606 4.05878 14.4951 3.91797C14.4656 3.85458 14.4455 3.82321 14.4346 3.80859C14.4298 3.80219 14.4276 3.79865 14.4268 3.79785L14.4258 3.79688H14.4248C14.424 3.79649 14.42 3.7946 14.4121 3.79199C14.3949 3.78629 14.3586 3.7765 14.2891 3.76855C14.1347 3.75094 13.9237 3.75 13.5586 3.75H10.4414Z"
                    fill="black"
                  />
                </svg>
              </VButton>
            </div>
          </td>
          <td>
            <div>
              {{ product.id }}
            </div>
          </td>
          <td>
            <div class="">
              <div
                v-if="!product.images.length"
                class="bg-cover bg-center w-20 aspect-square rounded-xl overflow-hidden shadow"
                style="background-image: url(https://placehold.co/600x400)"
              ></div>
              <div
                class="bg-cover bg-center w-20 aspect-square rounded-xl overflow-hidden shadow"
                :style="'background-image: url(' + product.images[0].image_url + ')'"
                v-else
              ></div>
            </div>
          </td>
          <td>
            <div>
              {{ product.name }}
            </div>
          </td>
          <td>
            <div class="">
              <span v-if="!product.is_active" class="bg-red-100 py-2 px-4 rounded-xl text-red-900">
                Не активен
              </span>
              <span
                v-if="product.is_active"
                class="bg-green-100 py-2 px-4 rounded-xl text-green-900"
              >
                Активен
              </span>
            </div>
          </td>
          <td>
            <div>
              {{ product.category }}
            </div>
          </td>
          <td>
            <div>
              {{ product.price }}
            </div>
          </td>
          <td>
            <div>
              {{ product.stock }}
            </div>
          </td>
        </tr>
      </tbody>
    </table>
    <div
      class="flex items-center justify-center bg-neutral-200 text-neutral-950 p-10 rounded-xl"
      v-if="!products.length && !isLoading"
    >
      <div class="font-bold">Товаров нет</div>
    </div>
    <div class="flex justify-center items-center gap-4">
      <div v-if="prevAvailable" class="flex justify-center mt-10">
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
      <div v-if="moreAvailable" class="flex justify-center mt-10">
        <VButton :disabled="isLoading" @click="loadMoreProducts">
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
</template>

<style scoped>
th {
  border-bottom: 4px solid transparent;
}
th:nth-child(1) {
  width: 1px;
}
th:nth-child(2) {
  width: 1px;
  text-align: center;
}
th:nth-child(3) {
  width: 100px;
}
th:nth-child(4) {
  width: 100%;
  text-align: left;
}
th:nth-child(5) {
  width: 100px;
}
th:nth-child(6) {
  /* width: 100px; */
  text-align: center;
}
th:nth-child(7) {
  /* width: 100px; */
  text-align: center;
  white-space: nowrap;
}
th {
  text-align: center;
}
th:nth-child(8) {
  text-align: center;
  /* width: 100px; */
}
td:nth-child(5) div,
td:nth-child(6) div,
td:nth-child(7) div,
td:nth-child(8) div {
  white-space: nowrap;
  display: flex;
  justify-content: center;
}

td:nth-child(6) div {
  display: flex;
  justify-content: center;
}
th:first-child div {
  border-top-left-radius: 8px;
  border-bottom-left-radius: 8px;
}
th:last-child div {
  border-top-right-radius: 8px;
  border-bottom-right-radius: 8px;
}
td:first-child div {
  border-top-left-radius: 8px;
  border-bottom-left-radius: 8px;
}
td:last-child div {
  border-top-right-radius: 8px;
  border-bottom-right-radius: 8px;
}
td div {
  padding: 8px;
  transition: 0.2s;
  display: flex;
  align-items: center;
  max-width: 100%;
  overflow: hidden;
  margin: auto;
}
tr:nth-child(2n + 2):not(thead tr) > td > div {
  background: oklch(97% 0 0);
}
tr:hover:not(thead tr) > td > div {
  background: oklch(92.2% 0 0);
}
table,
tr,
td {
  height: 100%;
}
td > div {
  height: 100%;
}
</style>

<script>
import axios from 'axios'
import VButton from '@/components/forms/VButton.vue'

import VBlockLoader from '@/components/loaders/VBlockLoader.vue'

export default {
  name: 'AdminCategoriesView',

  components: {
    VButton,

    VBlockLoader,
  },

  data() {
    return {
      categories: [],
      currentPage: 1,
      itemsPerPage: 100,
      categoryLoading: false,
      isDeleting: false,
      isUpdating: false,

      loadingCategories: true,
      categoryCurrent: {
        value: null,
        label: 'Выберите категорию',
      },
      addNewCategory: '',
    }
  },

  computed: {
    parsedCategories() {
      return this.categories.map((category) => {
        const parent = this.categories.find((c) => c.ID === category.parent_id)
        return {
          ...category,
          parentName: parent ? parent.name : '—',
        }
      })
    },
  },

  methods: {
    async getCategories() {
      this.categoryLoading = true
      try {
        const response = await axios.get('/api/admin/category/all', {
          params: {
            page: this.currentPage,
            limit: this.itemsPerPage,
            sort: 'created_desc',
          },
          headers: {
            Authorization: localStorage.getItem('token'),
          },
        })
        this.categories = response.data.categories
      } catch (error) {
        console.log('Ошибка при запросе категорий: ', error)
      } finally {
        this.categoryLoading = false
      }
    },

    async updateParent(categoryId, newParentId) {
      if (this.isUpdating) return
      try {
        this.isUpdating = true
        await axios.put(
          `/api/admin/category/${categoryId}`,
          { parent_id: newParentId },
          {
            headers: {
              Authorization: localStorage.getItem('token'),
            },
          },
        )
        this.getCategories()
      } catch (error) {
        console.error(`Ошибка при обновлении родителя категории ${categoryId}:`, error)
      } finally {
        this.isUpdating = false
      }
    },
    async deleteCat(id) {
      if (this.isDeleting) return
      try {
        this.isDeleting = true
        await axios.delete('/api/admin/category/' + id, {
          headers: {
            Authorization: localStorage.getItem('token'),
          },
        })
        this.categories = this.categories.filter((cat) => cat.ID !== id)
      } catch (error) {
        console.log(error)
      } finally {
        this.isDeleting = false
      }
    },
    async updateCategoryName(categoryId, newName) {
      if (this.isUpdating) return
      try {
        this.isUpdating = true
        await axios.put(
          `/api/admin/category/${categoryId}`,
          { name: newName },
          {
            headers: {
              Authorization: localStorage.getItem('token'),
            },
          },
        )
        // Можно обновить список или просто ничего не делать,
        // если изменения уже применены локально
      } catch (error) {
        console.error(`Ошибка при обновлении названия категории ${categoryId}:`, error)
      } finally {
        this.isUpdating = false
      }
    },

    async addCategory() {
      try {
        this.loadingCategories = true
        await axios.post(
          '/api/admin/category/',
          {
            name: this.addNewCategory,
          },
          {
            headers: {
              Authorization: localStorage.getItem('token'),
            },
          },
        )
        this.addNewCategory = ''
        // this.getCategories()
      } catch (error) {
        console.log(error)
      } finally {
        this.loadingCategories = false
      }
    },
  },

  created() {
    this.getCategories()
  },
}
</script>

<template>
  <div class="max-w-xl mx-auto">
    <div class="mb-4">
      <form @submit.prevent="addCategory" class="mt-4">
        <div>
          <label for="addNewCategory" class="block text-sm font-medium text-gray-700"
            >Не нашли категорию? Добавьте</label
          >
          <div class="flex items-end gap-2">
            <input
              v-if="!categoryLoading"
              id="addNewCategory"
              v-model="addNewCategory"
              type="text"
              required
              class="mt-1 flex-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
            />
            <VBlockLoader v-else class="flex-1 h-[42px] mt-1" />
            <VButton :disabled="categoryLoading" class="!p-2 w-[42px] h-[42px]">
              <svg
                class="w-7 *:fill-blue-100"
                viewBox="0 0 24 24"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  d="M12 5.25C12.4142 5.25 12.75 5.58579 12.75 6V11.25H18C18.4142 11.25 18.75 11.5858 18.75 12C18.75 12.4142 18.4142 12.75 18 12.75H12.75V18C12.75 18.4142 12.4142 18.75 12 18.75C11.5858 18.75 11.25 18.4142 11.25 18V12.75H6C5.58579 12.75 5.25 12.4142 5.25 12C5.25 11.5858 5.58579 11.25 6 11.25H11.25V6C11.25 5.58579 11.5858 5.25 12 5.25Z"
                  fill="black"
                />
              </svg>
            </VButton>
          </div>
        </div>
      </form>
    </div>
    <div
      class="p-2 rounded-xl flex gap-4 items-center justify-between odd:bg-neutral-100"
      v-for="cat in parsedCategories"
      :key="cat.ID"
    >
      <div class="">
        <VButton @click="deleteCat(cat.ID)" :variant="'danger'" class="!p-2">
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
      <div class="w-50 flex-1">
        <input
          type="text"
          v-model="cat.name"
          class="border px-2 py-1 rounded w-full"
          :placeholder="'(без названия)'"
          @blur="updateCategoryName(cat.ID, cat.name)"
          @keyup.enter="$event.target.blur()"
        />
      </div>

      <div>
        <select
          class="border px-2 py-1 rounded w-40 overflow-hidden *:w-full"
          v-model="cat.parent_id"
          @change="updateParent(cat.ID, cat.parent_id)"
        >
          <option :value="null" class="font-bold">Без родителя</option>
          <option
            v-for="option in categories"
            :key="option.ID"
            :value="option.ID"
            :disabled="option.ID === cat.ID"
          >
            {{ option.name || '(без названия)' }}
          </option>
        </select>
      </div>
    </div>
  </div>
</template>

<style scoped></style>

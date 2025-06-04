<script>
import axios from 'axios'
import VButton from '../forms/VButton.vue'
import VSelect from '../forms/VSelect.vue'
import VFlashMessage from '../forms/VFlashMessage.vue'
import VBlockLoader from '../loaders/VBlockLoader.vue'

export default {
  name: 'VPopupCreateProduct',

  components: { VButton, VSelect, VFlashMessage, VBlockLoader },

  props: {
    modelValue: {
      type: Boolean,
      default: false,
    },
  },

  data() {
    return {
      form: {
        name: '',
        description: '',
        category: '',
      },
      sortOptions: [],
      loadingCategories: true,
      categoryCurrent: {
        value: null,
        label: 'Выберите категорию',
      },
      addNewCategory: '',
      flashMessage: null,

      currentStep: 0,
      steps: [{ title: 'Основное' }, { title: 'Фото' }, { title: 'Варианты' }],
      uploadedPhotos: [],
      formImages: [],
      variants: [{ name: '', price: 0, stock: 0 }],
      isSubmitting: false,
    }
  },

  computed: {},

  methods: {
    clearForm() {
      this.form = {
        name: '',
        description: '',
        category: '',
      }
      this.categoryCurrent = {
        value: null,
        label: 'Выберите категорию',
      }
      this.addNewCategory = ''
      this.uploadedPhotos = []
      this.variants = [{ name: '', price: 0, stock: 0 }]
      this.currentStep = 0
      this.flashMessage = null
      this.isSubmitting = false
    },
    nextStep() {
      if (this.currentStep < this.steps.length - 1) {
        this.currentStep++
      }
    },
    prevStep() {
      if (this.currentStep > 0) {
        this.currentStep--
      }
    },
    handleFileUpload(e) {
      const files = Array.from(e.target.files)
      files.forEach((file) => {
        const reader = new FileReader()
        reader.onload = (e) => {
          this.uploadedPhotos.push({
            file,
            preview: e.target.result,
          })
        }
        reader.readAsDataURL(file)
      })
    },
    removePhoto(index) {
      this.uploadedPhotos.splice(index, 1)
    },
    addVariant() {
      this.variants.push({ name: '', price: 0, stock: 0 })
    },
    removeVariant(index) {
      if (this.variants.length > 1) {
        this.variants.splice(index, 1)
      }
    },
    async uploadPhotos(productId) {
      try {
        const formData = new FormData()
        this.uploadedPhotos.forEach((photo) => {
          formData.append('images', photo.file)
        })
        console.log('Форма с фото:', formData)

        const response = await axios.post(
          `/api/admin/products/${productId}/images`,
          formData,
          {
            headers: {
              Authorization: localStorage.getItem('token'),
              'Content-Type': 'multipart/form-data',
            },
          },
        )

        console.log(response.data)
      } catch (error) {
        console.error('Ошибка загрузки фото:', error)
        throw error
      }
    },

    async submitAll() {
      try {
        this.isSubmitting = true

        // 1. Сначала создаем основной товар
        const productResponse = await axios.post(
          '/api/admin/products/',
          {
            name: this.form.name,
            description: this.form.description,
            category_id: this.form.category,
          },
          {
            headers: {
              Authorization: localStorage.getItem('token'),
            },
          },
        )
        console.log(productResponse.data.id)

        const productId = productResponse.data.id

        // 2. Загружаем фото
        if (this.uploadedPhotos.length > 0) {
          await this.uploadPhotos(productId)
        }

        // 3. Добавляем варианты
        const variantsResponse = await axios.post(
          `/api/admin/products/${productId}/variants`,
          {
            variants: this.variants,
          },
          {
            headers: {
              Authorization: localStorage.getItem('token'),
            },
          },
        )
        console.log(variantsResponse.data)

        this.showSuccess('Товар успешно создан!')
        this.closePopup()
      } catch (error) {
        console.error('Ошибка:', error)
        this.showError(error.response?.data?.error || 'Ошибка при создании товара')
      } finally {
        this.isSubmitting = false
      }
    },

    closePopup() {
      this.$emit('update:modelValue', false)

      this.clearForm()
    },
    showSuccess(message) {
      message = message || 'все гуд'
      this.flashMessage = {
        text: message,
        type: 'success',
      }
    },
    showError(message) {
      message = message || 'все гуд'
      this.flashMessage = {
        text: message,
        type: 'danger',
      }
    },
    async getCategories() {
      try {
        this.loadingCategories = true
        const response = await axios.get('/api/category/', {
          headers: {
            Authorization: localStorage.getItem('token'),
          },
        })
        this.sortOptions = response.data.categories.map((category) => ({
          value: category.id,
          label: category.name,
        }))
      } catch (error) {
        console.log(error)
      } finally {
        this.loadingCategories = false
      }
    },
    async addCategory() {
      try {
        this.loadingCategories = true
        const response = await axios.post(
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
        this.showSuccess(response.data.message)
        this.addNewCategory = ''
        this.getCategories()
      } catch (error) {
        console.log(error)
      } finally {
        this.loadingCategories = false
      }
    },
  },

  watch: {
    categoryCurrent: {
      handler(newValue) {
        this.form.category = newValue.value
      },
      immediate: true,
    },
  },

  created() {},
  mounted() {
    this.getCategories()
  },
  updated() {},
  beforeUnmount() {},
  unmounted() {},
}
</script>

<template>
  <div
    v-if="modelValue"
    class="bg-zinc-950/80 min-h-screen w-full fixed top-0 left-0 z-10 p-4 flex items-center justify-center"
  >
    <div class="container max-w-lg p-5 bg-white rounded-2xl flex flex-col gap-4">
      <div class="text-2xl font-bold text-center">
        <div class="">Добавить товар</div>
      </div>
      <div class="">
        <div class="" v-if="currentStep === 0">
          <form
            id="createProduct"
            @keydown.enter.prevent
            @submit.prevent="nextStep"
            class="flex flex-col gap-4"
          >
            <div>
              <label for="name" class="block text-sm font-medium text-gray-700"
                >Введите название товара</label
              >
              <input
                id="name"
                v-model="form.name"
                type="text"
                required
                class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              />
            </div>
            <div>
              <label for="description" class="block text-sm font-medium text-gray-700"
                >Введите описание товара</label
              >
              <input
                id="description"
                v-model="form.description"
                type="text"
                required
                class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              />
            </div>
            <div>
              <div class="block text-sm font-medium text-gray-700">Выберите категорию товара</div>
              <div v-if="!loadingCategories" class="grid mt-1 grid-cols-2 gap-2">
                <VSelect
                  @click.stop.prevent
                  class="*:w-full"
                  v-model="categoryCurrent"
                  :options="sortOptions"
                />
              </div>
              <VBlockLoader v-else class="w-1/2 h-9" />
            </div>

            <div class="grid grid-cols-2 gap-2">
              <div class="flex justify-start">
                <VButton @click="closePopup" :disabled="loadingCategories" :variant="'danger'">
                  <span v-if="loadingCategories" class="animate-spin"
                    ><svg
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
                    </svg>
                  </span>
                  <span v-else>Отменить</span>
                </VButton>
              </div>
              <div class="flex justify-end">
                <VButton :type="'submit'" :disabled="loadingCategories">
                  <span v-if="loadingCategories" class="animate-spin"
                    ><svg
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
                    </svg>
                  </span>
                  <span v-else>К добавлению фото</span>
                </VButton>
              </div>
            </div>
          </form>
        </div>

        <div v-if="currentStep === 1">
          <form @submit.prevent="nextStep" class="flex flex-col gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700">Фотографии товара</label>
              <div class="mt-1 border-2 border-dashed border-gray-300 rounded-md p-4">
                <input
                  type="file"
                  multiple
                  accept="image/*"
                  @change="handleFileUpload"
                  class="hidden"
                  ref="fileInput"
                />
                <div @click="$refs.fileInput.click()" class="cursor-pointer text-center py-8">
                  <div class="text-gray-500">Перетащите файлы сюда или кликните для выбора</div>
                </div>

                <!-- Превью загруженных фото -->
                <div v-if="uploadedPhotos.length" class="grid grid-cols-3 gap-2 mt-4">
                  <div v-for="(photo, index) in uploadedPhotos" :key="index" class="relative">
                    <img :src="photo.preview" class="w-full h-24 object-cover rounded" />
                    <button
                      @click.stop="removePhoto(index)"
                      class="absolute top-1 right-1 bg-red-500 text-white rounded-full w-5 h-5 flex items-center justify-center"
                    >
                      ×
                    </button>
                  </div>
                </div>
              </div>
            </div>

            <div class="grid grid-cols-2 gap-2">
              <VButton @click="prevStep" type="button"> Назад </VButton>
              <VButton type="submit"> Далее: Варианты товара </VButton>
            </div>
          </form>
        </div>

        <div v-if="currentStep === 2">
          <form @submit.prevent="submitAll" class="flex flex-col gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700">Варианты товара</label>

              <!-- Список вариантов -->
              <div
                v-for="(variant, index) in variants"
                :key="index"
                class="border rounded-md p-3 mb-3"
              >
                <div class="grid grid-cols-3 gap-3 mb-3">
                  <div>
                    <label class="block text-xs text-gray-500">Название</label>
                    <input
                      v-model="variant.name"
                      type="text"
                      required
                      class="mt-1 block w-full px-2 py-1 border border-gray-300 rounded-md text-sm"
                    />
                  </div>
                  <div>
                    <label class="block text-xs text-gray-500">Цена</label>
                    <input
                      v-model="variant.price"
                      type="number"
                      min="0"
                      required
                      class="mt-1 block w-full px-2 py-1 border border-gray-300 rounded-md text-sm"
                    />
                  </div>
                  <div>
                    <label class="block text-xs text-gray-500">Наличие</label>
                    <input
                      v-model="variant.stock"
                      type="number"
                      min="0"
                      required
                      class="mt-1 block w-full px-2 py-1 border border-gray-300 rounded-md text-sm"
                    />
                  </div>
                </div>
                <button @click="removeVariant(index)" type="button" class="text-red-500 text-sm">
                  Удалить вариант
                </button>
              </div>

              <button @click="addVariant" type="button" class="text-indigo-600 text-sm">
                + Добавить вариант
              </button>
            </div>

            <div class="grid grid-cols-2 gap-2">
              <VButton @click="prevStep" type="button"> Назад </VButton>
              <VButton type="submit" :disabled="isSubmitting">
                {{ isSubmitting ? 'Сохранение...' : 'Сохранить товар' }}
              </VButton>
            </div>
          </form>
        </div>

        <form v-if="currentStep == 0" @submit.prevent="addCategory" class="mt-4">
          <VFlashMessage
            v-if="flashMessage"
            :message="flashMessage.text"
            :type="flashMessage.type"
          ></VFlashMessage>
          <div>
            <label for="addNewCategory" class="block text-sm font-medium text-gray-700"
              >Не нашли категорию? Добавьте</label
            >
            <div class="flex items-end gap-2">
              <input
                v-if="!loadingCategories"
                id="addNewCategory"
                v-model="addNewCategory"
                type="text"
                required
                class="mt-1 flex-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              />
              <VBlockLoader v-else class="flex-1 h-[42px] mt-1" />
              <VButton :disabled="loadingCategories" class="!p-2 w-[42px] h-[42px]">
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
    </div>
  </div>
</template>

<style scoped></style>

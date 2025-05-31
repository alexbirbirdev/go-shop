<script>
import axios from 'axios'
import VButton from '@/components/forms/VButton.vue'
import VBlockLoader from '@/components/loaders/VBlockLoader.vue'

export default {
  name: 'ProfileView',

  components: { VButton, VBlockLoader },

  props: {},

  data() {
    return {
      user: {
        name: '',
        last_name: '',
        role: '',
        username: '',
        registrationDate: '',
      },
      form: {
        name: '',
        last_name: '',
      },
      loadingProfile: true,
      showUpdateProfile: false,
      loadingUpdateProfile: false,
      formCreateAddress: {
        title: '',
        city: '',
        street: '',
        house: '',
        apartment: '',
        floor: '',
        comment: '',
        is_default: true,
      },
      showCreateAddressForm: false,
      addresses: [],
      loadingAddress: true,
      isEditPopupVisible: false,
      editForm: {
        ID: null,
        title: '',
        city: '',
        street: '',
        house: '',
        apartment: '',
        floor: '',
        comment: '',
      },
    }
  },

  computed: {
    fullName() {
      if (this.user.name != '' || this.user.last_name != '') {
        return `${this.user.name} ${this.user.last_name}`
      }
      return `Не заполнено`
    },
    userInitials() {
      if (this.user.name != '' || this.user.last_name != '') {
        return `${this.user.name[0]}${this.user.last_name[0]}`
      }
      return `НН`
    },
    role() {
      if (this.user.role == 'user') {
        return `Рядовой пользователь`
      }
      if (this.user.role == 'admin') {
        return 'Админ'
      }
      return `Ноунейм галимый`
    },
  },

  methods: {
    openEditPopup(address) {
      this.editForm = { ...address }
      this.isEditPopupVisible = true
    },
    closeEditPopup() {
      this.isEditPopupVisible = false
    },
    cancelUpdateProfile() {
      this.showUpdateProfile = false
      this.form.name = this.user.name
      this.form.last_name = this.user.last_name
      this.loadingProfile = false
    },
    cancelCreateAddress() {
      this.showCreateAddressForm = false
    },
    handleSubmit() {
      alert('adfa')
    },
    editProfile() {
      this.showUpdateProfile = true
    },
    showCreateForm() {
      this.showCreateAddressForm = true
    },
    logout() {
      localStorage.removeItem('token')
      localStorage.removeItem('role')
      this.$store.commit('SET_AUTH', false)
      this.$router.push('/auth/signin')
    },
    async getProfile() {
      try {
        this.loadingProfile = true
        const response = await axios.get('http://localhost:8080/profile/', {
          headers: {
            Authorization: localStorage.getItem('token'),
          },
        })
        this.user.name = response.data.user.name
        this.user.last_name = response.data.user.last_name
        this.form.name = response.data.user.name
        this.form.last_name = response.data.user.last_name
        this.user.username = response.data.user.username
        this.user.registrationDate = response.data.user.created_at
        this.user.role = response.data.user.role
      } catch (error) {
        console.log(error)
      } finally {
        this.loadingProfile = false
      }
    },
    async updateProfile() {
      try {
        this.loadingUpdateProfile = true
        await axios.put(
          'http://localhost:8080/profile/',
          {
            name: this.form.name,
            last_name: this.form.last_name,
          },
          {
            headers: {
              Authorization: localStorage.getItem('token'),
            },
          },
        )
        this.user.name = this.form.name
        this.user.last_name = this.form.last_name
        this.showUpdateProfile = false
      } catch (error) {
        console.log(error)
      } finally {
        this.loadingProfile = false
        this.loadingUpdateProfile = false
      }
    },
    async createAddress() {
      try {
        this.loadingAddress = true
        await axios.post(
          'http://localhost:8080/user-addresses/',
          {
            title: this.formCreateAddress.title,
            city: this.formCreateAddress.city,
            street: this.formCreateAddress.street,
            house: this.formCreateAddress.house,
            apartment: this.formCreateAddress.apartment,
            floor: this.formCreateAddress.floor,
            comment: this.formCreateAddress.comment,
            is_default: this.formCreateAddress.is_default,
          },
          {
            headers: {
              Authorization: localStorage.getItem('token'),
            },
          },
        )

        this.showCreateAddressForm = false
        this.getAddresses()
      } catch (error) {
        console.log(error)
      } finally {
        this.loadingAddress = false
      }
    },
    async updateAddress() {
      console.log(this.editForm.ID)
      try {
        this.loadingAddress = true
        await axios.put(
          'http://localhost:8080/user-addresses/' + this.editForm.ID,
          this.editForm,
          {
            headers: {
              Authorization: localStorage.getItem('token'),
            },
          },
        )

        this.isEditPopupVisible = false
        this.getAddresses()
      } catch (error) {
        console.log(error)
      } finally {
        this.loadingAddress = false
      }
    },
    async makeDefaultAddress(id) {
      try {
        this.loadingAddress = true
        await axios.put(
          'http://localhost:8080/user-addresses/' + id,
          {
            is_default: true,
          },
          {
            headers: {
              Authorization: localStorage.getItem('token'),
            },
          },
        )
        this.getAddresses()
      } catch (error) {
        console.log(error)
      } finally {
        this.loadingAddress = false
      }
    },
    async getAddresses() {
      try {
        this.loadingAddress = true
        const response = await axios.get('http://localhost:8080/user-addresses/', {
          headers: {
            Authorization: localStorage.getItem('token'),
          },
        })
        this.addresses = response.data.user_addresses
      } catch (error) {
        console.log(error)
      } finally {
        this.loadingAddress = false
      }
    },
    async deleteAddress(id) {
      try {
        await axios.delete('http://localhost:8080/user-addresses/' + id, {
          headers: {
            Authorization: localStorage.getItem('token'),
          },
        })
        this.addresses = this.addresses.filter((address) => address.ID !== id)
      } catch (error) {
        console.log(error)
      }
    },
  },

  watch: {},

  created() {
    this.getProfile()
    this.getAddresses()
  },
  mounted() {},
  updated() {},
  beforeUnmount() {},
  unmounted() {},
}
</script>

<template>
  <div class="flex flex-col gap-10">
    <div
      v-if="showUpdateProfile"
      class="bg-zinc-950/80 min-h-screen w-full fixed top-0 left-0 z-10 p-4 flex items-center justify-center"
    >
      <div class="container max-w-sm bg-white p-5 rounded-xl flex flex-col gap-4">
        <div class="text-xl font-bold text-center">Редактирование профиля</div>
        <form class="space-y-6" @submit.prevent="updateProfile">
          <div class="space-y-4">
            <div>
              <label for="name" class="block text-sm font-medium text-gray-700">Введите имя</label>
              <input
                id="name"
                placeholder="Введите имя"
                v-model="form.name"
                type="text"
                required
                class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              />
            </div>
            <div>
              <label for="last_name" class="block text-sm font-medium text-gray-700"
                >Введите фамилию</label
              >
              <input
                id="last_name"
                placeholder="Введите фамилию"
                v-model="form.last_name"
                type="text"
                required
                class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              />
            </div>
          </div>
          <div class="grid grid-cols-3 gap-2">
            <VButton
              class="col-span-2"
              @submit.prevent="handleSubmit"
              :disabled="loadingUpdateProfile"
              ><span v-if="!loadingUpdateProfile">Изменить</span>
              <span v-else>Загрузка...</span></VButton
            >
            <VButton
              @click="cancelUpdateProfile"
              :disabled="loadingUpdateProfile"
              :variant="'danger'"
              ><span v-if="!loadingUpdateProfile">Отменить</span>
              <span v-else>Загрузка...</span></VButton
            >
          </div>
        </form>
      </div>
    </div>
    <div
      v-if="showCreateAddressForm"
      class="bg-zinc-950/80 min-h-screen w-full fixed top-0 left-0 z-10 p-4 flex items-center justify-center"
    >
      <div class="container max-w-xl bg-white p-5 rounded-xl flex flex-col gap-4">
        <div class="text-xl font-bold text-center">Добавить новый адрес</div>
        <form class="space-y-6" @submit.prevent="createAddress">
          <div class="grid grid-cols-2 gap-4 *:mt-auto">
            <div>
              <label for="formCreateAddressTitle" class="block text-sm font-medium text-gray-700"
                >Введите название адреса</label
              >
              <input
                id="formCreateAddressTitle"
                placeholder="Введите название адреса"
                v-model="formCreateAddress.title"
                type="text"
                required
                class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              />
            </div>
            <div>
              <label for="formCreateAddressCity" class="block text-sm font-medium text-gray-700"
                >Введите город</label
              >
              <input
                id="formCreateAddressCity"
                placeholder="Введите город"
                v-model="formCreateAddress.city"
                type="text"
                required
                class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              />
            </div>
            <div>
              <label for="formCreateAddressStreet" class="block text-sm font-medium text-gray-700"
                >Введите улицу</label
              >
              <input
                id="formCreateAddressStreet"
                placeholder="Введите улицу"
                v-model="formCreateAddress.street"
                type="text"
                required
                class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              />
            </div>
            <div>
              <label for="formCreateAddressHouse" class="block text-sm font-medium text-gray-700"
                >Введите дом (с корпусом | подъездом | строением)</label
              >
              <input
                id="formCreateAddressHouse"
                placeholder="Введите дом (с корпусом | подъездом | строением)"
                v-model="formCreateAddress.house"
                type="text"
                required
                class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              />
            </div>
            <div>
              <label
                for="formCreateAddressApartment"
                class="block text-sm font-medium text-gray-700"
                >Введите квартиру</label
              >
              <input
                id="formCreateAddressApartment"
                placeholder="Введите квартиру"
                v-model="formCreateAddress.apartment"
                type="text"
                required
                class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              />
            </div>
            <div>
              <label for="formCreateAddressFloor" class="block text-sm font-medium text-gray-700"
                >Введите этаж</label
              >
              <input
                id="formCreateAddressFloor"
                placeholder="Введите этаж"
                v-model="formCreateAddress.floor"
                type="text"
                required
                class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              />
            </div>
            <div class="col-span-2">
              <label for="formCreateAddressComment" class="block text-sm font-medium text-gray-700"
                >Введите комментарий</label
              >
              <input
                id="formCreateAddressComment"
                placeholder="Введите комментарий"
                v-model="formCreateAddress.comment"
                type="text"
                required
                class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              />
            </div>
            <div class="flex items-center col-span-2" v-if="addresses.length > 0">
              <input
                id="formCreateAddressDefault"
                v-model="formCreateAddress.is_default"
                type="checkbox"
                checked
                class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
              />
              <label for="formCreateAddressDefault" class="ml-2 block text-sm text-gray-700">
                Сделать основным
              </label>
            </div>
          </div>
          <div class="grid grid-cols-3 gap-2">
            <VButton class="col-span-2" :disabled="loadingUpdateProfile"
              ><span v-if="!loadingUpdateProfile">Добавить</span>
              <span v-else>Загрузка...</span></VButton
            >
            <VButton
              @click="cancelCreateAddress"
              :disabled="loadingUpdateProfile"
              :variant="'danger'"
              ><span v-if="!loadingUpdateProfile">Отменить</span>
              <span v-else>Загрузка...</span></VButton
            >
          </div>
        </form>
      </div>
    </div>
    <div
      v-if="isEditPopupVisible"
      class="bg-zinc-950/80 min-h-screen w-full fixed top-0 left-0 z-10 p-4 flex items-center justify-center"
    >
      <div class="container max-w-xl bg-white p-5 rounded-xl flex flex-col gap-4">
        <div class="text-xl font-bold text-center">Добавить новый адрес</div>
        <form class="space-y-6" @submit.prevent="updateAddress">
          <div class="grid grid-cols-2 gap-4 *:mt-auto">
            <div>
              <label for="editFormTitle" class="block text-sm font-medium text-gray-700"
                >Введите название адреса</label
              >
              <input
                id="editFormTitle"
                placeholder="Введите название адреса"
                v-model="editForm.title"
                type="text"
                required
                class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              />
            </div>
            <div>
              <label for="editFormCity" class="block text-sm font-medium text-gray-700"
                >Введите город</label
              >
              <input
                id="editFormCity"
                placeholder="Введите город"
                v-model="editForm.city"
                type="text"
                required
                class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              />
            </div>
            <div>
              <label for="editFormStreet" class="block text-sm font-medium text-gray-700"
                >Введите улицу</label
              >
              <input
                id="editFormStreet"
                placeholder="Введите улицу"
                v-model="editForm.street"
                type="text"
                required
                class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              />
            </div>
            <div>
              <label for="editFormHouse" class="block text-sm font-medium text-gray-700"
                >Введите дом (с корпусом | подъездом | строением)</label
              >
              <input
                id="editFormHouse"
                placeholder="Введите дом (с корпусом | подъездом | строением)"
                v-model="editForm.house"
                type="text"
                required
                class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              />
            </div>
            <div>
              <label
                for="editFormApartment"
                class="block text-sm font-medium text-gray-700"
                >Введите квартиру</label
              >
              <input
                id="editFormApartment"
                placeholder="Введите квартиру"
                v-model="editForm.apartment"
                type="text"
                required
                class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              />
            </div>
            <div>
              <label for="editFormFloor" class="block text-sm font-medium text-gray-700"
                >Введите этаж</label
              >
              <input
                id="editFormFloor"
                placeholder="Введите этаж"
                v-model="editForm.floor"
                type="text"
                required
                class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              />
            </div>
            <div class="col-span-2">
              <label for="editFormComment" class="block text-sm font-medium text-gray-700"
                >Введите комментарий</label
              >
              <input
                id="editFormComment"
                placeholder="Введите комментарий"
                v-model="editForm.comment"
                type="text"
                required
                class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              />
            </div>
          </div>
          <div class="grid grid-cols-3 gap-2">
            <VButton class="col-span-2" :disabled="loadingAddress"
              ><span v-if="!loadingAddress">Добавить</span> <span v-else>Загрузка...</span></VButton
            >
            <VButton @click="closeEditPopup" :disabled="loadingAddress" :variant="'danger'"
              ><span v-if="!loadingAddress">Отменить</span> <span v-else>Загрузка...</span></VButton
            >
          </div>
        </form>
      </div>
    </div>
    <div class="flex flex-col gap-5">
      <div class="bg-indigo-600 p-6 text-white rounded-2xl">
        <div class="flex items-center space-x-4">
          <div class="flex-shrink-0">
            <div
              class="h-16 w-16 rounded-full bg-indigo-400 flex items-center justify-center text-2xl font-bold"
            >
              <VBlockLoader v-if="loadingProfile || loadingUpdateProfile" :class="'w-8 h-4'" />
              <span v-else>{{ userInitials }}</span>
            </div>
          </div>
          <div class="flex flex-col gap-2">
            <VBlockLoader v-if="loadingProfile || loadingUpdateProfile" :class="'w-30 h-5'" />
            <h1 v-else class="text-xl font-bold leading-[120%]">{{ fullName }}</h1>
            <VBlockLoader v-if="loadingProfile" :class="'w-15 h-4'" />
            <p v-else class="text-indigo-100 leading-[120%]">{{ role }}</p>
          </div>
        </div>
      </div>
      <div class="border-b border-gray-200 pb-4">
        <h2 class="text-md font-medium text-black">Основная информация</h2>
        <div class="mt-2 space-y-2">
          <div class="flex justify-between">
            <div class="text-sm text-gray-500">Логин</div>
            <div class="text-sm text-gray-900">
              <VBlockLoader v-if="loadingProfile" :class="'w-30 h-5'" />
              <span v-else>{{ user.username }}</span>
            </div>
          </div>
          <div class="flex justify-between">
            <div class="text-sm text-gray-500">Дата регистрации</div>
            <div class="text-sm text-gray-900">
              <VBlockLoader v-if="loadingProfile" :class="'w-30 h-5'" />
              <span v-else>{{ user.registrationDate }}</span>
            </div>
          </div>
        </div>
      </div>
      <div class="">
        <VButton v-if="!loadingProfile" class="block" :variant="'warning'" @click="editProfile"
          >Редактировать</VButton
        >
      </div>
    </div>

    <div class="flex flex-col gap-5">
      <div class="border-b border-gray-200 pb-4 space-y-4">
        <h2 class="text-md font-medium text-black">Сохраненные адреса</h2>
        <div v-if="loadingAddress" class="grid grid-cols-3 gap-2">
          <VBlockLoader v-for="i in 6" :key="i" class="h-40" />
        </div>
        <div v-else class="grid grid-cols-3 gap-2">
          <div
            v-for="ad in addresses"
            :key="ad.id"
            class="flex flex-col justify-between gap-2 bg-white rounded-xl p-4"
            :class="ad.is_default ? 'shadow-2xl' : ''"
          >
            <div class="">
              <div class="text-md">{{ ad.title }}</div>
              <div class="text-gray-400">
                {{ ad.city }}, {{ ad.street }}, {{ ad.house }}, {{ ad.apartment }}, {{ ad.floor }}
              </div>
            </div>

            <div class="grid grid-cols-2 gap-2 *:text-xs">
              <VButton v-if="!ad.is_default" @click="deleteAddress(ad.ID)" :variant="'danger'"
                >Удалить</VButton
              >
              <VButton v-if="!ad.is_default" @click="makeDefaultAddress(ad.ID)"
                >Сделать основным</VButton
              >
              <VButton @click="openEditPopup(ad)" class="col-span-2" :variant="'warning'"
                >Изменить</VButton
              >
            </div>
          </div>
          <div
            @click="showCreateForm"
            class="min-h-40 flex items-center justify-center bg-white rounded-xl border border-green-500 text-green-600 duration-200 hover:scale-105 cursor-pointer"
          >
            <span>Добавить адрес</span>
          </div>
        </div>
      </div>
    </div>
    <div class="">
      <VButton v-if="!loadingProfile" :variant="'danger'" @click="logout">Выйти</VButton>
    </div>
  </div>
</template>

<style scoped></style>

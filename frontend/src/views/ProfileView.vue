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
        return `Админ`
      }
      return `Ноунейм галимый`
    },
  },

  methods: {
    cancelUpdateProfile() {
      this.showUpdateProfile = false
      this.form.name = this.user.name
      this.form.last_name = this.user.last_name
      this.loadingProfile = false
    },
    handleSubmit() {
      alert('adfa')
    },
    editProfile() {
      this.showUpdateProfile = true
      this.loadingProfile = true
    },
    logout() {
      localStorage.removeItem('token')
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
        this.loadingProfile = true
        this.loadingUpdateProfile = true
        const response = await axios.put(
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
        console.log(response.data.message)
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
  },

  watch: {},

  created() {
    this.getProfile()
  },
  mounted() {},
  updated() {},
  beforeUnmount() {},
  unmounted() {},
}
</script>

<template>
  <div class="flex flex-col gap-5">
    <div
      v-if="showUpdateProfile"
      class="bg-zinc-950/80 h-screen w-full fixed top-0 left-0 z-10 p-4 flex items-center justify-center"
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
    <div class="bg-indigo-600 p-6 text-white rounded-2xl">
      <div class="flex items-center space-x-4">
        <div class="flex-shrink-0">
          <div
            class="h-16 w-16 rounded-full bg-indigo-400 flex items-center justify-center text-2xl font-bold"
          >
            <VBlockLoader v-if="loadingProfile" :class="'w-8 h-4'" />
            <span v-else>{{ userInitials }}</span>
          </div>
        </div>
        <div class="flex flex-col gap-2">
          <VBlockLoader v-if="loadingProfile" :class="'w-30 h-5'" />
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
    <div class="flex flex-col gap-2 items-start">
      <VButton v-if="!loadingProfile" :variant="'warning'" @click="editProfile"
        >Редактировать</VButton
      >
      <VButton v-if="!loadingProfile" :variant="'danger'" @click="logout">Выйти</VButton>
    </div>
  </div>
</template>

<style scoped></style>

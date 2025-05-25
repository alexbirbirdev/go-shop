<script>
import axios from 'axios'
import VButton from '@/components/forms/VButton.vue'

export default {
  name: 'ProfileView',

  components: { VButton },

  props: {},

  data() {
    return {
      user: {
        firstName: '',
        lastName: '',
        role: '',
        username: '',
        registrationDate: '',
      },
      loadingProfile: true,
    }
  },

  computed: {
    fullName() {
      if (this.user.firstName != '' || this.user.lastName != '') {
        return `${this.user.firstName} ${this.user.lastName}`
      }
      return `Не заполнено`
      
    },
    userInitials() {
      if (this.user.firstName != '' || this.user.lastName != '') {
        return `${this.user.firstName[0]}${this.user.lastName[0]}`
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
    editProfile() {
      alert('Редактирование профиля')
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
        this.user.firstName = response.data.user.name
        this.user.lastName = response.data.user.last_name
        this.user.username = response.data.user.username
        this.user.registrationDate = response.data.user.created_at
        this.user.role = response.data.user.role
      } catch (error) {
        console.log(error)
      } finally {
        this.loadingProfile = false
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
    <div class="bg-indigo-600 p-6 text-white rounded-2xl">
      <div class="flex items-center space-x-4">
        <div class="flex-shrink-0">
          <div
            class="h-16 w-16 rounded-full bg-indigo-400 flex items-center justify-center text-2xl font-bold"
          >
            {{ userInitials }}
          </div>
        </div>
        <div>
          <h1 class="text-xl font-bold">{{ fullName }}</h1>
          <p class="text-indigo-100">{{ role }}</p>
        </div>
      </div>
    </div>
    <div class="border-b border-gray-200 pb-4">
      <h2 class="text-md font-medium text-black">Основная информация</h2>
      <div class="mt-2 space-y-2">
        <div class="flex justify-between">
          <div class="text-sm text-gray-500">Логин</div>
          <div class="text-sm text-gray-900">{{ user.username }}</div>
        </div>
        <div class="flex justify-between">
          <div class="text-sm text-gray-500">Дата регистрации</div>
          <div class="text-sm text-gray-900">{{ user.registrationDate }}</div>
        </div>
      </div>
    </div>
    <div class="flex flex-col gap-2 items-start">
      <VButton variant="warning" @click="editProfile">Редактировать</VButton>
      <VButton variant="danger" @click="logout">Выйти</VButton>
    </div>
  </div>
</template>

<style scoped></style>

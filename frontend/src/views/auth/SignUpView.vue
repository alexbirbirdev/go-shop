<script>
import VFlashMessage from '@/components/forms/VFlashMessage.vue'
import axios from 'axios'
export default {
  name: 'SignUpView',

  components: { VFlashMessage },

  props: {},

  data() {
    return {
      loading: false,
      form: {
        email: '',
        password: '',
        confirmPassword: '',
        agreeTerms: false,
      },
      flashMessage: null,
    }
  },

  computed: {},

  methods: {
    showSuccess(message) {
      message = message || 'все гуд'
      this.flashMessage = {
        text: message,
        type: 'success',
      }
    },
    showError(message) {
      message = message || 'Произошла ошибка'
      this.flashMessage = {
        text: message,
        type: 'error',
      }
    },
    async handleSubmit() {
      if (this.form.password !== this.form.confirmPassword) {
        this.showError('Пароли не совпадают')
        return
      }

      this.loading = true
      try {
        const response = await axios.post(
          '/api/auth/signup',
          {
            email: this.form.email,
            password: this.form.password,
          },
          {
            headers: {
              'Content-Type': 'application/json',
            },
          },
        )
        console.log(response.data)

        // const redirect = this.$route.query.redirect || '/'
        // await this.$router.push(redirect)
        this.showSuccess(response.data.message + '. Авторизируйтесь!')
        this.form.email = null
        this.form.password = null
        this.form.confirmPassword = null
        this.form.agreeTerms = null
      } catch (error) {
        // console.error('Ошибка регистрации:', error.response.data.error)
        this.showError(error.response.data.error)
      } finally {
        this.loading = false
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
  <div class="min-h-screen flex items-center justify-center bg-gray-50">
    <VFlashMessage v-if="flashMessage" :message="flashMessage.text" :type="flashMessage.type" />
    <div class="w-full max-w-md p-8 space-y-8 bg-white rounded-lg shadow">
      <h2 class="text-2xl font-bold text-center text-gray-900">Создать аккаунт</h2>
      <form class="mt-8 space-y-6" @submit.prevent="handleSubmit">
        <div class="space-y-4">
          <div>
            <label for="email" class="block text-sm font-medium text-gray-700">Email</label>
            <input
              id="email"
              v-model="form.email"
              type="email"
              required
              class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
            />
          </div>
          <div>
            <label for="password" class="block text-sm font-medium text-gray-700">Пароль</label>
            <input
              id="password"
              v-model="form.password"
              type="password"
              required
              class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
            />
          </div>
          <div>
            <label for="confirmPassword" class="block text-sm font-medium text-gray-700"
              >Подтвердите пароль</label
            >
            <input
              id="confirmPassword"
              v-model="form.confirmPassword"
              type="password"
              required
              class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
            />
          </div>
        </div>

        <div class="flex items-center">
          <input
            id="terms"
            v-model="form.agreeTerms"
            type="checkbox"
            required
            class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
          />
          <label for="terms" class="ml-2 block text-sm text-gray-700">
            Я согласен с условиями использования
          </label>
        </div>

        <div>
          <button
            type="submit"
            class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
            :disabled="loading"
          >
            <span v-if="!loading">Зарегистрироваться</span>
            <span v-else>Загрузка...</span>
          </button>
        </div>
      </form>
      <div class="flex flex-col items-center gap-2 text-xs">
        <div class="">
          Есть аккаунт?
          <router-link class="text-blue-700 duration-200 hover:text-blue-500" to="/auth/signin"
            >Войти</router-link
          >
        </div>
        <router-link class="text-blue-700 duration-200 hover:text-blue-500" to="/"
          >Вернуться на главную</router-link
        >
      </div>
    </div>
  </div>
</template>

<style scoped></style>

<script>
import axios from 'axios'

export default {
  name: 'SignInView',

  components: {},

  props: {},

  data() {
    return {
      loading: false,
      form: {
        email: '',
        password: '',
      },
      message: '',
    }
  },

  computed: {},

  methods: {
    async handleSubmit() {
      this.loading = true
      try {
        const response = await axios.post('http://localhost:8080/auth/signin', this.form, {
          headers: {
            'Content-Type': 'application/json',
          },
        })

        localStorage.setItem('token', response.data.token, response.data.token)

        if (this.$store) {
          this.$store.commit('SET_AUTH', {
            isAuth: true,
            token: response.data.token,
          })
        } else {
          console.warn('Vuex store не доступен')
        }

        const redirect = this.$route.query.redirect || '/'
        await this.$router.push(redirect)
      } catch (error) {
        console.error('Ошибка авторизации:', error.response.data.error)
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
    <div class="w-full max-w-md p-8 space-y-8 bg-white rounded-lg shadow">
      <h2 class="text-2xl font-bold text-center text-gray-900">Вход в аккаунт</h2>
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
        </div>

        <div>
          <button
            type="submit"
            class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
            :disabled="loading"
          >
            <span v-if="!loading">Войти</span>
            <span v-else>Загрузка...</span>
          </button>
        </div>
      </form>
      <div class="flex flex-col items-center gap-2 text-xs">
        <div class="">
          Нет аккаунта?
          <router-link class="text-blue-700 duration-200 hover:text-blue-500" to="/auth/signup"
            >Зарегистрироваться</router-link
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

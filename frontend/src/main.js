import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import store from './stores'

import App from './App.vue'
import router from './router'

const app = createApp(App)

app.use(createPinia())
app.use(store)
app.use(router)

app.mount('#app')

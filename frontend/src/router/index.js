import CartView from '@/views/CartView.vue'
import CatalogView from '@/views/CatalogView.vue'
import FavoriteView from '@/views/FavoriteView.vue'
import OrderDetailView from '@/views/OrderDetailView.vue'
import OrdersView from '@/views/OrdersView.vue'
import ProductView from '@/views/ProductView.vue'
import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: ProductView,
      redirect: "/catalog",
    },
    {
      path: '/catalog',
      name: 'catalog',
      component: CatalogView,
    },
    {
      path: '/catalog/:id',
      name: 'product',
      component: ProductView,
    },
    {
      path: '/cart',
      name: 'cart',
      component: CartView,
    },
    {
      path: '/favorite',
      name: 'favorite',
      component: FavoriteView,
    },
    {
      path: '/orders',
      name: 'orders',
      component: OrdersView,
    },
    {
      path: '/orders/:id',
      name: 'order',
      component: OrderDetailView,
    },
    {
      path: '/auth/signin',
      name: 'signin',
      component: ProductView,
    },
    {
      path: '/auth/signup',
      name: 'signup',
      component: ProductView,
    },
    
  ],
})

export default router

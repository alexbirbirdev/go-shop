import AdminCategoriesView from '@/views/admin/AdminCategoriesView.vue'
import AdminOrdersView from '@/views/admin/AdminOrdersView.vue'
import AdminProductsView from '@/views/admin/AdminProductsView.vue'
import AdminView from '@/views/admin/AdminView.vue'
import UsersView from '@/views/admin/UsersView.vue'
import SignInView from '@/views/auth/SignInView.vue'
import SignUpView from '@/views/auth/SignUpView.vue'
import CartView from '@/views/CartView.vue'
import CatalogView from '@/views/CatalogView.vue'
import FavoriteView from '@/views/FavoriteView.vue'
import OrderDetailView from '@/views/OrderDetailView.vue'
import OrdersView from '@/views/OrdersView.vue'
import ProductView from '@/views/ProductView.vue'
import ProfileView from '@/views/ProfileView.vue'
import { createRouter, createWebHistory } from 'vue-router'
// import stores from '@/stores'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: ProductView,
      redirect: '/catalog',
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
      meta: {
        requiresAuth: true,
      },
    },
    {
      path: '/favorite',
      name: 'favorite',
      component: FavoriteView,
      meta: {
        requiresAuth: true,
      },
    },
    {
      path: '/profile',
      name: 'profile',
      component: ProfileView,
      meta: {
        requiresAuth: true,
      },
    },
    {
      path: '/orders',
      name: 'orders',
      component: OrdersView,
      meta: {
        requiresAuth: true,
      },
    },
    {
      path: '/orders/:id',
      name: 'order',
      component: OrderDetailView,
      meta: {
        requiresAuth: true,
      },
    },
    {
      path: '/auth/',
      name: 'redirectAuth',
      component: SignInView,
      redirect: '/auth/signin/',
    },
    {
      path: '/auth/signin',
      name: 'signin',
      component: SignInView,
      meta: {
        authPage: true,
      },
    },
    {
      path: '/auth/signup',
      name: 'signup',
      component: SignUpView,
      meta: {
        authPage: true,
      },
    },
    {
      path: '/admin',
      name: 'admin',
      component: AdminView,
      meta: {
        requiresAuth: true,
        adminPage: true,
      },
    },
    {
      path: '/admin/users',
      name: 'adminUsers',
      component: UsersView,
      meta: {
        requiresAuth: true,
        adminPage: true,
      },
    },
    {
      path: '/admin/orders',
      name: 'adminOrders',
      component: AdminOrdersView,
      meta: {
        requiresAuth: true,
        adminPage: true,
      },
    },
    {
      path: '/admin/products',
      name: 'adminProducts',
      component: AdminProductsView,
      meta: {
        requiresAuth: true,
        adminPage: true,
      },
    },
    {
      path: '/admin/categories',
      name: 'adminCategories',
      component: AdminCategoriesView,
      meta: {
        requiresAuth: true,
        adminPage: true,
      },
    },
  ],
})

router.beforeEach((to, from, next) => {
  const isAuth = !!localStorage.getItem('token')

  if (isAuth && (to.name === 'signin' || to.name === 'signup')) {
    next('/')
    return 
  }

  if (to.meta.requiresAuth && !isAuth) {
    next({ name: 'signin', query: { redirect: to.path } })
  } else {
    next()
  }
})

export default router

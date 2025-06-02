<script>
import VButton from '@/components/forms/VButton.vue'
import axios from 'axios'
import { mapMutations } from 'vuex'
export default {
  name: 'CheckoutView',

  components: {
    VButton,
  },

  props: {},

  data() {
    return {
      isCartLoading: true,
      products: [],
      total: 0,
      loadingProfile: true,
      form: {
        address: '',
        addressCom: '',
        phone: '',
        first_name: '',
        last_name: '',
      },
      loadingAddress: false,
      addresses: [],
      selectedAddressId: null,

      isLoading: false,
    }
  },

  computed: {},

  methods: {
    ...mapMutations('cart', ['INCREMENT', 'DECREMENT', 'SET_COUNT']),
    async getCart() {
      try {
        this.isCartLoading = true
        const response = await axios.get('http://localhost:8080/cart/', {
          headers: {
            Authorization: localStorage.getItem('token'),
          },
        })
        this.total = response.data.total
        this.products = response.data.cart_items
      } catch (error) {
        console.log(error)
      } finally {
        this.isCartLoading = false
      }
    },
    async getProfile() {
      if (!localStorage.getItem('role')) {
        this.$router.push('/auth/signin?redirect=/profile/')
        return
      }
      try {
        this.loadingProfile = true
        const response = await axios.get('http://localhost:8080/profile/', {
          headers: {
            Authorization: localStorage.getItem('token'),
          },
        })
        this.form.first_name = response.data.user.name
        this.form.last_name = response.data.user.last_name
      } catch (error) {
        console.log(error)
      } finally {
        this.loadingProfile = false
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

        const defaultAddr = this.addresses.find((addr) => addr.is_default)
        if (defaultAddr) {
          this.selectedAddressId = defaultAddr.ID
          this.applySelectedAddress(defaultAddr)
        }
      } catch (error) {
        console.log(error)
      } finally {
        this.loadingAddress = false
      }
    },

    applySelectedAddress(addr) {
      // Составляем строку "город, улица, дом, квартира, этаж"
      this.form.address = [addr.city, addr.street, addr.house, addr.apartment, addr.floor]
        .filter(Boolean)
        .join(', ')

      // Сохраняем комментарий
      this.form.addressCom = addr.comment || ''
    },

    onAddressChange() {
      const addr = this.addresses.find((a) => a.ID === this.selectedAddressId)
      if (addr) {
        this.applySelectedAddress(addr)
      }
    },
    async handleForm() {
      try {
        this.isLoading = true
        await axios.post(
          'http://localhost:8080/orders/',
          {
            address: this.form.address,
            address_com: this.form.addressCom,
            phone: this.form.phone,
            first_name: this.form.first_name,
            last_name: this.form.last_name,
          },
          {
            headers: {
              Authorization: localStorage.getItem('token'),
            },
          },
        )
        this.$store.commit('cart/SET_COUNT', 0)
        this.$router.push('/orders')
      } catch (error) {
        console.log(error)
      } finally {
        this.isLoading = false
      }
    },
  },

  watch: {},

  created() {
    this.getProfile()
    this.getCart()
    this.getAddresses()
  },
  mounted() {},
  updated() {},
  beforeUnmount() {},
  unmounted() {},
}
</script>

<template>
  <div class="max-w-3xl mx-auto">
    <div class="text-4xl font-bold mb-10">Оформление заказа</div>
    <form @submit.prevent.stop="handleForm" class="grid grid-cols-2 gap-10">
      <div class="flex flex-col gap-5">
        <div>
          <label for="phone" class="block text-sm font-medium text-gray-700"
            >Номер телефона получателя</label
          >
          <input
            id="phone"
            v-model="form.phone"
            type="text"
            required
            class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
          />
        </div>
        <div>
          <label for="first_name" class="block text-sm font-medium text-gray-700"
            >Имя получателя</label
          >
          <input
            id="first_name"
            v-model="form.first_name"
            type="text"
            required
            class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
          />
        </div>
        <div>
          <label for="last_name" class="block text-sm font-medium text-gray-700"
            >Фамилия получателя</label
          >
          <input
            id="last_name"
            v-model="form.last_name"
            type="text"
            required
            class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
          />
        </div>
        <!-- Здесь добавляем селект для выбора адреса -->
        <div>
          <label for="addressSelect" class="block text-sm font-medium text-gray-700">
            Выберите адрес доставки
          </label>
          <select
            id="addressSelect"
            v-model="selectedAddressId"
            @change="onAddressChange"
            class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
          >
            <!-- Если адреса пока не загрузились, можно показать один пустой option -->
            <option value="" disabled>Загрузка адресов...</option>

            <!-- Иначе перечисляем их -->
            <option v-for="addr in addresses" :key="addr.ID" :value="addr.ID">
              {{ addr.city }}, {{ addr.street }}, д.{{ addr.house }}, кв.{{ addr.apartment }}, эт.{{
                addr.floor
              }}
            </option>
          </select>
        </div>
      </div>
      <div class="flex flex-col gap-10">
        <div class="grid gap-2">
          <div class="flex gap-2" v-for="product in products" :key="product">
            <div
              class="aspect-square flex items-center select-none justify-center border-2 border-neutral-300 rounded-lg w-25 overflow-hidden"
            >
              <img
                v-if="!product.product_variant.Product.images.length"
                src="https://placehold.co/600x400"
                :alt="product.product_variant.Product.name"
                class="max-h-full max-w-full"
              />
              <img
                v-else
                :src="product.product_variant.Product.images[0].image_url"
                :alt="product.product_variant.Product.name"
                class="max-h-full max-w-full"
              />
            </div>
            <div class="flex-1 leading-[120%] flex flex-col justify-between">
              <div class="">
                <div class="text-md font-medium mb-2 line-clamp-1">
                  {{ product.product_variant.Product.name }}
                </div>
                <div class="grid gap-1">
                  <div class="text-xs text-gray-600">
                    Вариант:
                    <span
                      class="py-1 px-3 text-xs rounded-2xl bg-neutral-200 duration-200 select-none"
                    >
                      {{ product.product_variant.name }}</span
                    >
                  </div>
                  <div class="text-xs text-gray-600">
                    Кол-во: <b>{{ product.quantity }}</b>
                  </div>
                </div>
              </div>

              <div class="mt-auto text-xs">
                <span class="!text-sm font-bold">{{ product.product_variant.price }} ₽</span> / шт.
              </div>
            </div>
          </div>
        </div>
        <VButton :type="'submit'" class="!justify-between">
          <span>Оплатить:</span>
          <span class="text-base">{{ total }} ₽</span>
        </VButton>
      </div>
    </form>
  </div>
</template>

<style scoped></style>

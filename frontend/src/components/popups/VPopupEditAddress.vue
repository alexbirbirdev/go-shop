<script>
import VButton from '@/components/forms/VButton.vue'
export default {
  name: 'VPopupEditAddress',

  components: { VButton },

  props: {
    visible: {
      type: Boolean,
      required: true,
    },
    address: {
      type: Object,
      default: () => ({}),
    },
  },

  data() {
    return {
      //   isVisible: false,
      form: {
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

  computed: {},

  methods: {
    handleSubmit() {
      this.$emit('save', this.form)
      //   this.$emit('update:visible', false); // Закрываем попап после сохранения
    },
    handleCancel() {
      this.$emit('update:visible', false) // Закрываем попап при отмене
    },
  },

  watch: {
    address: {
      immediate: true,
      handler(newVal) {
        if (newVal) {
          this.form = { ...newVal }
        }
      },
    },
  },

  created() {},
  mounted() {},
  updated() {},
  beforeUnmount() {},
  unmounted() {},
}
</script>

<template>
  <div
    v-if="visible"
    class="bg-zinc-950/80 min-h-screen w-full fixed top-0 left-0 z-10 p-4 flex items-center justify-center"
  >
    <div class="container max-w-xl bg-white p-5 rounded-xl flex flex-col gap-4">
      <div class="text-xl font-bold text-center">Редактирование адрес</div>
      <form class="space-y-6" >
        <div class="grid grid-cols-2 gap-4 *:mt-auto">
          <div>
            <label for="formCreateAddressTitle" class="block text-sm font-medium text-gray-700"
              >Введите название адреса</label
            >
            <input
              id="formCreateAddressTitle"
              placeholder="Введите название адреса"
              v-model="form.title"
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
              v-model="form.city"
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
              v-model="form.street"
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
              v-model="form.house"
              type="text"
              required
              class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
            />
          </div>
          <div>
            <label for="formCreateAddressApartment" class="block text-sm font-medium text-gray-700"
              >Введите квартиру</label
            >
            <input
              id="formCreateAddressApartment"
              placeholder="Введите квартиру"
              v-model="form.apartment"
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
              v-model="form.floor"
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
              v-model="form.comment"
              type="text"
              required
              class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
            />
          </div>
        </div>
        <div class="grid grid-cols-3 gap-2">
          <VButton @submit.prevent="handleSubmit" class="col-span-2"><span>Изменить</span> </VButton>
          <VButton @click="handleCancel" :variant="'danger'">Отменить </VButton>
        </div>
      </form>
    </div>
  </div>
</template>

<style scoped></style>

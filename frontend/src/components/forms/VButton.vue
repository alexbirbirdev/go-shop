<script>
export default {
  name: 'VButton',

  components: {},

  props: {
    variant: {
      type: String,
      default: 'default', // 'default', 'danger', 'warning'
      validator: (value) => ['default', 'danger', 'warning'].includes(value),
    },
  },

  data() {
    return {}
  },

  computed: {
    buttonClasses() {
      const base = 'flex items-center justify-center'
      const variants = {
        default: 'bg-blue-600 text-white hover:bg-blue-700 focus:ring-blue-500',
        danger: 'bg-red-600 text-white hover:bg-red-700 focus:ring-red-500',
        warning: 'bg-yellow-500 text-gray-800 hover:bg-yellow-600 focus:ring-yellow-400',
      }
      return `${base} ${variants[this.variant]}`
    },
  },

  methods: {},

  watch: {},

  created() {},
  mounted() {},
  updated() {},
  beforeUnmount() {},
  unmounted() {},
}
</script>

<template>
  <button
    :class="buttonClasses"
    class="px-4 cursor-pointer py-2 rounded-md font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2"
    @click="$emit('click')"
  >
    <slot></slot>
  </button>
</template>

<style scoped></style>

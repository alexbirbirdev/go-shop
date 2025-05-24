<script>
export default {
  name: 'VFlashMessage',

  components: {},

  props: {
    message: {
      type: [String, Object],
      required: true,
    },
    type: {
      type: String,
      default: 'info',
      validator: (value) => ['success', 'error', 'info'].includes(value),
    },
    duration: {
      type: Number,
      default: 3000,
    },
  },

  data() {
    return {
      isVisible: true,
      timeoutId: null,
    }
  },

  computed: {
    messageClasses() {
      return {
        success: 'bg-green-100 text-green-800',
        error: 'bg-red-100 text-red-800',
        info: 'bg-blue-100 text-blue-800',
      }[this.type]
    },
    iconColor() {
      return {
        success: 'text-green-500',
        error: 'text-red-500',
        info: 'text-blue-500',
      }[this.type]
    },
    displayMessage() {
      return typeof this.message === 'object'
        ? this.message.message || JSON.stringify(this.message)
        : this.message
    },
  },

  methods: {
    setupAutoHide() {
      this.timeoutId = setTimeout(() => {
        this.isVisible = false
      }, this.duration)
    },
    resetTimer() {
      clearTimeout(this.timeoutId)
      this.isVisible = true
      this.setupAutoHide()
    },
  },

  watch: {
    message() {
      this.resetTimer()
    },
  },

  created() {},
  mounted() {
    this.setupAutoHide()
  },
  updated() {},
  beforeUnmount() {},
  unmounted() {},
}
</script>

<template>
  <transition
    enter-active-class="transition-opacity duration-300"
    enter-from-class="opacity-0"
    enter-to-class="opacity-100"
    leave-active-class="transition-opacity duration-300"
    leave-from-class="opacity-100"
    leave-to-class="opacity-0"
  >
    <div
      v-if="isVisible"
      class="fixed top-4 left-1/2 -translate-x-1/2 max-w-xs p-4 rounded-lg shadow-lg"
      :class="messageClasses"
    >
      <div class="flex items-start">
        <div>
          <p class="text-sm font-medium">
            {{ displayMessage }}
          </p>
        </div>
      </div>
    </div>
  </transition>
</template>

<style scoped></style>

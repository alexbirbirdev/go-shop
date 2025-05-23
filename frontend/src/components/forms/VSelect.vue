<script>
export default {
  name: 'VSelect',

  components: {},

  props: {
    options: {
      type: Array,
      required: true,
      validator: (value) => {
        return value.every((option) => {
          return 'value' in option && 'label' in option && 'direction' in option
        })
      },
    },
    modelValue: {
      type: Object,
      required: true,
    },
  },

  data() {
    return {
      isOpen: false,
    }
  },

  computed: {
    activeOption() {
      return (
        this.options.find(
          (opt) =>
            opt.value === this.modelValue.value && opt.direction === this.modelValue.direction,
        ) || this.options[0]
      )
    },
  },

  methods: {
    selectOption(option) {
      this.$emit('update:modelValue', {
        value: option.value,
        direction: option.direction,
      })
      this.isOpen = false
    },
    closeOnClickOutside(e) {
      if (!this.$el.contains(e.target)) {
        this.isOpen = false
      }
    },
  },

  watch: {},

  created() {},
  mounted() {
    document.addEventListener('click', this.closeOnClickOutside)
  },
  updated() {},
  beforeUnmount() {
    document.removeEventListener('click', this.closeOnClickOutside)
  },
  unmounted() {},
}
</script>

<template>
  <div class="relative inline-block text-left">
    <button
      @click="isOpen = !isOpen"
      class="flex items-center gap-1 px-3 py-1 text-sm border rounded-md hover:bg-gray-50"
    >
      <span>{{ activeOption.label }}</span>
    </button>

    <transition
      enter-active-class="transition duration-100 ease-out"
      enter-from-class="transform scale-95 opacity-0"
      enter-to-class="transform scale-100 opacity-100"
      leave-active-class="transition duration-75 ease-in"
      leave-from-class="transform scale-100 opacity-100"
      leave-to-class="transform scale-95 opacity-0"
    >
      <ul
        v-if="isOpen"
        class="absolute overflow-hidden right-0 z-10 mt-1 w-40 origin-top-right rounded-md bg-white shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none"
      >
        <li
          v-for="option in options"
          :key="option.value"
          @click="selectOption(option)"
          :class="activeOption.value === option.value ? 'bg-gray-100 font-bold' : ''"
          class="px-3 py-2 text-sm cursor-pointer hover:bg-gray-100 flex justify-between items-center"
        >
          <span>{{ option.label }}</span>
        </li>
      </ul>
    </transition>
  </div>
</template>

<style scoped></style>

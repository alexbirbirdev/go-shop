export default {
  namespaced: true,
  state: () => ({
    count: 0,
  }),
  mutations: {
    SET_COUNT(state, value) {
      state.count = value
    },
    INCREMENT(state) {
      state.count++
    },
    DECREMENT(state) {
      if (state.count > 0) state.count--
    },
  },
  getters: {
    count: (state) => state.count,
  },
}

import Vuex from 'vuex';
import favorites from './modules/favorites';

export default new Vuex.Store({
    modules: {
        favorites,
    },
    state: {
        isAuth: false,
        token: '',
    },
    mutations: {
        SET_AUTH(state, payload) {
            state.isAuth = payload.isAuth
            state.token = payload.token
        }
    },
    actions: {
        checkAuth({ commit }) {
            const token = localStorage.getItem('token')
            commit('SET_AUTH', !!token, token)
        }
    },
    getters: {
        isAuth: state => state.isAuth
    }
})
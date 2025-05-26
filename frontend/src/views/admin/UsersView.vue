<script>
import axios from 'axios'
import VButton from '@/components/forms/VButton.vue'

export default {
  name: 'AdminUsersView',

  components: { VButton },

  props: {},

  data() {
    return {
      users: [],
      tableHeaders: ['Действия', 'ID', 'Имя', 'Фамилия', 'Почта', 'Роль', 'Создан', 'Обновлен'],
      isLoading: true,
      currentPage: 1,
      itemsPerPage: 20,
      moreAvailable: true,
    }
  },

  computed: {},

  methods: {
    async getUsers(reset = false) {
      try {
        if (reset) {
          this.currentPage = 1
          this.users = []
        }
        this.isLoading = true
        const response = await axios.get('http://localhost:8080/admin/users/', {
          params: {
            page: this.currentPage,
            limit: this.itemsPerPage,
          },
          headers: {
            Authorization: localStorage.getItem('token'),
          },
        })
        this.users.push(...response.data.users)
        this.$router.replace({
          query: {
            page: this.currentPage,
            limit: this.itemsPerPage,
          },
        })
        this.moreAvailable = response.data.users.length === this.itemsPerPage
      } catch (error) {
        console.log(error)
      } finally {
        this.isLoading = false
      }
    },
    loadMoreUsers() {
      if (this.moreAvailable && !this.isLoading) {
        this.currentPage++
        this.getUsers()
      }
    },
  },

  watch: {},

  created() {
    this.currentPage = this.$route.query.page
    this.getUsers()
  },
  mounted() {},
  updated() {},
  beforeUnmount() {},
  unmounted() {},
}
</script>

<template>
  <div class="">
    <!-- Таблица с пользователями -->
    <table class="w-full">
      <thead class="">
        <tr class="text-left">
          <th v-for="header in tableHeaders" :key="header" class="">
            <div class="bg-blue-100 py-4 px-2">{{ header }}</div>
          </th>
        </tr>
      </thead>
      <tbody v-if="!isLoading">
        <tr v-for="user in users" :key="user.ID">
          <td class="">
            <div class="flex items-center gap-2">
              <VButton :variant="'warning'" class="!p-2">
                <svg class="w-5" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path
                    d="M12 2.25034C12.4142 2.25034 12.75 2.58612 12.75 3.00034C12.7498 3.41441 12.4141 3.75034 12 3.75034H11C9.09324 3.75034 7.73859 3.75183 6.71094 3.88999C5.8306 4.00835 5.27638 4.21718 4.86816 4.55112L4.70215 4.70249C4.27903 5.12568 4.0249 5.70531 3.88965 6.71128C3.75154 7.73889 3.75 9.09376 3.75 11.0003V13.0003C3.75 14.9069 3.7515 16.2618 3.88965 17.2894C4.02492 18.2953 4.27898 18.875 4.70215 19.2982L4.86816 19.4496C5.27638 19.7834 5.83068 19.9923 6.71094 20.1107C7.73859 20.2488 9.09327 20.2503 11 20.2503H13C14.9067 20.2503 16.2614 20.2488 17.2891 20.1107C18.295 19.9754 18.8747 19.7213 19.2979 19.2982L19.4492 19.1322C19.7831 18.724 19.992 18.1696 20.1104 17.2894C20.2485 16.2618 20.25 14.9069 20.25 13.0003V12.0003C20.25 11.5861 20.5858 11.2503 21 11.2503C21.4142 11.2503 21.75 11.5861 21.75 12.0003V13.0003C21.75 14.8646 21.7514 16.3388 21.5967 17.4896C21.4489 18.5884 21.1477 19.4911 20.4941 20.2162L20.3584 20.3587C19.6101 21.1069 18.6615 21.4394 17.4893 21.597C16.3384 21.7517 14.8644 21.7503 13 21.7503H11C9.13561 21.7503 7.66158 21.7517 6.51074 21.597C5.41183 21.4492 4.50926 21.148 3.78418 20.4945L3.6416 20.3587C2.89336 19.6105 2.56098 18.6618 2.40332 17.4896C2.2486 16.3388 2.25 14.8646 2.25 13.0003V11.0003C2.25 9.13608 2.24863 7.66188 2.40332 6.51108C2.56095 5.33885 2.8934 4.39022 3.6416 3.64194L3.78418 3.5062C4.50928 2.85261 5.41178 2.55144 6.51074 2.40366C7.66159 2.24893 9.13559 2.25034 11 2.25034H12ZM17.3535 2.62631C18.3976 1.84438 19.8664 1.90052 20.8496 2.79428L21.0303 2.97006L21.2061 3.15073L21.374 3.35385C22.1036 4.32821 22.1037 5.67252 21.374 6.64682L21.2061 6.84995L21.0303 7.03061L17.0977 10.9642C16.441 11.6208 15.9989 12.0684 15.5088 12.391L15.2949 12.5218C14.7128 12.8513 14.0615 13.0077 13.0322 13.265L11.1816 13.7279C10.9262 13.7916 10.6559 13.7167 10.4697 13.5306C10.2836 13.3445 10.2088 13.0741 10.2725 12.8187L10.7354 10.9681C10.9927 9.93886 11.1491 9.28749 11.4785 8.70542L11.6094 8.49155C11.932 8.00141 12.3795 7.55931 13.0361 6.90268L16.9697 2.97006L17.1504 2.79428L17.3535 2.62631ZM14.0977 7.9642C13.3919 8.66998 13.0729 8.99481 12.8662 9.30893L12.7842 9.4437C12.5781 9.80773 12.4671 10.2257 12.1904 11.3324L12.0303 11.9691L12.668 11.8099L13.3789 11.6292C13.9802 11.4717 14.2836 11.3707 14.5566 11.2162L14.6914 11.1341C15.0055 10.9275 15.3304 10.6085 16.0361 9.90268L18.3926 7.54624C17.5069 7.20046 16.7987 6.49258 16.4531 5.60678L14.0977 7.9642ZM19.8408 3.90366C19.3938 3.49735 18.7255 3.47182 18.251 3.82749L18.1592 3.90366L18.0303 4.03061L17.7227 4.33725C17.6984 5.41363 18.5867 6.30021 19.6631 6.27573L19.9697 5.97006L20.0967 5.84116L20.1729 5.74936C20.5048 5.3064 20.5047 4.69433 20.1729 4.25131L20.0967 4.15952L19.9697 4.03061L19.8408 3.90366Z"
                    fill="black"
                  />
                </svg>
              </VButton>
            </div>
          </td>
          <td v-for="values in user" :key="values">
            <div v-if="values">{{ values }}</div>
            <div v-else>-</div>
          </td>
        </tr>
      </tbody>
    </table>
    <div v-if="moreAvailable" class="flex justify-center mt-10">
      <VButton :disabled="isLoading" @click="loadMoreUsers">
        <span v-if="isLoading">Ожидайте</span><span v-else>Загрузить еще</span>
      </VButton>
    </div>
  </div>
</template>

<style scoped>
th {
  border-bottom: 4px solid transparent;
}
th:nth-child(1) {
  width: 1px;
}
th:nth-child(2) {
  width: 1px;
  text-align: center;
}
th:nth-child(3) {
  width: 300px;
}
th:nth-child(4) {
  width: 300px;
}
th:nth-child(5) {
  width: 300px;
}
th:nth-child(6) {
  width: 100px;
}
th:nth-child(7) {
  width: 100px;
  text-align: center;
  white-space: nowrap;
}
th:nth-child(8) {
  text-align: center;
  width: 100px;
}
td:nth-child(7) div,
td:nth-child(8) div {
  white-space: nowrap;
}
th:first-child div {
  border-top-left-radius: 8px;
  border-bottom-left-radius: 8px;
}
th:last-child div {
  border-top-right-radius: 8px;
  border-bottom-right-radius: 8px;
}
td:first-child div {
  border-top-left-radius: 8px;
  border-bottom-left-radius: 8px;
}
td:last-child div {
  border-top-right-radius: 8px;
  border-bottom-right-radius: 8px;
}
td div {
  padding: 8px;
  transition: 0.2s;
  display: flex;
  align-items: center;
  max-width: 100%;
  overflow: hidden;
}
tr:nth-child(2n + 2):not(thead tr) div {
  background: oklch(97% 0 0);
}
tr:hover:not(thead tr) div {
  background: oklch(92.2% 0 0);
}
table,
tr,
td {
  height: 100%;
}
td > div {
  height: 100%;
}
</style>

<template>
  <div class="container">
    <Routes />
    <h1>プロフィール</h1>
    <p>ここはプロフィール画面</p>
    <!-- タブによるページ切り替えのサンプル -->
    <ul class="nav nav-tabs">
      <li class="nav-item" @click="onClickTab(1)">
        <a class="nav-link" :class="{ active: activeTab === 1 }"
          >アカウント情報</a
        >
      </li>
      <li class="nav-item" @click="onClickTab(2)">
        <a class="nav-link" :class="{ active: activeTab === 2 }"
          >作成した質問一覧</a
        >
      </li>
      <li class="nav-item" @click="onClickTab(3)">
        <a class="nav-link" :class="{ active: activeTab === 3 }"
          >回答した質問一覧</a
        >
      </li>
    </ul>
    <!-- 実際のコンテンツ -->
    <p v-if="activeTab === 1">ここはアカウント情報画面</p>
    <p v-if="activeTab === 2">ここは作成した質問一覧画面</p>
    <p v-if="activeTab === 3">ここは回答した質問一覧画面</p>
  </div>
</template>

<script lang="ts">
import Routes from '../components/Routes.vue'
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'Profile',
  components: { Routes },
  data() {
    return {
      activeTab: 1
    }
  },
  methods: {
    onClickTab(num) {
      this.activeTab = num
    }
  }
})
</script>

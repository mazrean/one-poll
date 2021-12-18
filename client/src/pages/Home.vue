<template>
  <div class="container">
    <Routes />
    <h1>ホーム</h1>
    <p>ここはホーム画面</p>
    <!-- 以下にスクリプトを記述 -->
    <!-- chart.jsのテスト -->
    <SampleChart class="chart" :chart-data="chartData" />
    <!-- Bootstrap Iconsのテスト -->
    <h2>
      <em class="bi-alarm"></em>
      <em class="bi bi-book"></em>
      <em class="bi bi-calendar"></em>
      <em class="bi bi-folder"></em>
      <em class="bi bi-forward"></em>
      <em class="bi bi-house-fill"></em>
    </h2>
    <h1>
      <svg
        xmlns="http://www.w3.org/2000/svg"
        width="64"
        height="64"
        fill="currentColor"
        class="bi bi-bar-chart-line-fill"
        viewBox="0 0 16 16">
        <path
          d="M11 2a1 1 0 0 1 1-1h2a1 1 0 0 1 1 1v12h.5a.5.5 0 0 1 0 1H.5a.5.5 0 0 1 0-1H1v-3a1 1 0 0 1 1-1h2a1 1 0 0 1 1 1v3h1V7a1 1 0 0 1 1-1h2a1 1 0 0 1 1 1v7h1V2z" />
      </svg>
      <svg
        xmlns="http://www.w3.org/2000/svg"
        width="128"
        height="128"
        fill="currentColor"
        class="bi bi-chat-left-dots-fill"
        viewBox="0 0 16 16">
        <path
          d="M0 2a2 2 0 0 1 2-2h12a2 2 0 0 1 2 2v8a2 2 0 0 1-2 2H4.414a1 1 0 0 0-.707.293L.854 15.146A.5.5 0 0 1 0 14.793V2zm5 4a1 1 0 1 0-2 0 1 1 0 0 0 2 0zm4 0a1 1 0 1 0-2 0 1 1 0 0 0 2 0zm3 1a1 1 0 1 0 0-2 1 1 0 0 0 0 2z" />
      </svg>
      <svg
        xmlns="http://www.w3.org/2000/svg"
        width="256"
        height="256"
        fill="currentColor"
        class="bi bi-person-fill"
        viewBox="0 0 16 16">
        <path
          d="M3 14s-1 0-1-1 1-4 6-4 6 3 6 4-1 1-1 1H3zm5-6a3 3 0 1 0 0-6 3 3 0 0 0 0 6z" />
      </svg>
    </h1>
  </div>
</template>

<script lang="ts">
import Routes from '../components/Routes.vue'
import SampleChart from '../components/SampleChart.vue'
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'HomePage',
  components: { Routes, SampleChart }
})
</script>

<template>
  <div class="container">
    <Routes />
    <h1>サインイン</h1>
    <p>ここはサインイン画面</p>
    <!-- フォームサンプル -->
    <form action="/hoge" method="POST">
      <div class="my-3">
        <label for="user_name" class="form-label">ユーザー名</label>
        <input
          id="user_name"
          type="text"
          class="form-control"
          name="user_name"
          placeholder="ユーザー名を入力"
        />
        <label for="password" class="form-label">パスワード</label>
        <input
          id="password"
          type="password"
          class="form-control"
          name="password"
          placeholder="パスワードを入力"
        />
      </div>
      <button type="submit" class="btn btn-primary">サインイン</button>
    </form>
  </div>
</template>

<script lang="ts">
import Routes from '../components/Routes.vue'
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'Signin',
  components: { Routes }
})
</script>

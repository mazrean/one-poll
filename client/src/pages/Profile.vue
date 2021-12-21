<template>
  <div class="container">
    <h1>プロフィール</h1>
    <p>ここはプロフィール画面</p>
    <ul id="myTab" class="nav nav-tabs" role="tablist">
      <li class="nav-item" role="presentation">
        <button
          id="home-tab"
          class="nav-link active"
          data-bs-toggle="tab"
          data-bs-target="#home"
          type="button"
          role="tab"
          aria-controls="home"
          aria-selected="true">
          アカウント情報
        </button>
      </li>
      <li class="nav-item" role="presentation">
        <button
          id="profile-tab"
          class="nav-link"
          data-bs-toggle="tab"
          data-bs-target="#profile"
          type="button"
          role="tab"
          aria-controls="profile"
          aria-selected="false">
          作成した質問一覧
        </button>
      </li>
      <li class="nav-item" role="presentation">
        <button
          id="contact-tab"
          class="nav-link"
          data-bs-toggle="tab"
          data-bs-target="#contact"
          type="button"
          role="tab"
          aria-controls="contact"
          aria-selected="false">
          回答した質問一覧
        </button>
      </li>
    </ul>
    <div id="myTabContent" class="tab-content">
      <div
        id="home"
        class="tab-pane fade show active"
        role="tabpanel"
        aria-labelledby="home-tab">
        ここはアカウント情報画面
      </div>
      <div
        id="profile"
        class="tab-pane fade"
        role="tabpanel"
        aria-labelledby="profile-tab">
        ここは作成した質問一覧画面
      </div>
      <div
        id="contact"
        class="tab-pane fade"
        role="tabpanel"
        aria-labelledby="contact-tab">
        ここは回答した質問一覧画面
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'ProfilePage',
  components: {}
})
</script>

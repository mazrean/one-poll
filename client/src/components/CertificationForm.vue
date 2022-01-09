<template>
  <div class="card">
    <div class="card-body">
      <h5 class="card-title">PollQ</h5>
      <form class="needs-validation" @submit="Validation" novalidate>
        <div>
          <label for="user_name" class="card-text">User</label><br />
          <input
            id="user_name"
            type="text"
            class="form-control"
            placeholder="ユーザー名を入力"
            required
            pattern="[0-9a-zA-Z_]{4,16}" />
          <div class="invalid-feedback">
            英数＋アンダーバー込で4~16文字を入力
          </div>
        </div>
        <div>
          <label for="password" class="card-text">password</label><br />
          <input
            id="password"
            type="password"
            class="form-control"
            placeholder="パスワードを入力"
            required
            pattern="[0-9a-zA-Z_]{8,50}" />
          <div class="invalid-feedback">
            英数＋アンダーバー込で8~50文字を入力
          </div>
        </div>
        <button type="submit" class="btn btn-primary">{{ Sign }}</button>
      </form>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'CertificationFormComponent',
  props: ['Sign'],
  methods: {
    Validation: function (event: Event) {
      const forms =
        document.querySelectorAll<HTMLInputElement>('.needs-validation')
      forms.forEach(form => {
        if (!form.checkValidity()) {
          event.preventDefault()
          event.stopPropagation()
        }
        form.classList.add('was-validated')
      }, false)
    }
  },
  setup() {
    return {}
  }
})
</script>

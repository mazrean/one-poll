<template>
    <div class="card">
        <div class="card-body">
            <h5 class="card-title">PollQ</h5>
            <form action="/hoge" method="post">
            <div>
            <label for="user_name" class="card-text">User</label><br>
            <input id="user_name" type="text" placeholder="ユーザー名を入力">
            </div>
            <div>
            <label for="password" class="card-text">password</label><br>
            <input id="password" type="password" placeholder="パスワードを入力">
            </div>
            <button type="submit" class="btn btn-primary">{{ Sign }}</button>
            </form>
        </div>
    </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
    name:'CertificationFormComponent',
    props:['Sign'],
    setup() {
        return {}
    }
})
</script>

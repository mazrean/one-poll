<template>
  <div class="card m-3 p-2" style="border-radius: 1em">
    <div class="card-header text-start bg-white card-title">
      <h5>{{ passkey.name }}</h5>
    </div>
    <div class="card-body">
      <div class="d-flex align-items-center">
        <div class="text-start" style="width: 105px">作成日時</div>
        <time :datetime="passkey.createdAt">{{ createdAtText }}</time>
      </div>
      <div class="d-flex align-items-center">
        <div class="text-start" style="width: 105px">最終使用日時</div>
        <time :datetime="passkey.lastUsedAt">{{ lastUsedAtText }}</time>
      </div>
      <div class="d-flex align-items-center justify-content-end">
        <button type="button" class="btn btn-danger" @click="onDelete">
          削除
        </button>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue'
import { WebAuthnCredential } from '/@/lib/apis'

export default defineComponent({
  props: {
    passkey: {
      type: Object as PropType<WebAuthnCredential>,
      required: true
    }
  },
  emits: ['delete'],
  setup(props, context) {
    const createdAt = new Date(props.passkey.createdAt)
    const createdAtText = `${createdAt.toLocaleDateString()} ${createdAt.toLocaleTimeString()}`

    const lastUsedAt = new Date(props.passkey.lastUsedAt)
    const lastUsedAtText = `${lastUsedAt.toLocaleDateString()} ${lastUsedAt.toLocaleTimeString()}`

    const onDelete = () => {
      context.emit('delete', props.passkey.id)
    }

    return {
      createdAtText,
      lastUsedAtText,
      onDelete
    }
  }
})
</script>

<style>
.card {
}
</style>

<template>
  <!-- Modal -->
  <div
    id="newPollModal"
    class="modal fade"
    tabindex="-1"
    aria-labelledby="newPollModalLabel"
    aria-hidden="true">
    <div class="modal-dialog modal-lg">
      <div class="modal-content">
        <div class="modal-header">
          <h4 class="modal-title">新しい質問</h4>
          <button
            type="button"
            class="btn-close"
            data-bs-dismiss="modal"
            aria-label="Close"></button>
        </div>
        <div class="modal-body">
          <form class="needs-validation" @submit.prevent>
            <label for="title" class="col-sm-2 col-form-label">質問文</label>
            <textarea
              id="title"
              v-model="state.title"
              class="form-control mb-3"
              name="title"
              maxlength="140"
              required
              placeholder="質問する内容を入力" />
            <label for="options" class="col-sm-2 col-form-label">選択肢</label>
            <ul
              v-for="(v, i) in state.options"
              :key="i"
              class="nav nav-pills flex-column">
              <li class="nav-item">
                <div class="d-inline-block d-flex">
                  <input
                    :id="'option-' + i"
                    v-model="state.options[i]"
                    type="text"
                    class="form-control d-flex my-1"
                    :name="'option-' + i"
                    maxlength="25"
                    required
                    placeholder="オプションの内容を入力" />
                  <button
                    :disabled="i < 2"
                    type="button"
                    class="btn btn-close d-flex my-auto ml-3"
                    @Click="deleteOption(i)"></button>
                </div>
              </li>
            </ul>
            <button
              v-show="state.options.length < 5"
              class="btn link"
              type="button"
              @Click="insertOption()">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="32"
                height="32"
                fill="currentColor"
                class="bi bi-plus"
                viewBox="0 0 16 16">
                <path
                  d="M8 4a.5.5 0 0 1 .5.5v3h3a.5.5 0 0 1 0 1h-3v3a.5.5 0 0 1-1 0v-3h-3a.5.5 0 0 1 0-1h3v-3A.5.5 0 0 1 8 4z" />
              </svg>
              <span class="align-middle">選択肢を追加する</span></button
            ><br />

            <label for="tag" class="col-sm-2 col-form-label">タグ</label>
            <div>
              <ul
                v-for="v in state.newTags"
                :key="v"
                class="nav nav-pills d-inline-block">
                <li class="nav-item d-flex">
                  <span
                    class="rounded-3 px-3 py-1 bg-light text-dark d-flex mx-1 my-1 fw-normal"
                    ><em class="bi bi-tags-fill" /> {{ v
                    }}<button
                      type="button"
                      class="btn btn-sm btn-close"
                      @Click="deleteTag(v)"></button
                  ></span>
                </li>
              </ul>
            </div>
            <input
              id="tag"
              v-model="state.newTag"
              :disabled="state.newTags.size >= 10"
              type="text"
              class="form-control mt-2"
              name="tag"
              maxlength="16"
              placeholder="新しく追加するタグ名を入力"
              @Input="calculateFilter()"
              @keydown.enter="insertTag()" />
            <ul v-for="v in state.autocompletes" :key="v" class="list-group">
              <button
                class="list-group-item list-group-item-action p-1"
                type="button"
                @click="onAutocomplete(v.name)">
                <em class="bi bi-tags-fill" /> {{ v.name }}
              </button>
            </ul>
            <label for="deadline" class="col-sm-2 col-form-label">締切</label>
            <div class="d-flex flex-wrap justify-content-center m-auto">
              <select
                id="day"
                v-model="state.day"
                name="day"
                class="form-select form-select-lg mb-3 w-25">
                <template v-for="n of 8" :key="n">
                  <option :value="n - 1">{{ n - 1 }}</option>
                </template>
              </select>
              <span class="m-auto">日</span>
              <select
                id="hour"
                v-model="state.hour"
                name="hour"
                class="form-select form-select-lg mb-3 w-25">
                <template v-for="n of 24" :key="n">
                  <option :value="n - 1">{{ n - 1 }}</option>
                </template>
              </select>
              <span class="m-auto">時間</span>
              <select
                id="minute"
                v-model="state.minute"
                name="minute"
                class="form-select form-select-lg mb-3 w-25">
                <template v-for="n of 60" :key="n">
                  <option :value="n - 1">{{ n - 1 }}</option>
                </template>
              </select>
              <span class="m-auto">分後</span>
            </div>
            <button
              type="button"
              class="btn btn-lg btn-primary"
              :disabled="!validationCheck()"
              data-bs-dismiss="modal"
              @click="submitPoll()">
              <strong>投稿する</strong></button
            ><br />
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, reactive, onMounted } from 'vue'
import api, { NewPoll, PollType, PollTag } from '/@/lib/apis'

interface State {
  title: string
  deadline: string
  day: number
  hour: number
  minute: number
  options: string[]
  newTags: Set<string>
  newTag: string
  tags: PollTag[]
  autocompletes: PollTag[]
  validated: boolean
}

export default defineComponent({
  name: 'NewPollComponent',
  setup() {
    const state = reactive<State>({
      title: '',
      deadline: '0',
      day: 1,
      hour: 0,
      minute: 0,
      options: ['', ''],
      newTags: new Set(),
      newTag: '',
      tags: [],
      autocompletes: [],
      validated: false
    })
    onMounted(async () => {
      state.tags = (await api.getTags()).data
    })
    const insertTag = () => {
      if (state.newTag.length === 0) return
      state.newTags.add(state.newTag)
      state.newTag = ''
    }
    const deleteTag = (str: string) => {
      state.newTags.delete(str)
    }
    const insertOption = () => {
      state.options.push('')
    }
    const deleteOption = (idx: number) => {
      state.options.splice(idx, 1)
    }
    const calculateFilter = () => {
      if (state.newTag.length === 0) {
        state.autocompletes = []
      } else {
        state.autocompletes = state.tags
          .filter((v: PollTag) => {
            return v.name.indexOf(state.newTag) === 0
          })
          .slice(0, 5)
      }
    }
    const onAutocomplete = (str: string) => {
      if (str.length === 0) return
      state.newTags.add(str)
      state.newTag = ''
      calculateFilter()
    }
    const validationCheck = () => {
      if (state.title === '') return false
      for (const opt of state.options) {
        if (opt === '') return false
      }
      return true
    }
    const submitPoll = async () => {
      // deadlineの計算
      const time: Date = new Date()
      time.setDate(time.getDate() + state.day)
      time.setHours(time.getHours() + state.hour)
      time.setMinutes(time.getMinutes() + state.minute)
      const poll: NewPoll = {
        title: state.title,
        type: PollType.Radio,
        deadline: time.toISOString(),
        tags: Array.from(state.newTags),
        question: state.options
      }
      try {
        await api.postPolls(poll)
      } catch {
        alert(
          '質問を投稿できませんでした。時間を空けてもう一度お試しください。'
        )
        return
      }
      state.title = ''
      state.options = ['', '']
      state.newTags = new Set()
      state.newTag = ''
      calculateFilter()
      state.deadline = '0'
      location.href = '/'
    }
    return {
      state,
      insertTag,
      deleteTag,
      insertOption,
      deleteOption,
      calculateFilter,
      onAutocomplete,
      validationCheck,
      submitPoll
    }
  }
})
</script>

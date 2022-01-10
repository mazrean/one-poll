<template>
  <div class="card">
    <div class="card-header text-start">
      <h2 class="card-title">{{ title }}</h2>
      <div class="card-tags bi bi-tags-fill text-muted d-flex flex-wrap">
        <span v-for="tag in tags" :key="tag.id" class="ms-1">
          {{ tag.name }},
        </span>
      </div>
    </div>
    <div class="card-body">
      <div v-if="only_browsable">
        <button
          v-for="q in question"
          :key="q.id"
          type="button"
          class="vote-button btn btn-outline-secondary mb-1"
          disabled>
          {{ q.choice }}
        </button>
      </div>
      <div v-if="can_answer">
        <button
          v-for="q in question"
          :key="q.id"
          type="button"
          class="vote-button btn btn-outline-secondary mb-1">
          {{ q.choice }}
        </button>
      </div>
      <div v-if="can_access_details">
        <PollResultComponent :result="result" />
        <div class="d-flex justify-content-around">
          <div class="card-link"><a href="#">詳細を見る</a></div>
        </div>
      </div>
    </div>
    <div class="footer d-flex justify-content-around">
      <div>投票数 : {{ count }}</div>
      <div>
        作成者 : <a href="#">{{ owner.name }}</a>
      </div>
      <div>作成日 : {{ createdAt }}</div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType } from 'vue'
import PollResultComponent from '/@/components/PollResult.vue'
import result from '/@/assets/poll_result_data.json'

interface Tag {
  id: number
  name: string
}

interface Question {
  id: number
  choice: string
}

interface Owner {
  uuid: string
  name: string
}

interface UserStatus {
  isOwner: boolean
  accessMode: string
}

interface Poll {
  pollId: string
  title: string
  type: string
  deadline: string
  tags: Array<Tag>
  question: Array<Question>
  count: number
  createdAt: string
  owner: Owner
  userStatus: UserStatus
}

export default defineComponent({
  components: { PollResultComponent },
  props: {
    pollId: {
      type: String,
      default: 'poll_id',
      required: true
    },
    title: {
      type: String,
      default: 'Title',
      required: true
    },
    type: {
      type: String,
      default: 'type',
      required: false
    },
    deadline: {
      type: String,
      default: 'deadline',
      required: false
    },
    tags: {
      type: Array as PropType<Tag[]>,
      default() {
        return [
          {
            id: -1,
            name: 'name'
          }
        ]
      },
      required: true
    },
    question: {
      type: Array as PropType<Question[]>,
      default() {
        return [
          {
            id: -1,
            choice: 'choice'
          }
        ]
      },
      required: true
    },
    count: {
      type: Number,
      default: 0,
      required: true
    },
    createdAt: {
      type: String,
      default: 'created_at',
      required: true
    },
    owner: {
      type: Object as PropType<Owner>,
      default() {
        return {
          uuid: 'uuid',
          name: 'name'
        }
      },
      required: true
    },
    userStatus: {
      type: Object as PropType<UserStatus>,
      default() {
        return {
          isOwner: true,
          accessMode: 'only_browsable'
        }
      },
      required: true
    }
  },
  setup(props: Poll) {
    const only_browsable = props.userStatus.accessMode == 'only_browsable'
    const can_answer = props.userStatus.accessMode == 'can_answer'
    const can_access_details =
      props.userStatus.accessMode == 'can_access_details'
    return {
      only_browsable,
      can_answer,
      can_access_details,
      result,
      PollResultComponent
    }
  }
})
</script>

<style>
.card {
  width: 32rem;
}
.vote-button {
  width: 30rem;
  height: 2.7rem;
}
</style>

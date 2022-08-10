<script setup lang="ts">
import Multiselect from '@vueform/multiselect'
import '@vueform/multiselect/themes/default.css'
import { storeToRefs } from 'pinia'
import { useStore } from '../store'
import useWebSocketEvent from '../composables/useWebSocketEvent'
import { nextTick } from 'vue'

const store = useStore()
const { operators, selectedOperator, isWaitingTripUpdates } = storeToRefs(store)
const { clearTripUpdates, waitForTripUpdates } = store
const { sendTripUpdatesRequest } = useWebSocketEvent()
let intervalId: any = null

async function handleOperatorChange() {
  clearInterval(intervalId)
  await nextTick()
  clearTripUpdates()
  // If operator is not selected, do not send request
  if (selectedOperator.value) {
    waitForTripUpdates()
    sendTripUpdatesRequest()

    intervalId = setInterval(() => {
      waitForTripUpdates()
      sendTripUpdatesRequest()
    }, import.meta.env.VITE_TRIP_UPDATE_INTERVAL)
  }
}
</script>

<template>
  <div class="w-[100%]">
    <Multiselect
      v-model="selectedOperator"
      placeholder="Choose the transit operator"
      :options="operators"
      :searchable="true"
      :required="true"
      track-by="name"
      label="name"
      value-prop="id"
      :disabled="operators.length === 0 || isWaitingTripUpdates"
      @change="handleOperatorChange"
    />
  </div>
</template>
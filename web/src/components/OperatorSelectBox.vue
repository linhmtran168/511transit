<script setup lang="ts">
import Multiselect from '@vueform/multiselect'
import '@vueform/multiselect/themes/default.css'
import { storeToRefs } from 'pinia'
import { useStore } from '../store'
import useWebSocketEvent from '../composables/useWebSocketEvent'

const store = useStore()
const { operators, selectedOperator } = storeToRefs(store)
const { clearSelectedOperator } = store
const { sendTripUpdatesRequest } = useWebSocketEvent()
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
      :disabled="operators.length === 0"
      @clear="clearSelectedOperator"
      @change="sendTripUpdatesRequest"
    />
  </div>
</template>
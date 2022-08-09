<script setup lang="ts">
import useWebSocketEvent from './composables/useWebSocketEvent'
import OperatorSelectBox from './components/OperatorSelectBox.vue'
import TripUpdateList from './components/TripUpdateList.vue'

useWebSocketEvent()
</script>

<template>
  <div>
    <div class="container mx-auto">
      <div class="flex flex-col">
        <div class="prose prose-xl mt-64 mx-auto">
          <h1 class="text-primary">511 SF Bay's Transit Update</h1>
        </div>
        <div class="w-[40%] mx-auto mt-10">
          <OperatorSelectBox />
        </div>
        <div class="w-[40%] mx-auto mt-10">
          <TripUpdateList />
        </div>
      </div>
    </div>
  </div>
</template>

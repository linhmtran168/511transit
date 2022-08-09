<script setup lang="ts">
import { TripUpdate } from '../types/global';
import { fromUnixTime, formatISO9075 } from 'date-fns'

const props = defineProps<{
  tripUpdate: TripUpdate
}>()
const trip = props.tripUpdate.trip
const stopTimes = props.tripUpdate.stopTimeUpdate

function formatTime(time: number | undefined | null): string {
  if (!time) {
    return '???'
  }
  return formatISO9075(fromUnixTime(time))
}
</script>
<template>
  <div class="mt-5">
    <div tabindex="0" class="collapse collapse-plus border border-base-300 bg-base-100 rounded-box">
      <input type="checkbox" class="peer" /> 
      <div class="collapse-title text-xl font-medium">
        Trip: {{ trip.tripId }} - Route: {{ trip.routeId }} - Direction: {{ trip.directionId }}
      </div>
      <div class="collapse-content"> 
        <ul class="steps steps-vertical">
          <li class="step" v-for="stop in stopTimes" :key="stop.stopId" :data-content="stop.stopSequence">
            Arrival: {{ formatTime(stop.arrival?.time) }} - Departure: {{ formatTime(stop.departure?.time) }}
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>
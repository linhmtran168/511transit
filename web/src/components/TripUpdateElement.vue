<script setup lang="ts">
import { TripUpdate } from '../types/global';
import useUtils from '../composables/useUtils';

const props = defineProps<{
  tripUpdate: TripUpdate
}>()
const trip = props.tripUpdate.trip
const vehicle = props.tripUpdate.vehicle
const stopTimes = props.tripUpdate.stop_time_update
const { formatUnixTime } = useUtils()

</script>
<template>
  <div class="mt-5">
    <div tabindex="0" class="collapse collapse-plus border border-base-300 bg-base-100 rounded-box">
      <input type="checkbox" class="peer" /> 
      <div class="collapse-title text-xl font-medium">
        Trip: {{ trip.trip_id }} - Route: {{ trip.route_id }}
      </div>
      <div class="collapse-content"> 
        <div v-if="vehicle">
          <b>Vehicle</b>: ID: {{ vehicle.id }} - Label: {{ vehicle.label }} - Plate: {{ vehicle.license_plate }}
        </div>
        <ul class="steps steps-vertical">
          <li class="step" v-for="stop in stopTimes" :key="stop.stop_id" :data-content="stop.stop_sequence">
            <p>
              <b>Arrival:</b> {{ formatUnixTime(stop.arrival?.time) }} - <b>Departure:</b> {{ formatUnixTime(stop.departure?.time) }}
            </p>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>
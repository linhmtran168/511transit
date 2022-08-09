import { defineStore } from 'pinia'
import { MainStore, Operator, TripUpdatesData } from '../types/global'

const state: MainStore = {
  operators: [],
  tripUpdates: [],
  selectedOperator: null,
  isWaitingTripUpdates: false,
}

export const useStore = defineStore('main', {
  state: () => state,
  actions: {
    setOperators(operators: Operator[]) {
      this.operators = operators
    },
    setTripUpdates(data: TripUpdatesData) {
      // If wrong response for other operator, do nothing
      if (data.operatorId !== this.selectedOperator) {
        return
      }
      if (data.tripUpdates) {
        this.tripUpdates = data.tripUpdates    
      } else {
        this.tripUpdates = []
      }
      this.isWaitingTripUpdates = false
    },
    clearTripUpdates() {
      this.tripUpdates = []
    },
    clearSelectedOperator() {
      this.selectedOperator = null
    },
    waitForTripUpdates() {
      this.isWaitingTripUpdates = true
    }
  }
})
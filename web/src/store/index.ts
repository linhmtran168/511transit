import { defineStore } from 'pinia'
import { MainStore, Operator, TripUpdatesData } from '../types/global'

const mainState: MainStore = {
  operators: [],
  tripUpdates: [],
  selectedOperator: null,
  isWaitingTripUpdates: false,
}

export const useStore = defineStore('main', {
  state: () => mainState,
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
      this.isWaitingTripUpdates = false
    },
    clearSelectedOperator() {
      this.selectedOperator = null
    },
    waitForTripUpdates() {
      this.isWaitingTripUpdates = true
    },
    clearWaitingState() {
      this.isWaitingTripUpdates = false
    }
  }
})
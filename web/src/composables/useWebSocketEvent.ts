import { useStore } from "../store"
import { MessageResponse, MessageType, Operator, MessageRequest, TripUpdatesData } from "../types/global"
import useWebSocketConnection from "./useWebSocketConnection"
import camelcaseKeys from 'camelcase-keys';

export default function useWebSocketEvent() {
  const { initConnection } = useWebSocketConnection()
  const conn = initConnection()
  const store = useStore()
  const { setOperators, setTripUpdates, clearTripUpdates, waitForTripUpdates } = store

  if (!conn.onmessage) {
    conn.onmessage = (event: MessageEvent) => {
      const resp = JSON.parse(event.data) as MessageResponse
      if (resp.responseType === MessageType.operators) {
        setOperators(resp.data as Operator[])
      } else {
        const tripUpdateDatas = camelcaseKeys(resp.data, { deep: true }) as TripUpdatesData
        setTripUpdates(tripUpdateDatas)
      }
    }
  }

  if (!conn.onopen) {
    conn.onopen = () => sendOperatorsRequest()
  }

  const sendTripUpdatesRequest = (operatorId: string) => {
    clearTripUpdates()
    const request: MessageRequest = { requestType: MessageType.tripUpdates, data: { operatorId: operatorId } }
    waitForTripUpdates()
    conn.send(JSON.stringify(request))
  }

  const sendOperatorsRequest = () => {
    const request: MessageRequest = { requestType: MessageType.operators }
    conn.send(JSON.stringify(request))
  }

  return {
    sendTripUpdatesRequest
  }
}
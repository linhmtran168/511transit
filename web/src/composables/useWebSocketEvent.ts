import { useStore } from "../store"
import { MessageResponse, MessageType, Operator, MessageRequest, TripUpdatesData } from "../types/global"
import useWebSocketConnection from "./useWebSocketConnection"

export default function useWebSocketEvent() {
  const { initConnection } = useWebSocketConnection()
  const conn = initConnection()
  const store = useStore()
  const { setOperators, setTripUpdates } = store

  if (!conn.onmessage) {
    conn.onmessage = (event: MessageEvent) => {
      const resp = JSON.parse(event.data) as MessageResponse

      switch (resp.responseType) {
        case MessageType.operators:
          setOperators(resp.data as Operator[])
          break
        case MessageType.tripUpdates:
          setTripUpdates(resp.data as TripUpdatesData)
          break
        default:
          return
      }
    }
  }

  if (!conn.onopen) {
    conn.onopen = () => sendOperatorsRequest()
  }

  const sendTripUpdatesRequest = () => {
    if (!store.selectedOperator) {
      return
    }
    const request: MessageRequest = { requestType: MessageType.tripUpdates, data: { operatorId: store.selectedOperator } }
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
import { useStore } from "../store"
import { MessageResponse, MessageType, Operator, MessageRequest, TripUpdatesData } from "../types/global"
import useWebSocketConnection from "./useWebSocketConnection"

export default function useWebSocketEvent() {
  const { initConnection, sendMessage, closeConnection, setupEventHandlers } = useWebSocketConnection()
  const store = useStore()
  const { setOperators, setTripUpdates } = store

  const sendTripUpdatesRequest = () => {
    if (!store.selectedOperator) {
      return
    }
    const request: MessageRequest = { requestType: MessageType.tripUpdates, data: { operatorId: store.selectedOperator } }
    sendMessage(JSON.stringify(request))
  }

  const sendOperatorsRequest = () => {
    const request: MessageRequest = { requestType: MessageType.operators }
    sendMessage(JSON.stringify(request))
  }

  const openHandler = () => sendOperatorsRequest()

  const messageHandler = (event: MessageEvent) => {
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

  const errorHandler = () => {
    // close the connection
    closeConnection()
  } 

  const closeHandler = () => {
    // reconnect if closed
    setTimeout(() => {
      console.log('Connection close, reconnecting...')
      initConnection()
      setupEventHandlers(openHandler, messageHandler, errorHandler, closeHandler)
      store.clearWaitingState()
    }, 1000)
  }

  initConnection()
  // Initial setup event handlers
  setupEventHandlers(openHandler, messageHandler, errorHandler, closeHandler)

  return {
    sendTripUpdatesRequest
  }
}
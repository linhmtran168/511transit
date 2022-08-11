let connection: WebSocket

function initConnection() {
  if (!connection || connection.readyState === WebSocket.CLOSED) {
    let protocol = window.location.protocol === 'http:' ? 'ws' : 'wss'
    let host = `${protocol}://${window.location.host}/ws`
    if (import.meta.env.DEV) {
      host = import.meta.env.VITE_LOCAL_WS_HOST
    }
    connection = new WebSocket(host, ['realtime-transit'])
  }
}

function closeConnection() {
  connection.close()
}

function setupEventHandlers(
  openHandler: () => void,
  messageHandler: (event: MessageEvent) => void,
  errorHandler: () => void,
  closeHandler: () => void) {
    if (!connection || connection.readyState === WebSocket.CLOSED) {
      return
    }

    if (!connection.onopen) {
      connection.onopen = openHandler
    }
    if (!connection.onmessage) {
      connection.onmessage = messageHandler
    }
    if (!connection.onerror) {
      connection.onerror = errorHandler
    }
    if (!connection.onclose) {
      connection.onclose = closeHandler
    }
}

function sendMessage(message: string) {
  connection.send(message)
}


export default function useWebSocketConnection() {
  return {
    initConnection,
    sendMessage,
    closeConnection,
    setupEventHandlers,
  }
}
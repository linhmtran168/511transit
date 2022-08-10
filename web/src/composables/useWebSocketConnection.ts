let connection: WebSocket

function initConnection() {
  if (!connection || connection.readyState === WebSocket.CLOSED) {
    let host = `ws://${window.location.host}/ws`
    if (import.meta.env.DEV) {
      host = import.meta.env.VITE_LOCAL_WS_HOST
    }
    connection = new WebSocket(host, ['realtime-transit'])
  }

  return connection
}

export default function useWebSocketConnection() {
  return {
    initConnection
  }
}
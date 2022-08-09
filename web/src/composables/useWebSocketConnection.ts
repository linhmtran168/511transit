let connection: WebSocket

function initConnection() {
  if (!connection || connection.readyState === WebSocket.CLOSED) {
    connection = new WebSocket("ws://localhost:5050/ws", ['realtime-transit'])
  }

  return connection
}

export default function useWebSocketConnection() {
  return {
    initConnection
  }
}
// src/wsClient.js
export function createWS(onEvent) {
  const { protocol, host } = window.location;
  const url =
    (protocol === "https:" ? "wss://" : "ws://") +
    host.replace(/:\d+$/, ":8080") +
    "/ws";

  const ws = new WebSocket(url);

  ws.onopen = () => console.log("✅ WebSocket conectado:", url);
  ws.onmessage = (ev) => {
    try {
      const data = JSON.parse(ev.data);
      if (onEvent) onEvent(data);
    } catch (err) {
      console.error("Erro ao decodificar mensagem WebSocket:", err);
    }
  };
  ws.onclose = () => console.log("🔴 WebSocket desconectado");
  ws.onerror = (e) => console.error("⚠️ Erro WebSocket:", e);

  return {
    send: (event, data) => ws.send(JSON.stringify({ event, data })),
    raw: ws,
  };
}

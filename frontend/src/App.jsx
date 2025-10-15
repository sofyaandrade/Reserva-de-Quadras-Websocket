import React, { useEffect, useState } from "react";
import { createWS } from "./wsClient";
import "./styles/palette.css"; // importa o tema pastel

function nowUnix() {
  return Math.floor(Date.now() / 1000);
}

export default function App() {
  const [wsClient, setWsClient] = useState(null);
  const [courts, setCourts] = useState([]);
  const [reservations, setReservations] = useState([]);
  const [log, setLog] = useState([]);
  const [systemStarted, setSystemStarted] = useState(false);
  const [newCourt, setNewCourt] = useState({ name: "", capacity: 1 });
  const [timers, setTimers] = useState({});

  useEffect(() => {
    const client = createWS(handleEvent);
    setWsClient(client);
    return () => {
      if (client?.raw) client.raw.close();
    };
  }, []);

  useEffect(() => {
    const interval = setInterval(() => {
      setTimers((t) =>
        Object.fromEntries(
          Object.entries(t)
            .map(([id, v]) => [id, { ...v, seconds: Math.max(0, v.seconds - 1) }])
            .filter(([_, v]) => v.seconds > 0)
        )
      );
    }, 1000);
    return () => clearInterval(interval);
  }, []);

  function addLog(msg) {
    setLog((l) => [msg, ...l].slice(0, 50));
  }

  function handleEvent(ev) {
    const evt = ev.event;
    const data = ev.data;
    addLog(`${evt} — ${JSON.stringify(data).slice(0, 100)}`);

    switch (evt) {
      case "state.snapshot":
        setCourts(data.courts || []);
        setReservations(data.reservations || []);
        setSystemStarted(data.systemStarted || false);
        break;

      case "sistema.iniciado":
        setCourts(data.courts || []);
        setSystemStarted(true);
        addLog("✅ Sistema iniciado.");
        break;

      case "sistema.parado":
        setSystemStarted(false);
        addLog("🚫 Sistema parado pelo administrador.");
        break;

      case "quadra.adicionada":
        setCourts((c) => [data, ...c]);
        addLog(`🏟️ Nova quadra adicionada: ${data.name}`);
        break;

      case "horario.reservado":
        setReservations((r) => [data, ...r.filter((x) => x.id !== data.id)]);
        setTimers((t) => ({ ...t, [data.id]: { type: "confirm", seconds: 5 } }));
        break;

      case "horario.confirmado":
        setReservations((r) => [data, ...r.filter((x) => x.id !== data.id)]);
        setTimers((t) => ({ ...t, [data.id]: { type: "use", seconds: 300 } }));
        break;

      case "reserva.cancelada":
      case "reserva.expirada":
      case "jogo.finalizado":
        setReservations((r) => [data, ...r.filter((x) => x.id !== data.id)]);
        setTimers((t) => {
          const newT = { ...t };
          delete newT[data.id];
          return newT;
        });
        break;

      case "reserva.negada": {
        const motivo = data.reason || "Motivo desconhecido";
        addLog(`🚫 Reserva negada: ${motivo}`);
        const fakeReservation = {
          id: "negada-" + Date.now(),
          courtId: data.courtId || "-",
          user: data.user || "—",
          status: "Negada",
          motivo,
        };
        setReservations((r) => [fakeReservation, ...r]);
        break;
      }
      default:
        break;
    }
  }

  // ==== AÇÕES ====
  const startSystem = () => wsClient.send("sistema.iniciado", {});
  const stopSystem = () => {
    if (window.confirm("Deseja realmente parar o sistema?"))
      wsClient.send("sistema.parado", {});
  };
  const addCourt = () => {
    if (!newCourt.name.trim()) return alert("Informe o nome da quadra!");
    wsClient.send("quadra.adicionada", {
      name: newCourt.name,
      capacity: Number(newCourt.capacity) || 1,
    });
    setNewCourt({ name: "", capacity: 1 });
  };
  const reserve = (courtId) => {
    const user = prompt("Nome do usuário para reserva:", "Aluno");
    if (!user) return;
    wsClient.send("horario.reservado", { courtId, user, startTime: nowUnix() });
  };
  const cancel = (resId) => wsClient.send("reserva.cancelada", { reservationId: resId });
  const finalize = (resId) => wsClient.send("jogo.finalizado", { reservationId: resId });

  function formatTime(s) {
    const m = Math.floor(s / 60);
    const sec = s % 60;
    return `${String(m).padStart(2, "0")}:${String(sec).padStart(2, "0")}`;
  }

  return (
    <>
      <header className="header">
        <div className="header-inner page">
          <div className="brand">
            <span className="brand-dot" />
            <span>Reserva de Quadras</span>
          </div>
          <span className="badge info">Sistema em tempo real</span>
        </div>
      </header>

      <main className="page">
        {!systemStarted && (
          <div className="panel warn" style={{ marginBottom: 16 }}>
            ⚠️ O sistema ainda não foi iniciado. As reservas estão bloqueadas.
          </div>
        )}

        <div className="columns">
          {/* === Área do Usuário === */}
          <section className="panel user">
            <div className="section-title">
              <span>Usuário</span>
              <span className="badge ok">
                {systemStarted ? "Ativo" : "Aguardando início"}
              </span>
            </div>

            <h3 style={{ marginTop: 0 }}>Quadras Disponíveis</h3>
            {courts.length === 0 && <p>Nenhuma quadra cadastrada.</p>}
            <ul>
              {courts.map((c) => (
                <li key={c.id} style={{ marginBottom: 8 }}>
                  <b>{c.name}</b> — Capacidade: {c.capacity}
                  <button
                    className="btn btn-primary"
                    style={{ marginLeft: 10 }}
                    onClick={() => reserve(c.id)}
                    disabled={!systemStarted}
                  >
                    Reservar
                  </button>
                </li>
              ))}
            </ul>

            <h3>Reservas</h3>
            {reservations.length === 0 && <p>Nenhuma reserva até o momento.</p>}
            <div className="panel" style={{ background: "var(--panel-2)" }}>
              <table className="table">
                <thead>
                  <tr>
                    <th>ID</th>
                    <th>Quadra</th>
                    <th>Usuário</th>
                    <th>Status</th>
                    <th>Motivo</th>
                    <th>Tempo</th>
                    <th>Ações</th>
                  </tr>
                </thead>
                <tbody>
                  {reservations.map((r) => {
                    const timer = timers[r.id];
                    let timeLabel = "";
                    if (timer) {
                      if (timer.type === "confirm")
                        timeLabel = `Confirma em ${formatTime(timer.seconds)}`;
                      else if (timer.type === "use")
                        timeLabel = `Uso: ${formatTime(timer.seconds)}`;
                    }
                    return (
                      <tr key={r.id || Math.random()}>
                        <td>{r.id}</td>
                        <td>{r.courtId}</td>
                        <td>{r.user}</td>
                        <td>{r.status}</td>
                        <td>{r.motivo || "-"}</td>
                        <td>{timeLabel || "-"}</td>
                        <td>
                          {(r.status === "reserved" || r.status === "confirmed") && (
                            <button className="btn" onClick={() => cancel(r.id)}>
                              Cancelar
                            </button>
                          )}
                          {r.status === "confirmed" && (
                            <button className="btn btn-primary" onClick={() => finalize(r.id)}>
                              Finalizar
                            </button>
                          )}
                        </td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>
          </section>

          {/* === Área do Admin === */}
          <section className="panel admin">
            <div className="section-title">
              <span>Administração</span>
              <span className="badge warn">Painel</span>
            </div>

            <div style={{ display: "flex", gap: 8, marginBottom: 14, flexWrap: "wrap" }}>
              <button className="btn btn-primary" onClick={startSystem}>
                📡 Iniciar sistema
              </button>
              <button className="btn" onClick={stopSystem}>
                🛑 Parar sistema
              </button>
            </div>

            <div className="panel" style={{ background: "var(--panel-2)", marginBottom: 20 }}>
              <h4 style={{ marginTop: 0 }}>Adicionar nova quadra</h4>
              <div style={{ display: "flex", gap: 10, alignItems: "center", flexWrap: "wrap" }}>
                <input
                  className="input"
                  type="text"
                  placeholder="Nome da quadra"
                  value={newCourt.name}
                  onChange={(e) => setNewCourt({ ...newCourt, name: e.target.value })}
                />
                <input
                  className="input"
                  type="number"
                  min="1"
                  placeholder="Capacidade"
                  style={{ width: 120 }}
                  value={newCourt.capacity}
                  onChange={(e) => setNewCourt({ ...newCourt, capacity: e.target.value })}
                />
                <button className="btn btn-primary" onClick={addCourt}>
                  ➕ Adicionar
                </button>
              </div>
            </div>

            <div className="panel" style={{ background: "var(--panel-2)" }}>
              <div className="section-title" style={{ marginBottom: 6 }}>
                <span>Logs do Sistema</span>
              </div>
              <div
                style={{
                  maxHeight: 400,
                  overflow: "auto",
                  background: "rgba(255,255,255,.03)",
                  padding: 8,
                  fontSize: 13,
                }}
              >
                {log.map((l, i) => (
                  <div key={i}>{l}</div>
                ))}
              </div>
            </div>
          </section>
        </div>
      </main>

      <footer className="footer">

        <div>
            Desenvolvido por: 
        <br /><b>Matheus Ferrari dos Santos</b>
        <br />&
        <br /><b>Sofya BBS de Andrade</b>      
        </div>
        
      </footer>
    </>
  );
}

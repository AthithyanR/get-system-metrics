import { useEffect, useRef, useState } from 'react'
import './App.css'

function App() {
  const [stats, setStats] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [disConnected, setDisconnected] = useState(false);
  const firstMount = useRef(true);

  // const wsUrl = "ws://localhost/ws";
  const wsUrl = "ws://athi.fun/ws";

  useEffect(() => {
    if (firstMount.current) {
      firstMount.current = false;
      connect(wsUrl);
    }
  }, []);

  function connect(url: string): void {
    const ws = new WebSocket(url);
    ws.onopen = onOpen;
    ws.onmessage = onMessage;
    ws.onerror = () => { };
    ws.onclose = onClose;
  }

  function onOpen(): void {
    setLoading(false);
  }

  function onMessage(event: any): void {
    setStats(JSON.parse(event.data))
  }

  function onClose(event: any): void {
    setLoading(false);
    setDisconnected(true);
  }

  if (loading || !stats) return <>Loading...</>
  if (disConnected) return <>Socket disconnected...</>

  return (
    <div className="App">
      <h2 className="title">CPU</h2>
      <div className="container">
        <div className="card" style={{ backgroundColor: stats.CpuUser < 80 ? 'green' : 'red' }}>
          <label>User</label>
          <p>{stats.CpuUser.toFixed(2)}%</p>
        </div>
        <div className="card" style={{ backgroundColor: stats.CpuSystem < 80 ? 'green' : 'red' }}>
          <label>System</label>
          <p>{stats.CpuSystem.toFixed(2)}%</p>
        </div>
        <div className="card" style={{ backgroundColor: stats.CpuIdle > 80 ? 'green' : 'red' }}>
          <label>Idle</label>
          <p>{stats.CpuIdle.toFixed(2)}%</p>
        </div>
      </div>
      <h2 className="title">MEMORY</h2>
      <div className="container">

        <div className="card" style={{ backgroundColor: 'green' }}>
          <label>Total</label>
          <p>{Math.floor(stats.MemoryTotal)} MB</p>
        </div>
        <div
          className="card"
          style={{
            backgroundColor: stats.MemoryUsed < (stats.MemoryTotal * 0.8) ? 'green' : 'red',
          }}
        >
          <label>Used</label>
          <p>{Math.floor(stats.MemoryUsed)} MB</p>
        </div>
        <div 
          className="card"
          style={{
            backgroundColor: stats.MemoryCached < (stats.MemoryTotal * 0.8) ? 'green' : 'red',
          }}
        >
          <label>Cached</label>
          <p>{Math.floor(stats.MemoryCached)} MB</p>
        </div>
        <div
          className="card"
          style={{
            backgroundColor: stats.MemoryFree > (stats.MemoryTotal * 0.8) ? 'green' : 'red',
          }}
        >
          <label>Free</label>
          <p>{Math.floor(stats.MemoryFree)} MB</p>
        </div>
      </div>
    </div>
  )
}

export default App

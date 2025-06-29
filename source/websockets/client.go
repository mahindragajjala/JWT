<!DOCTYPE html>
<html>
<body>
  <h3>WebSocket JWT Demo</h3>
  <script>
    async function connect() {
      // Get token
      const user = "mahindra";
      const token = await fetch(`http://localhost:8080/login?user=${user}`).then(res => res.text());

      // Connect to WebSocket
      const ws = new WebSocket(`ws://localhost:8080/ws?token=${token}`);

      ws.onopen = () => {
        console.log("Connected");
        ws.send("Hello from client!");
      };

      ws.onmessage = (msg) => {
        console.log("Received:", msg.data);
      };

      ws.onerror = (e) => console.error("WebSocket error:", e);
    }

    connect();
  </script>
</body>
</html>

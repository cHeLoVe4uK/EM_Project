export class WebSocketClient {
    private socket: WebSocket | null = null;
    private url: string;
  
    constructor(url: string) {
      this.url = url;
    }
  
    connect(onMessage: (data: any) => void) {
      this.socket = new WebSocket(this.url);
  
      this.socket.onopen = () => {
        console.log("WebSocket connected");
      };
  
      this.socket.onmessage = (event) => {
        const data = JSON.parse(event.data);
        onMessage(data);
      };
  
      this.socket.onerror = (error) => {
        console.error("WebSocket error:", error);
      };
  
      this.socket.onclose = () => {
        console.log("WebSocket disconnected, attempting to reconnect...");
        setTimeout(() => this.connect(onMessage), 3000); // Переподключение
      };
    }
  
    send(data: any) {
      if (this.socket?.readyState === WebSocket.OPEN) {
        this.socket.send(JSON.stringify(data));
      } else {
        console.error("WebSocket is not open");
      }
    }
  
    disconnect() {
      this.socket?.close();
    }
  }
  
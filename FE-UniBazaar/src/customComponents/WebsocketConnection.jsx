import { useEffect, useRef } from 'react'; 

const useWebSocket = (userId, onMessageReceived) => {
    const ws = useRef(null);
  
    useEffect(() => {
      if (!userId) return;
  
      const newWs = new WebSocket(`ws://13.218.174.66:8080/ws?user_id=${userId}`);
      ws.current = newWs;
  
      newWs.onopen = () => {
        console.log(`Connected as User ${userId}`);
      };
  
      newWs.onmessage = (event) => {
        const receivedMessage = JSON.parse(event.data);
        console.log("Received message:", receivedMessage);
        onMessageReceived(receivedMessage);
      };
  
      newWs.onerror = (error) => {
        console.error("WebSocket Error:", error);
      };
  
      newWs.onclose = (event) => {
        console.log("WebSocket closed:", event);
      };
  
      return () => {
        if (ws.current && ws.current.readyState === WebSocket.OPEN) {
          console.log("Closing WebSocket...");
          ws.current.close();
        }
      };
    }, [userId, onMessageReceived]);
  
    return ws;
  };

export default useWebSocket; 

import { useEffect, useRef } from 'react'; 
const CHAT_USERS_WS_URL = import.meta.env.VITE_CHAT_USERS_WS_URL;


const useWebSocket = (userId, onMessageReceived) => {
    const ws = useRef(null);
  
    useEffect(() => {
      if (!userId) return;
  
      const newWs = new WebSocket(`${CHAT_USERS_WS_URL}/ws?user_id=${userId}`);
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

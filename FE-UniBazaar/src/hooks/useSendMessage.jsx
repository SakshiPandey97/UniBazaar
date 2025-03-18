import { useCallback } from "react";
import { v4 as uuidv4 } from "uuid";

const useSendMessage = (userId, selectedUser, users, ws, input, setInput, setMessages) => {
  const sendMessage = useCallback(async () => {
    if (!userId || !selectedUser) {
      alert("Please select a user to chat with!");
      return;
    }

    if (input.trim() === "") {
      console.error("Message is empty.");
      return;
    }

    let sender_name = "";
    const sender = users.find((u) => u.id === parseInt(userId));
    if (sender) sender_name = sender.name;

    const tempMessageId = uuidv4();
    const message = {
      ID: tempMessageId,
      sender_id: parseInt(userId),
      receiver_id: selectedUser.id,
      content: input,
      timestamp: Date.now(),
      read: false,
      sender_name: sender_name,
    };

    try {
      if (!ws.current || ws.current.readyState !== WebSocket.OPEN) {
        console.error("WebSocket not ready.");
        return;
      }
      ws.current.send(JSON.stringify(message));
      setInput("");
    } catch (e) {
      console.error("Error sending message:", e);
    }
  }, [userId, selectedUser, users, ws, input, setInput]);

  return sendMessage;
};

export default useSendMessage;

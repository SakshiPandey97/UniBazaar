import React, { useState, useEffect, useRef } from "react";
import { getCurrentUserId } from "../utils/getUserId";
import { getAllUsersAPI } from "../api/axios";

const Chat = () => {
  const [users, setUsers] = useState([]);
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState("");
  const [selectedUser, setSelectedUser] = useState(null);
  const [typing, setTyping] = useState(false);
  const [loadingUsers, setLoadingUsers] = useState(true);

  const ws = useRef(null);
  const usersRef = useRef([]);
  const messagesEndRef = useRef(null);
  const typingTimeoutRef = useRef(null);

  const userId = getCurrentUserId();

  // Fetch users on mount
  useEffect(() => {
    if (userId) {
      const fetchUsers = async () => {
        try {
          const data = await getAllUsersAPI(userId);
          setUsers(data);
          usersRef.current = data; // Keep users in ref to avoid stale state
        } catch (error) {
          console.error("Failed to fetch users:", error);
          setTimeout(fetchUsers, 5000);
        } finally {
          setLoadingUsers(false);
        }
      };

      fetchUsers();
    }
  }, [userId]);

  // Initialize WebSocket once
  useEffect(() => {
    if (!ws.current && userId) {
      ws.current = new WebSocket(`ws://localhost:8080/ws?user_id=${userId}`);

      ws.current.onopen = () => {
        console.log(`Connected as User ${userId}`);
      };

      ws.current.onmessage = (event) => {
        const receivedMessage = JSON.parse(event.data);
        console.log("Received message:", receivedMessage);

        setMessages((prevMessages) => {
          const isDuplicate = prevMessages.some(
            (msg) => msg.message_id === receivedMessage.id
          );
          if (isDuplicate) return prevMessages;

          // Use latest users list
          const senderUser = usersRef.current.find(
            (u) => u.id === receivedMessage.sender_id
          );
          const senderName =
            receivedMessage.sender_id === userId ? "You" : senderUser?.name || "Unknown";

          return [
            ...prevMessages,
            {
              message_id: receivedMessage.id,
              sender_id: receivedMessage.sender_id,
              receiver_id: receivedMessage.receiver_id,
              message_text: receivedMessage.content, // Fix blank message issue
              timestamp: receivedMessage.timestamp,
              read: receivedMessage.read,
              sender_name: senderName,
            },
          ];
        });

        scrollToBottom();
      };

      ws.current.onerror = (error) => {
        console.error("WebSocket Error:", error);
      };

      ws.current.onclose = () => {
        console.log("WebSocket closed.");
      };

      return () => {
        if (ws.current && ws.current.readyState === WebSocket.OPEN) {
          console.log("Closing WebSocket...");
          ws.current.close();
        }
      };
    }
  }, [userId]);

  // Fetch messages when a user is selected
  useEffect(() => {
    if (selectedUser && userId) {
      const fetchMessages = async () => {
        try {
          const url = `http://localhost:8080/api/conversation/${userId}/${selectedUser.id}`;
          const res = await fetch(url);

          if (!res.ok) {
            const errorMessage = `HTTP error! Status: ${res.status}`;
            throw new Error(errorMessage);
          }

          const data = await res.json();
          setMessages(data);
        } catch (error) {
          console.error("Failed to load messages:", error);
        }
      };

      fetchMessages();
    }
  }, [selectedUser, userId]);

  // Scroll to bottom when messages update
  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  // Send message and optimistically update UI
  const sendMessage = async () => {
    if (!userId || !selectedUser) {
      alert("Please select a user to chat with!");
      return;
    }

    if (input.trim() === "") {
      console.error("Message is empty.");
      return;
    }

    const message = {
      sender_id: parseInt(userId),
      receiver_id: selectedUser.id,
      content: input,
      timestamp: Date.now(),
      read: false,
    };

    // Optimistically update UI before sending
    setMessages((prevMessages) => [
      ...prevMessages,
      { 
        ...message, 
        sender_name: "You", 
        message_id: `temp-${Date.now()}`,
        message_text: message.content, // Ensure message text is present
      },
    ]);

    setInput("");
    scrollToBottom();

    try {
      if (!ws.current || ws.current.readyState !== WebSocket.OPEN) {
        console.error("WebSocket not ready.");
        return;
      }
      ws.current.send(JSON.stringify(message));
    } catch (e) {
      console.error("Error sending message:", e);
    }
  };

  const handleTyping = (e) => {
    setInput(e.target.value);
    if (!typing) setTyping(true);

    if (typingTimeoutRef.current) clearTimeout(typingTimeoutRef.current);
    typingTimeoutRef.current = setTimeout(() => setTyping(false), 2000);
  };

  const scrollToBottom = () => {
    if (messagesEndRef.current) {
      messagesEndRef.current.scrollIntoView({ behavior: "smooth" });
    }
  };

  return (
    <div className="flex h-screen">
      <div className="w-1/4 bg-gray-200 p-4">
        <h3 className="font-bold mb-4">Users</h3>
        {loadingUsers ? (
          <div>Loading users...</div>
        ) : (
          <ul>
            {users.map((user) => (
              <li
                key={user.id}
                onClick={() => setSelectedUser(user)}
                className={`cursor-pointer p-2 mb-2 rounded ${
                  selectedUser?.id === user.id
                    ? "bg-gray-400"
                    : "hover:bg-gray-300"
                }`}
              >
                {user.name}
              </li>
            ))}
          </ul>
        )}
      </div>

      <div className="w-3/4 p-4">
        {selectedUser ? (
          <div>
            <h2 className="text-xl font-bold mb-3">
              Chat with {selectedUser.name}
            </h2>

            <div className="h-96 overflow-y-auto border p-2 mb-3 bg-gray-100 rounded">
              {messages.map((msg) => (
                <div
                  key={msg.message_id || `${msg.sender_id}-${msg.timestamp}`}
                  className={`my-1 mb-2`}
                >
                  <div
                    className={`p-2 rounded-lg inline-block ${
                      msg.sender_id !== selectedUser.id
                        ? "bg-blue-200 text-right float-right mb-1"
                        : "bg-gray-300 float-left mb-1"
                    }`}
                    style={{ clear: "both" }}
                  >
                    <strong>
                      {msg.sender_id !== selectedUser.id ? "You" : msg.sender_name}
                    </strong>
                    : {msg.message_text}
                  </div>
                </div>
              ))}
              <div ref={messagesEndRef} />
            </div>

            {typing && <p className="text-sm text-gray-500">Typing...</p>}

            <div className="flex gap-2">
              <input
                type="text"
                value={input}
                onChange={handleTyping}
                className="flex-1 border p-2 rounded"
                placeholder="Type a message..."
              />
              <button onClick={sendMessage} className="bg-blue-500 text-white p-2 rounded">
                Send
              </button>
            </div>
          </div>
        ) : (
          <p className="text-center text-gray-500">Select a user to chat with</p>
        )}
      </div>
    </div>
  );
};

export default Chat;

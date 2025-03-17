import React, { useState, useEffect, useRef } from "react";
import { getCurrentUserId } from "../utils/getUserId";
import { getAllUsersAPI } from "../api/axios";
import { v4 as uuidv4 } from 'uuid';

const Chat = () => {
  const [users, setUsers] = useState([]);
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState("");
  const [selectedUser, setSelectedUser] = useState(null);
  const [typing, setTyping] = useState(false);
  const [loadingUsers, setLoadingUsers] = useState(true);

  const ws = useRef(null);
  const messagesEndRef = useRef(null);
  const typingTimeoutRef = useRef(null);

  const userId = getCurrentUserId();

  useEffect(() => {
    if (userId) {
      const fetchUsers = async () => {
        try {
          const data = await getAllUsersAPI(userId);
          setUsers(data);
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

  useEffect(() => {
    if (userId) {
      const newWs = new WebSocket(`ws://localhost:8080/ws?user_id=${userId}`);
      ws.current = newWs;

      newWs.onopen = async () => {
        console.log(`Connected as User ${userId}`);
      };

      newWs.onmessage = (event) => {
        const receivedMessage = JSON.parse(event.data);
        console.log("Received message:", receivedMessage);

        setMessages((prevMessages) => {
            return [...prevMessages, receivedMessage];
        });
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
    }
  }, [userId, users]);

  useEffect(() => {
    if (selectedUser && userId) {
      const fetchMessages = async () => {
        try {
          const url = `http://localhost:8080/api/conversation/${userId}/${selectedUser.id}`;
          console.log("Fetching messages from:", url);
          console.log("Selected User:", selectedUser);
          console.log("User ID:", userId);
          console.log("Receiver ID", selectedUser.id);

          const res = await fetch(url);

          if (!res.ok) {
            const errorData = await res.json().catch(() => null);
            const errorMessage =
              errorData?.message || `HTTP error! Status: ${res.status}`;
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

  const sendMessage = async () => {
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


    const tempMessageId = uuidv4(); // Generate a temporary ID
    // console.log("Temporary ID:", tempMessageId);
    // console.log("sender_id:", typeof userId);
    // console.log("receiver_id:", typeof selectedUser.id);
    // console.log("content:", typeof input);
    // console.log("timestamp:", typeof Date.now());
    // console.log("sender_name:", typeof sender_name);

    const message = {
        ID: tempMessageId, // Include the temporary ID
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

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

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
                    {messages === null ? (
            <div className="text-center text-gray-500">
                Start a conversation with {selectedUser.name}!
            </div>
        ) : (
            messages.map((msg) => (
                <div
                    key={msg.ID} // Use msg.ID as the key
                    className={`my-1 mb-2`}
                >
                    <div
                        className={`p-2 rounded-lg inline-block ${
                            msg.sender_id === parseInt(userId)
                                ? "bg-blue-200 text-right float-right mb-1"
                                : "bg-gray-300 float-left mb-1"
                        }`}
                        style={{ clear: "both" }}
                    >
                        <strong>
                            {msg.sender_id === parseInt(userId)
                                ? "You"
                                : msg.sender_name}
                        </strong>
                        : {msg.content}
                    </div>
                </div>
            ))
        )}

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
              <button
                onClick={sendMessage}
                className="bg-blue-500 text-white p-2 rounded"
              >
                Send
              </button>
            </div>
          </div>
        ) : (
          <p className="text-center text-gray-500">
            Select a user to chat with
          </p>
        )}
      </div>
    </div>
  );
};

export default Chat;

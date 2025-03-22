import React, { useState, useCallback } from "react";
import { getCurrentUserId } from "../utils/getUserId";
import useWebSocket from "@/customComponents/WebsocketConnection";
import { useFetchMessages } from "@/hooks/useFetchMessages";
import { useTypingIndicator } from "@/hooks/useTypingIndicator";
import { MessageDisplay } from "@/customComponents/MessageDisplay";
import useFetchUsers from "@/hooks/useFetchUsers";
import useSendMessage from "@/hooks/useSendMessage";

const Chat = () => {
  const userId = getCurrentUserId();
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState("");
  const [selectedUser, setSelectedUser] = useState(null);

  const { users, loadingUsers } = useFetchUsers(userId);

  const handleMessageReceived = useCallback((message) => {
    setMessages((prevMessages) => {
        return [...prevMessages, message];
    });
  }, []);

  const ws = useWebSocket(userId, handleMessageReceived);
  useFetchMessages(userId, selectedUser, setMessages);

  const handleTyping = useTypingIndicator(setInput);

  const sendMessage = useSendMessage(userId, selectedUser, users, ws, input, setInput);

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

            <MessageDisplay
              messages={messages}
              userId={userId}
              selectedUser={selectedUser}
            />

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

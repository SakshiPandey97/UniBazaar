import React, { useState, useCallback, useEffect } from "react";
import { useLocation, useNavigate } from "react-router-dom"; 
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
  const location = useLocation(); 
  const navigate = useNavigate(); 
  const { users, loadingUsers } = useFetchUsers(userId);

  useEffect(() => {
    const params = new URLSearchParams(location.search);
    const recipientId = params.get("recipient");
    if (recipientId && users.length > 0) {
      const recipientUser = users.find(user => String(user.id) === String(recipientId));
      if (recipientUser) {
        if ((!selectedUser || String(selectedUser.id) !== String(recipientUser.id)) && String(recipientUser.id) !== String(userId)) {
          setSelectedUser(recipientUser);
          console.log("Selected user from URL:", recipientUser);

        } else if (String(recipientUser.id) === String(userId)) {
            console.warn("Attempted to select self from URL parameter.");
        }
      } else {
        console.warn(`Recipient user with ID ${recipientId} not found in user list.`);
      }
    }
  }, [location.search, users, userId, selectedUser, navigate]);

  const handleMessageReceived = useCallback((message) => {
     if (selectedUser &&
         ((String(message.sender_id) === String(userId) && String(message.receiver_id) === String(selectedUser.id)) ||
          (String(message.sender_id) === String(selectedUser.id) && String(message.receiver_id) === String(userId))))
     {
         setMessages((prevMessages) => [...prevMessages, message]);
     } else {
         console.log("Received message for a different conversation:", message);         
     }
  }, [userId, selectedUser]);

  const ws = useWebSocket(userId, handleMessageReceived);
  useFetchMessages(userId, selectedUser, setMessages);

  const handleTyping = useTypingIndicator(setInput);
  const sendMessage = useSendMessage(userId, selectedUser, users, ws, input, setInput);

  const handleSelectUser = (user) => {
     if (String(user.id) === String(userId)) {
         console.log("Cannot select yourself for chat.");
         return;
     }
     setSelectedUser(user);
     navigate('/messaging', { replace: true });
  };

  return (
    <div className="flex h-screen bg-gray-100">
      <div className="w-1/4 bg-white border-r border-gray-300 flex flex-col"> 
        <h3 className="font-bold text-lg p-4 border-b border-gray-300">Contacts</h3> 
        <div className="overflow-y-auto flex-grow"> 
          {loadingUsers ? (
            <div className="p-4 text-center text-gray-500">Loading users...</div>
          ) : (
            <ul>
              {users.map((user) => (
                <li
                  key={user.id}
                  onClick={() => handleSelectUser(user)} 
                  className={`flex items-center p-3 cursor-pointer border-b border-gray-200 ${
                    selectedUser?.id === user.id
                      ? "bg-blue-100 font-semibold" 
                      : "hover:bg-gray-100"
                  } ${String(user.id) === String(userId) ? 'text-gray-400 cursor-not-allowed italic' : ''}`} 
                >
                  <div className="w-8 h-8 rounded-full bg-gray-300 mr-3 flex items-center justify-center text-sm font-semibold">
                      {user.name ? user.name.charAt(0).toUpperCase() : '?'}
                  </div>
                  <span>
                    {user.name} {String(user.id) === String(userId) ? '(You)' : ''}
                  </span>
                </li>
              ))}
            </ul>
          )}
        </div>
      </div>

      {/* Chat Area Panel */}
      <div className="w-3/4 flex flex-col">
        {selectedUser ? (
          <>
            {/* Chat Header */}
            <div className="p-4 border-b border-gray-300 bg-white shadow-sm">
                <h2 className="text-xl font-bold">
                  Chat with {selectedUser.name}
                </h2>
            </div>

            {/* Message Display Area */}
            <div className="flex-grow overflow-y-auto p-4 bg-gray-50">
              <MessageDisplay
                messages={messages}
                userId={userId}
                selectedUser={selectedUser}
              />
            </div>

            <div className="p-4 border-t border-gray-300 bg-white">
              <div className="flex gap-2 items-center">
                <input
                  type="text"
                  value={input}
                  onChange={handleTyping}
                  onKeyPress={(e) => e.key === 'Enter' && !e.shiftKey && sendMessage()}

                  className="flex-1 border p-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  placeholder="Type a message..."
                />
                <button
                  onClick={sendMessage}
                  disabled={!input.trim()} 
                  className="bg-blue-500 text-white px-4 py-2 rounded-lg disabled:opacity-50 hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-opacity-50" // Improved styling
                >
                  Send
                </button>
              </div>
            </div>
          </>
        ) : (
          <div className="flex items-center justify-center h-full text-center text-gray-500">
            <div>
                <svg className="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
                    <path vectorEffect="non-scaling-stroke" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z" />
                </svg>
                <h3 className="mt-2 text-sm font-medium text-gray-900">Select a contact</h3>
                <p className="mt-1 text-sm text-gray-500">Choose someone from the list to start a conversation.</p>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default Chat;

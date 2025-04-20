import React, { useState, useCallback, useEffect } from "react";
import { useLocation }                  from "react-router-dom";
import { getCurrentUserId }             from "@/utils/getUserId";
import useWebSocket                     from "@/customComponents/WebsocketConnection";
import { useFetchMessages }             from "@/hooks/useFetchMessages";
import { useTypingIndicator }           from "@/hooks/useTypingIndicator";
import useSendMessage                   from "@/hooks/useSendMessage";

import ChatHeader   from "@/customComponents/Chat/ChatHeader";
import MessageList  from "@/customComponents/Chat/MessageList";
import MessageInput from "@/customComponents/Chat/MessageInput";

export default function ChatPanel({ users, selectedUser, setSelectedUser }) {
  const userId = getCurrentUserId();
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState("");
  const location = useLocation();

  useEffect(() => {
    const params = new URLSearchParams(location.search);
    const rid    = params.get("recipient");
    if (rid && users && users.length > 0) {
      const u = users.find((u) => String(u.id) === rid);
      if (u && String(u.id) !== String(userId)) {
        if (!selectedUser || String(u.id) !== String(selectedUser.id)) {
            setSelectedUser(u);
        }
      }
    }
  }, [location.search, users, userId, setSelectedUser, selectedUser]);

  const onReceive = useCallback(
    (msg) => {
      if (
        selectedUser &&
        ((String(msg.sender_id)   === String(userId)              &&
          String(msg.receiver_id) === String(selectedUser.id))   ||
         (String(msg.sender_id)   === String(selectedUser.id)    &&
          String(msg.receiver_id) === String(userId)))
      ) {
        const messageWithTimestamp = {
            ...msg,
            created_at: msg.created_at || new Date().toISOString(),
        };
        setMessages((prev) => [...prev, messageWithTimestamp]);
      }
    },
    [userId, selectedUser]
  );
  const ws = useWebSocket(userId, onReceive);
  useFetchMessages(userId, selectedUser, setMessages);
  const handleTyping = useTypingIndicator(setInput);
  const sendMessage = useSendMessage(
    userId,
    selectedUser,
    users,
    ws,
    input,
    setInput,
    setMessages 
  );

  if (!selectedUser) {
    return (
      <div className="flex-grow flex items-center justify-center h-full w-full text-center text-gray-500 bg-white">
        <div>
          <h3 className="mt-2 text-sm font-medium text-gray-900">Select a contact</h3>
          <p className="mt-1 text-sm text-gray-500">Choose someone to start chatting.</p>
        </div>
      </div>
    )
  }

  return (
    <div className="flex flex-col h-full w-full flex-grow bg-white">

      <div className="flex-shrink-0">
        <ChatHeader name={selectedUser.name} status={selectedUser.status} />
      </div>

      <div className="flex-grow min-h-0"> 
        <MessageList
          messages={messages}
          userId={userId}
          selectedUser={selectedUser}
        />
      </div>

      <div className="flex-shrink-0">
        <MessageInput
          input={input}
          onChange={handleTyping}
          onSend={sendMessage}
        />
      </div>
    </div>
  );
}

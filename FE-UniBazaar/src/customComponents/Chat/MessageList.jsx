import React, { useEffect, useRef, useCallback } from "react";
import  {MessageDisplay}  from "@/customComponents/Chat/MessageDisplay";

export default function MessageList({ messages, userId, selectedUser }) {
  const messagesEndRef = useRef(null); 

  const scrollToBottom = useCallback(() => {
    requestAnimationFrame(() => {
        messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
    });
  }, []);

  useEffect(() => {
    if (messages && messages.length > 0) {
        scrollToBottom();
    }
  }, [messages, scrollToBottom]);

  return (
    <div className="overflow-y-auto p-4 bg-gray-50 h-[calc(100vh-14rem)]">
      <MessageDisplay
        messages={messages}
        userId={userId}
        selectedUser={selectedUser}
      />
      <div ref={messagesEndRef} />
    </div>
  );
}

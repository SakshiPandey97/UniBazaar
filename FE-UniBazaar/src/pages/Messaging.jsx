import React, { useState, useEffect, useCallback } from "react";
import { getCurrentUserId } from "@/utils/getUserId";
import useFetchUsers from "../hooks/useFetchUsers";
import useWebSocket from "@/customComponents/WebsocketConnection"; 

import ContactList from "@/customComponents/Chat/ContactList";
import ChatPanel from "@/customComponents/Chat/ChatPanel";

export default function Messaging() {
  const userId = getCurrentUserId();
  const { users, loadingUsers } = useFetchUsers(userId);
  const [selectedUser, setSelectedUser] = useState(null);
  const [unreadSenders, setUnreadSenders] = useState(new Set()); 

  const handleNotificationReceive = useCallback((msg) => {
    if (
      msg &&
      String(msg.receiver_id) === String(userId) &&
      (!selectedUser || String(msg.sender_id) !== String(selectedUser.id))
    ) {
      setUnreadSenders((prev) => new Set(prev).add(String(msg.sender_id)));
    }
  }, [userId, selectedUser]); 

  useWebSocket(userId, handleNotificationReceive);

  // Function to clear notification for a specific user
  const clearNotification = (senderId) => {
    setUnreadSenders((prev) => {
      const newSet = new Set(prev);
      newSet.delete(String(senderId));
      return newSet;
    });
  };

  // When a user is selected, clear their notification
  const handleSelectUser = (user) => {
    setSelectedUser(user);
    clearNotification(user.id); // 
  };

  return (
    <div className="flex h-full w-full overflow-hidden bg-gray-100">
      <ContactList
        users={users}
        loading={loadingUsers}
        currentUserId={userId}
        selectedUserId={selectedUser?.id}
        onSelect={handleSelectUser} 
        unreadSenders={unreadSenders} 
      />

      <div className="w-3/4 flex">
        <ChatPanel
          users={users} 
          selectedUser={selectedUser}
          setSelectedUser={handleSelectUser} 
        />
      </div>
    </div>
  );
}

import React, { useState, useEffect, useCallback, useRef } from "react";
import { getCurrentUserId } from "@/utils/getUserId";
import useFetchUsers from "../hooks/useFetchUsers";
import { getUnreadSendersAPI } from "@/api/messagingAxios"; 

import ContactList from "@/customComponents/Chat/ContactList";
import ChatPanel from "@/customComponents/Chat/ChatPanel";

const POLLING_INTERVAL = 5000;

export default function Messaging() {
  const userId = getCurrentUserId();
  const { users, loadingUsers } = useFetchUsers(userId);
  const [selectedUser, setSelectedUser] = useState(null);
  const [unreadSenders, setUnreadSenders] = useState(new Set());
  const intervalRef = useRef(null);

  const fetchUnreadSenders = useCallback(async () => {
    // console.log("Attempting fetchUnreadSenders. User ID:", userId);

    if (!userId) { 
      // console.log("fetchUnreadSenders returning early, userId is falsy.");
      return;
    }

    try {
      // console.log("Calling getUnreadSendersAPI with userId:", userId);
      const senderIds = await getUnreadSendersAPI(userId); 

      if (Array.isArray(senderIds)) {
        setUnreadSenders((currentUnread) => {
            const newUnread = new Set(senderIds.map(String));
            if (currentUnread.size !== newUnread.size || ![...currentUnread].every(id => newUnread.has(id))) {
                console.log("Polling: Updating unread senders:", newUnread);
                return newUnread;
            }
            return currentUnread;
        });
      } else {
         console.warn("Received unexpected data format for unread senders:", senderIds);
      }
    } catch (error) {
      console.error("Failed to fetch unread senders:", error); 
      if (error.response && (error.response.status === 401 || error.response.status === 403)) {
          if (intervalRef.current) {
              clearInterval(intervalRef.current);
              console.log("Polling stopped due to authorization error.");
          }
      }
    }
  }, [userId]); 

  useEffect(() => {
    fetchUnreadSenders(); 

    if (intervalRef.current) {
      clearInterval(intervalRef.current);
    }

    intervalRef.current = setInterval(fetchUnreadSenders, POLLING_INTERVAL);

    return () => {
      if (intervalRef.current) {
        clearInterval(intervalRef.current);
      }
    };
  }, [fetchUnreadSenders]);


  const handleSelectUser = useCallback((user) => {
    if (user?.id !== selectedUser?.id) {
        setSelectedUser(user);
        if (user) {
          setUnreadSenders((prev) => {
            const newSet = new Set(prev);
            const deleted = newSet.delete(String(user.id));
            return deleted ? newSet : prev;
          });
        }
    }
  }, [selectedUser?.id]);

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

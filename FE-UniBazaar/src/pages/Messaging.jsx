import React, { useState } from "react";
import { getCurrentUserId } from "@/utils/getUserId";
import  useFetchUsers     from "../hooks/useFetchUsers";

import ContactList from "@/customComponents/Chat/ContactList";
import ChatPanel   from "@/customComponents/Chat/ChatPanel";

export default function Messaging() {
  const userId = getCurrentUserId();
  const { users, loadingUsers } = useFetchUsers(userId);
  const [selectedUser, setSelectedUser] = useState(null);

  return (
    <div className="flex h-full w-full overflow-hidden bg-gray-100">
      
      <ContactList
        users={users}
        loading={loadingUsers}
        currentUserId={userId}
        selectedUserId={selectedUser?.id}
        onSelect={setSelectedUser}
      />

      <div className="w-3/4 flex">
        <ChatPanel
          users={users}
          selectedUser={selectedUser}
          setSelectedUser={setSelectedUser}
        />
      </div>
    </div>
  );
}

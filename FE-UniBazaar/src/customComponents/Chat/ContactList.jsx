import React from "react";
import { useNavigate } from "react-router-dom";

export default function ContactList({
  users,
  loading,
  currentUserId,
  selectedUserId,
  onSelect,
}) {
  const navigate = useNavigate();

  const handleClick = (user) => {
    if (String(user.id) === String(currentUserId)) return;
    onSelect(user);
    navigate(`/messaging`, { replace: true });
  };

  return (
    <div className="w-1/4 bg-white border-r border-gray-300 flex flex-col">
      <h3 className="font-bold text-lg p-4 border-b border-gray-300">
        Contacts
      </h3>
      <div className="overflow-y-auto flex-grow">
        {loading ? (
          <div className="p-4 text-center text-gray-500">
            Loading users...
          </div>
        ) : (
          <ul>
            {users.map((user) => {
              const isSelf = String(user.id) === String(currentUserId);
              const isSelected = String(user.id) === String(selectedUserId);
              return (
                <li
                  key={user.id}
                  onClick={() => handleClick(user)}
                  className={`
                    flex items-center p-3 cursor-pointer border-b border-gray-200
                    ${isSelected ? "bg-blue-100 font-semibold" : "hover:bg-gray-100"}
                    ${isSelf ? "text-gray-400 cursor-not-allowed italic" : ""}
                  `}
                >
                  <div className="w-8 h-8 rounded-full bg-gray-300 mr-3 flex items-center justify-center text-sm font-semibold">
                    {user.name?.charAt(0).toUpperCase() || "?"}
                  </div>
                  <span>
                    {user.name} {isSelf && "(You)"}
                  </span>
                </li>
              );
            })}
          </ul>
        )}
      </div>
    </div>
  );
}

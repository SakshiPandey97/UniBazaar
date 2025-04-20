import React from "react";
import { useNavigate } from "react-router-dom";

export default function ContactList({
  users,
  loading,
  currentUserId,
  selectedUserId,
  onSelect,
  unreadSenders, 
}) {
  const navigate = useNavigate();

  const handleClick = (user) => {
    onSelect(user);
    navigate(`/messaging`, { replace: true });
  };

  const filteredUsers = users.filter(user => String(user.id) !== String(currentUserId));

  const getInitials = (nameString) => {
    if (!nameString || typeof nameString !== 'string' || nameString.trim().length === 0) {
      return "?";
    }
    return nameString.charAt(0).toUpperCase();
  };


  return (
    <div className="w-1/4 bg-white border-r border-gray-300 flex flex-col h-full"> 
      <h3 className="font-bold text-lg p-4 border-b border-gray-300 flex-shrink-0">
        Contacts
      </h3>
      <div className="overflow-y-auto flex-grow">
        {loading ? (
          <div className="p-4 text-center text-gray-500">
            Loading users...
          </div>
        ) : filteredUsers.length === 0 ? (
          <div className="p-4 text-center text-gray-500">
            No other contacts found.
          </div>
        ) : (
          <ul>
            {filteredUsers.map((user) => {
              const isSelected = String(user.id) === String(selectedUserId);
              const hasUnread = unreadSenders.has(String(user.id));

              return (
                <li
                  key={user.id}
                  onClick={() => handleClick(user)}
                  className={`
                    flex items-center p-3 cursor-pointer border-b border-gray-200 transition-colors duration-150 ease-in-out relative
                    ${isSelected ? "bg-blue-100 font-semibold text-blue-800" : "hover:bg-gray-100"} {/* Added text color for selected */}
                  `}
                  aria-current={isSelected ? 'page' : undefined} 
                >
                  <div className={`
                    flex-shrink-0 w-8 h-8 rounded-full mr-3 flex items-center justify-center text-sm font-semibold
                    ${isSelected ? 'bg-blue-500 text-white' : 'bg-gray-300 text-gray-700'} {/* Different style for selected */}
                  `}>
                    {getInitials(user.name)}
                  </div>
                  {/* User Name */}
                  <span className="truncate flex-grow mr-2"> {/* Added mr-2 for spacing before dot */}
                    {user.name || "Unknown User"} {/* Fallback for missing name */}
                  </span>
                  {/* Unread Dot */}
                  {hasUnread && (
                    <span
                      className="ml-auto flex-shrink-0 w-2.5 h-2.5 bg-red-500 rounded-full animate-pulse"
                      aria-label="Unread messages" // Accessibility label
                    ></span>
                  )}
                </li>
              );
            })}
          </ul>
        )}
      </div>
    </div>
  );
}

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
    onSelect(user);
    navigate(`/messaging`, { replace: true });
  };

  const filteredUsers = users.filter(user => String(user.id) !== String(currentUserId));

  return (
    <div className="w-1/4 bg-white border-r border-gray-300 flex flex-col">
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
              return (
                <li
                  key={user.id}
                  onClick={() => handleClick(user)}
                  className={`
                    flex items-center p-3 cursor-pointer border-b border-gray-200 transition-colors duration-150 ease-in-out {/* Added transition */}
                    ${isSelected ? "bg-blue-100 font-semibold" : "hover:bg-gray-100"}
                    {/* Removed isSelf styling logic */}
                  `}
                >
                  <div className="flex-shrink-0 w-8 h-8 rounded-full bg-gray-300 mr-3 flex items-center justify-center text-sm font-semibold"> {/* Added flex-shrink-0 */}
                    {user.name?.charAt(0).toUpperCase() || "?"}
                  </div>
                  <span className="truncate">
                    {user.name}
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

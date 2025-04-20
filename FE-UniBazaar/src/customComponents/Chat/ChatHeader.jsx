import React from "react";

export default function ChatHeader({ name, status }) {
  const getInitials = (nameString) => {
    if (!nameString || typeof nameString !== 'string' || nameString.trim().length === 0) {
      return "?";
    }

    const names = nameString.trim().split(' ').filter(n => n); 

    if (names.length === 0) {
      return "?";
    } else if (names.length === 1) {
      return names[0].charAt(0).toUpperCase();
    } else {
      const firstNameInitial = names[0].charAt(0);
      const lastNameInitial = names[names.length - 1].charAt(0); 
      return `${firstNameInitial}${lastNameInitial}`.toUpperCase();
    }
  };

  return (
    <div className="p-4 border-b border-gray-300 bg-white shadow-sm flex justify-center items-center space-x-3">
      <div className="flex-shrink-0 w-8 h-8 rounded-full bg-gray-300 flex items-center justify-center text-sm font-semibold">
        {getInitials(name)}
      </div>

      <h2 className="font-semibold text-lg">
        {name || "Chat"}
      </h2>

    </div>
  );
}

import React from "react";

export default function MessageInput({ input, onChange, onSend }) {
  const handleSend = (e) => {
    e?.preventDefault();
    if (input.trim()) { 
        onSend();
    }
  };

  const handleKeyPress = (e) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault(); 
      handleSend(); 
    }
  };

  return (
    <div className="p-4 border-t border-gray-300 bg-white">
      <form onSubmit={handleSend} className="flex gap-2 items-center">
        <input
          type="text"
          value={input}
          onChange={onChange}
          onKeyUp ={handleKeyPress} 
          className="flex-1 border p-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          placeholder="Type a message..."
          autoComplete="off"
        />
        <button
          type="submit" 
          disabled={!input.trim()}
          className="bg-blue-500 text-white px-4 py-2 rounded-lg disabled:opacity-50 hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-opacity-50"
        >
          Send
        </button>
      </form>
    </div>
  );
}

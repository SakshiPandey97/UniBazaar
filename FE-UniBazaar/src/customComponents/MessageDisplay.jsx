import React, { useEffect, useRef } from 'react';

export const MessageDisplay = ({ messages, userId, selectedUser }) => { // Add export here
    const messagesEndRef = useRef(null);

    useEffect(() => {
        if (messagesEndRef.current) {
            messagesEndRef.current.scrollIntoView({ behavior: "smooth" });
        }
    }, [messages]);

    return (
        <div className="h-96 overflow-y-auto border p-2 mb-3 bg-gray-100 rounded">
            {messages === null ? (
                <div className="text-center text-gray-500">
                    Start a conversation with {selectedUser.name}!
                </div>
            ) : (
                messages.map((msg) => (
                    <div key={msg.ID} className={`my-1 mb-2`}>
                        <div
                            className={`p-2 rounded-lg inline-block ${msg.sender_id === parseInt(userId) ? "bg-blue-200 text-right float-right mb-1" : "bg-gray-300 float-left mb-1"}`}
                            style={{ clear: "both" }}
                        >
                            <strong>
                                {msg.sender_id === parseInt(userId) ? "You" : msg.sender_name}
                            </strong>
                            : {msg.content}
                        </div>
                    </div>
                ))
            )}
            <div ref={messagesEndRef} />
        </div>
    );
};

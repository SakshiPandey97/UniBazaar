import React, { useEffect, useRef, useCallback } from 'react';

const MessageItem = React.memo(({ msg, userId }) => {
    return (
        <div className={`my-1 mb-2 ${msg.sender_id === parseInt(userId) ? "text-right" : "text-left"}`}>
            <div
                className={`p-2 rounded-lg inline-block ${msg.sender_id === parseInt(userId) ? "bg-blue-200" : "bg-gray-300"}`}
            >
                <strong>
                    {msg.sender_id === parseInt(userId) ? "You" : msg.sender_name}
                </strong>
                : {msg.content}
            </div>
            <div style={{ clear: "both" }}></div>
        </div>
    );
});

export const MessageDisplay = ({ messages, userId, selectedUser }) => {
    const messagesEndRef = useRef(null);
    const messagesContainerRef = useRef(null);
    const isInitialRender = useRef(true);
    
    const scrollToBottom = useCallback(() => {
        setTimeout(() => {
            if (messagesEndRef.current) {
                messagesEndRef.current.scrollIntoView({ behavior: 'smooth' });
            }
        }, 10);
    }, []);
    
    useEffect(() => {
        if (!messages) return;
        
        if (isInitialRender.current) {
            isInitialRender.current = false;
            scrollToBottom();
            return;
        }
        
        scrollToBottom();
    }, [messages, scrollToBottom]);
    
    return (
        <div
            className="h-96 overflow-y-auto border p-2 mb-3 bg-gray-100 rounded"
            ref={messagesContainerRef}
        >
            {messages === null ? (
                <div className="text-center text-gray-500">
                    Start a conversation with {selectedUser.name}!
                </div>
            ) : (
                messages.map((msg) => (
                    <MessageItem key={msg.id} msg={msg} userId={userId} />
                ))
            )}
            <div ref={messagesEndRef} style={{ clear: "both", float: "none" }}></div>
        </div>
    );
};
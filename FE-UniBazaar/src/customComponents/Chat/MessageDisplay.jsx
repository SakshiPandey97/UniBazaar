import React from 'react';

const MessageItem = React.memo(({ msg, userId }) => {
    const isSender = String(msg.sender_id) === String(userId);

    return (
        <div className={`my-1 mb-2 flex ${isSender ? 'justify-end' : 'justify-start'}`}>
            <div
                className={`p-2 rounded-lg max-w-xs sm:max-w-md md:max-w-lg ${
                    isSender ? "bg-blue-200 text-blue-900" : "bg-gray-300 text-gray-800"
                }`}
            >
                {/* {!isSender && msg.sender_name && (
                    <p className="text-xs font-semibold mb-1">{msg.sender_name}</p>
                )} */}
                <p className="text-sm break-words">{msg.content}</p>
            </div>
        </div>
    );
});

export const MessageDisplay = ({ messages, userId, selectedUser }) => {
    const noMessages = !messages || messages.length === 0;

    return (
        <div>
            {noMessages ? (
                <div className="text-center text-gray-500 py-10">
                    {selectedUser?.name
                        ? `Start a conversation with ${selectedUser.name}!`
                        : "No messages yet."}
                </div>
            ) : (
                messages.map((msg) => (
                    <MessageItem key={msg.id || msg.temp_id} msg={msg} userId={userId} />
                ))
            )}
        </div>
    );
};

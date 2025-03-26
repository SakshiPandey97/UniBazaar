import { useEffect } from 'react'; 

export const useFetchMessages = (userId, selectedUser, setMessages) => { 
    useEffect(() => {
        if (!selectedUser || !userId) return; // Ensure both IDs are available before making the request

        const fetchMessages = async () => {
            try {
                const url = `http://localhost:8080/api/conversation/${userId}/${selectedUser.id}`;
                console.log("Fetching messages from:", url);
                const res = await fetch(url);

                if (!res.ok) {
                    const errorData = await res.json().catch(() => null);
                    const errorMessage = errorData?.message || `HTTP error! Status: ${res.status}`;
                    throw new Error(errorMessage);
                }

                const data = await res.json();

                // Ensure the response is an array, even if it's empty
                if (Array.isArray(data)) {
                    setMessages(data.length > 0 ? data : []); // Set an empty array if no messages exist
                } else {
                    setMessages([]); // Default to empty array if unexpected response
                }
            } catch (error) {
                console.error("Failed to load messages:", error);
                setMessages([]); // Ensure the UI doesn't break if fetching fails
            }
        };

        fetchMessages();
    }, [selectedUser, userId, setMessages]);
};

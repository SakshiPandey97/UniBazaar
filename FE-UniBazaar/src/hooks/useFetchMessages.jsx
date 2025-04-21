import { useEffect } from 'react'; 
const CHAT_USERS_BASE_URL = import.meta.env.VITE_CHAT_USERS_BASE_URL;


export const useFetchMessages = (userId, selectedUser, setMessages) => { 
    useEffect(() => {
        if (!selectedUser || !userId) return; 

        const fetchMessages = async () => {
            try {
                const url = `${CHAT_USERS_BASE_URL}/api/conversation/${userId}/${selectedUser.id}`;
                console.log("Fetching messages from:", url);
                const res = await fetch(url);

                if (!res.ok) {
                    const errorData = await res.json().catch(() => null);
                    const errorMessage = errorData?.message || `HTTP error! Status: ${res.status}`;
                    throw new Error(errorMessage);
                }

                const data = await res.json();

                if (Array.isArray(data)) {
                    setMessages(data.length > 0 ? data : []); 
                } else {
                    setMessages([]); 
                }
            } catch (error) {
                console.error("Failed to load messages:", error);
                setMessages([]); 
            }
        };

        fetchMessages();
    }, [selectedUser, userId, setMessages]);
};

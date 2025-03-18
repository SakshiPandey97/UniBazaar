import { useEffect } from 'react'; 

export const useFetchMessages = (userId, selectedUser, setMessages) => { 
    useEffect(() => {
        if (!selectedUser || !userId) return;

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
                setMessages(data);
            } catch (error) {
                console.error("Failed to load messages:", error);
            }
        };

        fetchMessages();
    }, [selectedUser, userId, setMessages]);
};

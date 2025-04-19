// Rename the file to useStartChat.jsx (recommended)
import { useNavigate } from 'react-router-dom';
import { getCurrentUserId } from '@/utils/getUserId'; // Assuming this doesn't use hooks internally

// Define the custom hook
export const useStartChat = () => {
  const navigate = useNavigate(); // Call hook inside the custom hook
  const currentUserId = getCurrentUserId(); // Get the user ID

  // Define the handler function *inside* the hook
  const handleStartChat = (sellerId) => {
    if (!currentUserId) {
      console.error("Current user ID not found. Cannot start chat.");
      // Optionally redirect to login or show a message
      // navigate('/login'); // You can navigate here if needed
      return;
    }

    // Prevent user from messaging themselves
    // Compare as strings for reliability
    if (String(currentUserId) === String(sellerId)) {
      console.log("You cannot message yourself.");
      // Optionally show a user-friendly alert/toast message here
      alert("You cannot message yourself."); // Simple alert example
      return;
    }

    console.log(`Navigating to chat with seller: ${sellerId}`);
    // Navigate to the messaging page with the sellerId as a query parameter
    navigate(`/messaging?recipient=${sellerId}`);
  };

  // Return the handler function from the hook
  return handleStartChat;
};

// Default export is often preferred for hooks, but named export is fine too
export default useStartChat;

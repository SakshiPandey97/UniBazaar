import { useNavigate } from 'react-router-dom';
import { getCurrentUserId } from '@/utils/getUserId'; 

export const useStartChat = () => {
  const navigate = useNavigate(); 
  const currentUserId = getCurrentUserId(); 

  const handleStartChat = (sellerId) => {
    if (!currentUserId) {
      console.error("Current user ID not found. Cannot start chat.");
      return;
    }

    if (String(currentUserId) === String(sellerId)) {
      console.log("You cannot message yourself.");
      alert("You cannot message yourself."); 
      return;
    }

    console.log(`Navigating to chat with seller: ${sellerId}`);
    navigate(`/messaging?recipient=${sellerId}`);
  };

  return handleStartChat;
};

export default useStartChat;

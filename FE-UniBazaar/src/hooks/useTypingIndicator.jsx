import { useRef } from 'react'; 

export const useTypingIndicator = (setInput, setTyping) => { 
    const typingTimeoutRef = useRef(null);
  
    const handleTyping = (e) => {
      setInput(e.target.value);
      if (setTyping) setTyping(true);
  
      if (typingTimeoutRef.current) clearTimeout(typingTimeoutRef.current);
      typingTimeoutRef.current = setTimeout(() => setTyping(false), 2000);
    };
  
    return handleTyping;
  };

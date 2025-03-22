export const useTypingIndicator = (setInput) => {
  const handleTyping = (e) => {
    setInput(e.target.value);
  };
  
  return handleTyping;
};

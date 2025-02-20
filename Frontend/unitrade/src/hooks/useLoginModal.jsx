import { useEffect, useState } from "react";
import { useUserAuth } from "./useUserAuth";

function useAuthModal() {
  const [isModalOpen, setModalOpen] = useState(false);
  const userAuth = useUserAuth();

  const toggleModal = () => {
    setModalOpen((prev) => !prev);
  };

  useEffect(() => {
    document.body.style.overflow = isModalOpen ? "hidden" : "auto";
    return () => {
      document.body.style.overflow = "auto";
    };
  }, [isModalOpen]);

  return { isModalOpen, toggleModal };
}

export default useAuthModal;

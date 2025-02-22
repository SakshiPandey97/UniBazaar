import { useEffect, useState } from "react";
import { useUserAuth } from "./useUserAuth";

function useLoginModal() {
  const [isModalOpen, setModalOpen] = useState(false);
  const [isViewProfileOpen, setViewProfileOpen] = useState(false);
  const userAuth = useUserAuth();

  const toggleModal = () => {
    setModalOpen((prev) => !prev);
  };

  const openViewMyProfile = () => {
    setViewProfileOpen((prev) => !prev);
  };

  useEffect(() => {
    document.body.style.overflow = isModalOpen || isViewProfileOpen ? "hidden" : "auto";
    return () => {
      document.body.style.overflow = "auto";
    };
  }, [isModalOpen]);

  return { isModalOpen,isViewProfileOpen, toggleModal ,openViewMyProfile};
}

export default useLoginModal;

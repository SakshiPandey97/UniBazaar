import React from "react";
import { motion, AnimatePresence } from "framer-motion";

const ModalVariants = {
  hidden: { opacity: 0, scale: 0.8 },
  visible: {
    opacity: 1,
    scale: 1,
    transition: { type: "spring", stiffness: 120, damping: 10 },
  },
  exit: {
    opacity: 0,
    scale: 0.8,
    transition: { duration: 0.3, ease: "easeInOut" },
  },
};

const Modal = ({ isOpen, toggleModal, children }) => (
  <AnimatePresence>
    {isOpen && (
      <div className="modal-overlay" onClick={toggleModal}>
        <motion.div
          variants={ModalVariants}
          initial="hidden"
          animate="visible"
          exit="exit"
          className="modal"
          onClick={(e) => e.stopPropagation()}
        >
          {children}
        </motion.div>
      </div>
    )}
  </AnimatePresence>
);

export default Modal;

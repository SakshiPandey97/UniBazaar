import React from "react";
import { motion } from "framer-motion";

const ModalVariants = {
  hidden: { x: "-100vw" },
  visible: { opacity: 1, x: 0, transition: { transition: { type: "spring" } } },
};

const Modal = ({ isOpen, toggleModal, children }) =>
  isOpen && (
    <div className="modal-overlay" onClick={toggleModal}>
      <motion.div
        variants={ModalVariants}
        initial="hidden"
        animate="visible"
        className="modal"
        onClick={(e) => e.stopPropagation()}
      >
        {children}
      </motion.div>
    </div>
  );
export default Modal;

import React from 'react';
import { motion } from 'framer-motion';

function PageWrapper({ children }) {
  const transitionVariants = {
    initial: {
      x: "-100vw",
    },
    animate: {
      x: 0,
    },
    exit: {
      x: "100vw",
    }
  };

  return (
    <motion.div
      initial="initial"
      animate="animate"
      exit="exit"
      variants={transitionVariants}
    >
      {children}
    </motion.div>
  );
}

export default PageWrapper;

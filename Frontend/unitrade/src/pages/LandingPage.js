import React, { useState, useCallback, useRef, useEffect } from 'react';
import { useTransition, animated } from '@react-spring/web';
// import logo from './UniBazaar_logo.png';
import '../App.css';

const Banner = () => {
  const ref = useRef([]);
  const [items, set] = useState([]);

  const transitions = useTransition(items, {
    from: { opacity: 0, height: 0, innerHeight: 0, transform: 'perspective(600px) rotateX(0deg)', color: '#8fa5b6' },
    enter: [
      { opacity: 1, height: 80, innerHeight: 80 },
      { transform: 'perspective(600px) rotateX(90deg)', color: '#28d79f' },
      { transform: 'perspective(600px) rotateX(0deg)' },
    ],
    leave: [{ color: '#c23369' }, { innerHeight: 0 }, { opacity: 0, height: 0 }],
    update: { color: '#28b4d7' },
  });

  const reset = useCallback(() => {
    ref.current.forEach(clearTimeout);
    ref.current = [];
    const sequence = [['UNI', 'Bazaar'], ['UNI'], ['UNI', 'Bazaar']];
    sequence.forEach((item, index) => {
      ref.current.push(setTimeout(() => set(item), index * 3000));
    });
  }, []);

  useEffect(() => {
    reset();
    return () => ref.current.forEach(clearTimeout);
  }, [reset]);

  return (
    <div className="App-banner-container">
      <div className="App-background"></div>
      {/* <img src={logo} className="App-logo" alt="Logo" /> */}
      <div className="App-animated-text">
        {transitions(({ innerHeight, ...rest }, item) => (
          <animated.div className="App-text-item" style={rest} onClick={reset}>
            <animated.div style={{ overflow: 'hidden', height: innerHeight }}>{item}</animated.div>
          </animated.div>
        ))}
      </div>
    </div>
  );
};

export default Banner;
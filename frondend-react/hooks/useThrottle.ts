import { useState, useCallback } from "react";

export const useThrottle = (callback: () => void, delay: number) => {
  const [isThrottled, setIsThrottled] = useState(false);

  const throttledCallback = useCallback(() => {
    if (!isThrottled) {
      callback();
      setIsThrottled(true);
      setTimeout(() => setIsThrottled(false), delay);
    }
  }, [callback, delay, isThrottled]);

  return throttledCallback;
};

import React, { useCallback } from "react";

export const useThrottle = (callback: () => void, ms: number = 1000) => {
    const isCalled = React.useRef<boolean>( false );


    const throttledFunction = useCallback( () => {
        if (isCalled.current) {
            console.log( `Функция уже вызвана, подождите ${ ms } миллисекунд` );
            return;
        }

        isCalled.current = true;

        callback();

        setTimeout( () => {
            isCalled.current = false;
        }, ms );
    }, [ callback, ms ] );

    return {
        throttledFunction,
    };
};

import React, { forwardRef, HTMLAttributes } from 'react';

interface IButton extends HTMLAttributes<HTMLButtonElement> {
    className?: string;
}

const ButtonComponent = forwardRef<HTMLButtonElement, IButton>( (props, ref) => {
    const {
        className,
        children,
        ...otherProps
    } = props;

    return (
        <button ref={ ref } className={ className } { ...otherProps }>
            { children }
        </button>
    );
} );

export const Button = React.memo( ButtonComponent );

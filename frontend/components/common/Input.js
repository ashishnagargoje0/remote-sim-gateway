import { forwardRef } from 'react';
import clsx from 'clsx';

const Input = forwardRef(({
  className,
  type = 'text',
  label,
  error,
  helper,
  required,
  ...props
}, ref) => {
  const inputId = props.id || props.name;

  return (
    <div className="w-full">
      {label && (
        <label 
          htmlFor={inputId}
          className="block text-sm font-medium text-gray-700 mb-1"
        >
          {label}
          {required && <span className="text-red-500 ml-1">*</span>}
        </label>
      )}
      
      <input
        type={type}
        id={inputId}
        ref={ref}
        className={clsx(
          'block w-full rounded-md border-gray-300 shadow-sm transition-colors',
          'focus:border-primary-500 focus:ring-primary-500',
          'placeholder:text-gray-400',
          'disabled:bg-gray-50 disabled:text-gray-500 disabled:cursor-not-allowed',
          error && 'border-red-300 focus:border-red-500 focus:ring-red-500',
          className
        )}
        {...props}
      />
      
      {error && (
        <p className="mt-1 text-sm text-red-600">{error}</p>
      )}
      
      {helper && !error && (
        <p className="mt-1 text-sm text-gray-500">{helper}</p>
      )}
    </div>
  );
});

Input.displayName = 'Input';

export default Input;
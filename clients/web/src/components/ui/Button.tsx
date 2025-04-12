"use client";

import React from "react";

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: "primary" | "secondary";
  fullWidth?: boolean;
}

export const Button: React.FC<ButtonProps> = ({
  children,
  variant = "primary",
  fullWidth = false,
  className = "",
  ...props
}) => {
  const baseStyles = "font-bold py-3 px-5 rounded-lg text-base";
  
  const variantStyles = {
    primary: "bg-[#FBD5BD] text-black",
    secondary: "bg-white border border-gray-300 text-black",
  };
  
  const widthStyles = fullWidth ? "w-full" : "";
  
  return (
    <button
      className={`${baseStyles} ${variantStyles[variant]} ${widthStyles} ${className}`}
      {...props}
    >
      {children}
    </button>
  );
};
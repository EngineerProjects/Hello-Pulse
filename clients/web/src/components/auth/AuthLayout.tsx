import React from "react";

interface AuthLayoutProps {
  children: React.ReactNode;
  title: string;
  description: string;
}

export const AuthLayout: React.FC<AuthLayoutProps> = ({
  children,
  title,
  description,
}) => {
  return (
    <div className="bg-white h-screen overflow-hidden relative">
      <div className="grid grid-cols-1 lg:grid-cols-2 h-screen overflow-hidden">
        <div
          className="flex-col justify-center items-center relative h-full bg-svg-pattern bg-cover bg-center bg-no-repeat hidden lg:block"
          style={{ backgroundImage: "url('/images/login_background.svg')", width: "112%" }}
        >
          <div className="absolute inset-0 flex justify-center items-center w-full lg:w-[85%] px-8">
            <div className="items-center pl-14">
              <h1 className="text-4xl lg:text-7xl xl:text-8xl font-bold text-[#312783] w-[70%] font-impact">
                {title}
              </h1>
              <p className="mt-4 font-outfit text-sm w-[50%] text-black">
                {description}
              </p>
            </div>
          </div>
        </div>
        <div className="ml-10 flex flex-col md:w-[90%] lg:w-[70%] w-full items-center justify-center flex-1 h-screen px-4 lg:px-8">
          {children}
        </div>
      </div>
    </div>
  );
};
"use client";

import React from "react";
import Image from "next/image";

export const SocialLoginButtons: React.FC = () => {
  return (
    <div className="flex items-center justify-center gap-4 my-4">
      <button className="bg-white border border-gray-300 p-2 rounded">
        <Image src="/images/Google.png" alt="Google" width={40} height={40} />
      </button>
      <button className="bg-white border border-gray-300 p-2 rounded">
        <Image src="/images/Microsoft.png" alt="Microsoft" width={40} height={40} />
      </button>
      <button className="bg-white border border-gray-300 p-2 rounded">
        <Image src="/images/Apple.png" alt="Apple" width={40} height={40} />
      </button>
    </div>
  );
};
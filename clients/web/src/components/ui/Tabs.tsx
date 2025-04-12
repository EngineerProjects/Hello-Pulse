"use client";

import React from "react";
import Link from "next/link";

interface Tab {
  label: string;
  href: string;
  isActive: boolean;
}

interface TabsProps {
  tabs: Tab[];
}

export const Tabs: React.FC<TabsProps> = ({ tabs }) => {
  return (
    <div className="flex mb-3 shadow-md rounded-full font-bold border-[#FBD5BD] border-2 w-full text-black">
      {tabs.map((tab, index) => (
        <Link
          key={index}
          href={tab.href}
          className={`px-4 lg:px-8 py-4 rounded-full w-1/2 text-center ${
            tab.isActive ? "bg-[#FBD5BD]" : "bg-white"
          }`}
        >
          {tab.label}
        </Link>
      ))}
    </div>
  );
};
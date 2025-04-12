import React from "react";
import Link from "next/link";

export const AuthFooter: React.FC = () => {
  return (
    <div className="text-[#8CB2DF] bottom-0 text-xs lg:text-sm xl:text-lg text-center py-4 w-full pt-10">
      <Link href="/" className="mx-4">
        Conditions d&apos;utilisation
      </Link>{" "}
      |{" "}
      <Link href="/" className="mx-4">
        Politique de confidentialité
      </Link>
    </div>
  );
};
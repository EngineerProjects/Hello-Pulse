import React from "react";
import Image from "next/image";
import Link from "next/link";
import { RegisterForm } from "@/components/auth/RegisterForm";
import { SocialLoginButtons } from "@/components/auth/SocialLoginButtons";
import { AuthFooter } from "@/components/auth/AuthFooter";
import { AuthLayout } from "@/components/auth/AuthLayout";

export default function RegisterPage() {
  return (
    <AuthLayout
      title="À vos marque, prêt, feu, go ! Brainstormer"
      description="Entrez dans votre espace créatif et donnez vie à vos idées avec Hello Pulse. Connectez-vous pour rejoindre la prochaine session de brainstorming en direct, où l'inspiration est sans limites et les idées fusent en un éclair !"
    >
      <div className="lg:px-8 md:px-4 px-2 flex flex-col items-center md:w-[80%] lg:w-full gap-2">
        <Image
          src="/images/logo.png"
          alt="Hello Pulse"
          width={192}
          height={144}
          className="w-24 h-16 lg:w-32 lg:h-24 xl:w-48 xl:h-36 mb-6 md:mb-8 lg:mb-12"
        />
        
        <div className="flex mb-3 shadow-md rounded-full font-bold border-[#FBD5BD] border-2 w-full text-black">
          <Link
            href="/auth/login"
            className="px-4 lg:px-8 py-4 rounded-full w-1/2 text-center bg-white"
          >
            Se connecter
          </Link>
          <Link
            href="/auth/register"
            className="py-4 px-4 rounded-full text-center w-1/2 bg-[#FBD5BD]"
          >
            Créer un compte
          </Link>
        </div>
        
        <RegisterForm />
        
        <div className="flex w-full justify-center items-center gap-2">
          <div className="border-t-2 border-black w-[35%]"></div>
          <span className="text-black">OU</span>
          <div className="border-t-2 border-black w-[35%]"></div>
        </div>
        
        <SocialLoginButtons />
        
        <AuthFooter />
      </div>
    </AuthLayout>
  );
}
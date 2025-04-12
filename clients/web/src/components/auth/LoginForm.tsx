"use client";

import React, { useState } from "react";
import { Input } from "@/components/ui/Input";
import { Button } from "@/components/ui/Button";
import { useRouter } from "next/navigation";
import Link from "next/link";

export const LoginForm: React.FC = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const router = useRouter();

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    try {
      // Login API call would go here
      console.log("Login attempt", { email, password });
      // On success:
      // router.push("/dashboard");
    } catch (err) {
      setError("Une erreur s'est produite lors de la connexion.");
    }
  };

  return (
    <form className="flex w-[80%] flex-col items-center" onSubmit={handleLogin}>
      <Input
        type="email"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        placeholder="Identifiant"
        required
      />
      
      <Input
        type="password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        placeholder="Mot de passe"
        required
      />
      
      <Button type="submit">Connexion</Button>
      
      <div className="mt-4">
        <Link href="/auth/forgot-password" className="text-[#8CB2DF] text-sm underline">
          Mot de passe oublié ?
        </Link>
      </div>
    </form>
  );
};
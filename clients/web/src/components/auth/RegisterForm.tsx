"use client";

import React, { useState } from "react";
import { Input } from "@/components/ui/Input";
import { Button } from "@/components/ui/Button";
import Link from "next/link";

export const RegisterForm: React.FC = () => {
  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [error, setError] = useState<string | null>(null);

  const handleRegister = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    if (password !== confirmPassword) {
      setError("Les mots de passe ne correspondent pas.");
      return;
    }

    try {
      // Registration API call would go here
      console.log("Register attempt", { firstName, lastName, email, password });
      // On success:
      // router.push("/auth/login");
    } catch (err) {
      setError("Une erreur s'est produite lors de l'inscription.");
    }
  };

  return (
    <form className="flex w-[80%] flex-col items-center" onSubmit={handleRegister}>
      <Input
        type="text"
        value={lastName}
        onChange={(e) => setLastName(e.target.value)}
        placeholder="Nom"
        required
      />
      
      <Input
        type="text"
        value={firstName}
        onChange={(e) => setFirstName(e.target.value)}
        placeholder="Prénom"
        required
      />
      
      <Input
        type="email"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        placeholder="Mail"
        required
      />
      
      <Input
        type="password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        placeholder="Mot de passe"
        required
      />
      
      <Input
        type="password"
        value={confirmPassword}
        onChange={(e) => setConfirmPassword(e.target.value)}
        placeholder="Confirmer le mot de passe"
        required
      />
      
      <Button type="submit">S'enregistrer</Button>
      
      <div className="mt-2">
        <span className="text-sm text-center text-black">
          Vous avez déjà un compte ?{" "}
        </span>
        <Link href="/auth/login" className="text-[#8CB2DF] text-sm underline">
          Se connecter
        </Link>
      </div>
    </form>
  );
};
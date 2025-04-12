"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { LoginCredentials, RegisterCredentials, User } from "@/types/auth";

export function useAuth() {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const router = useRouter();

  useEffect(() => {
    // Check if user is already logged in
    const token = localStorage.getItem("token");
    if (token) {
      // Validate token and get user info
      setUser({
        id: "1",
        email: "user@example.com",
        firstName: "John",
        lastName: "Doe",
      });
    }
    setLoading(false);
  }, []);

  const login = async (credentials: LoginCredentials) => {
    setLoading(true);
    setError(null);
    try {
      // Call the login API
      const user = {
        id: "1",
        email: credentials.email,
        firstName: "John",
        lastName: "Doe",
      };
      const token = "fake-token";
      
      localStorage.setItem("token", token);
      setUser(user);
      router.push("/dashboard");
      return { user, token };
    } catch (err) {
      setError("Failed to login");
      throw err;
    } finally {
      setLoading(false);
    }
  };

  const register = async (credentials: RegisterCredentials) => {
    setLoading(true);
    setError(null);
    try {
      // Call the register API
      const user = {
        id: "1",
        email: credentials.email,
        firstName: credentials.firstName,
        lastName: credentials.lastName,
      };
      const token = "fake-token";
      
      localStorage.setItem("token", token);
      setUser(user);
      router.push("/dashboard");
      return { user, token };
    } catch (err) {
      setError("Failed to register");
      throw err;
    } finally {
      setLoading(false);
    }
  };

  const logout = () => {
    localStorage.removeItem("token");
    setUser(null);
    router.push("/auth/login");
  };

  return {
    user,
    loading,
    error,
    login,
    register,
    logout,
  };
}
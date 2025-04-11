'use client';

import React, { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import Image from 'next/image';
import { useAuth } from '@/lib/auth/use-auth';

interface LoginFormProps {
  initialSignUp?: boolean;
}

export default function LoginForm({ initialSignUp = false }: LoginFormProps) {
  const router = useRouter();
  const { login, register, error: authError, clearError } = useAuth();
  
  // Form state
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const [alertMessage, setAlertMessage] = useState<string | null>(null);
  const [signUp, setSignUp] = useState(initialSignUp);
  const [firstName, setFirstName] = useState('');
  const [lastName, setLastName] = useState('');
  const [confirmedPassword, setConfirmedPassword] = useState('');
  
  // Input styles
  const inputStyles = {
    WebkitBoxShadow: "0 0 0 30px #181c2e inset",
    WebkitTextFillColor: "white"
  };
  
  // Reset form fields
  const resetInfo = () => {
    setAlertMessage(null);
    setFirstName('');
    setLastName('');
    setEmail('');
    setPassword('');
    setConfirmedPassword('');
  };
  
  // Handle login submission
  const handleLogin = async (e: React.MouseEvent<HTMLButtonElement>) => {
    e.preventDefault();
    setAlertMessage(null);
    
    if (email && password) {
      try {
        setLoading(true);
        await login(email, password);
        // Router will handle redirect based on user state from auth context
      } catch (err: any) {
        setAlertMessage(err.message || 'Login email or password incorrect');
      } finally {
        setLoading(false);
      }
    } else {
      setAlertMessage('Please fill all fields');
    }
  };
  
  // Handle signup submission
  const handleSignUp = async (e: React.MouseEvent<HTMLButtonElement>) => {
    e.preventDefault();
    setAlertMessage(null);
    
    if (password !== confirmedPassword) {
      setAlertMessage('Passwords do not match');
      return;
    }
    
    if (firstName && lastName && email && password) {
      try {
        setLoading(true);
        await register({
          email,
          first_name: firstName,
          last_name: lastName,
          phone: "0", // Default values as in original
          address: "0", // Default values as in original
          password
        });
        // Router will handle redirect from auth context
      } catch (err: any) {
        setAlertMessage(err.message || 'Registration failed');
      } finally {
        setLoading(false);
      }
    } else {
      setAlertMessage('Please fill all fields');
    }
  };
  
  // Toggle between login and register forms
  const clickLogin = () => {
    setSignUp(false);
    resetInfo();
  };
  
  const clickRegister = () => {
    setSignUp(true);
    resetInfo();
  };
  
  // Handle Enter key press
  const handleKeydown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      if (signUp) {
        handleSignUp(e as unknown as React.MouseEvent<HTMLButtonElement>);
      } else {
        handleLogin(e as unknown as React.MouseEvent<HTMLButtonElement>);
      }
    }
  };
  
  useEffect(() => {
    window.addEventListener('keydown', handleKeydown as unknown as EventListener);
    
    return () => {
      window.removeEventListener('keydown', handleKeydown as unknown as EventListener);
    };
  }, [signUp, email, password, firstName, lastName, confirmedPassword]);
  
  return (
    <div className="lg:px-8 md:px-4 px-2 flex flex-col items-center md:w-[80%] lg:w-full gap-2">
      <Image
        src="/images/logo.png"
        alt="Hello Pulse"
        width={144}
        height={108}
        className="w-24 h-16 lg:w-32 lg:h-24 xl:w-48 xl:h-36 mb-3"
      />
      
      <div className="flex mb-3 shadow-md rounded-full font-bold border-[#FBD5BD] border-2 w-full text-black dark:text-white">
        <button
          className={`px-4 lg:px-8 py-4 rounded-full w-1/2 text-center ${!signUp ? 'bg-[#FBD5BD]' : 'bg-col'}`}
          onClick={clickLogin}
        >
          Log In
        </button>
        <button
          className={`py-4 px-4 rounded-full text-center w-1/2 ${signUp ? 'bg-[#FBD5BD] w-1/2 px-4 lg:px-6 py-4 ' : 'bg-white dark:bg-gray-900'}`}
          onClick={clickRegister}
        >
          Create an Account
        </button>
      </div>
      
      {alertMessage && (
        <div id="alert-2" className="flex items-center p-4 mb-4 text-red-800 rounded-lg bg-red-50 dark:text-red-400" role="alert">
          <svg className="flex-shrink-0 w-4 h-4" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="currentColor" viewBox="0 0 20 20">
            <path d="M10 .5a9.5 9.5 0 1 0 9.5 9.5A9.51 9.51 0 0 0 10 .5ZM9.5 4a1.5 1.5 0 1 1 0 3 1.5 1.5 0 0 1 0-3ZM12 15H8a1 1 0 0 1 0-2h1v-3H8a1 1 0 0 1 0-2h2a1 1 0 0 1 1 1v4h1a1 1 0 0 1 0 2Z" />
          </svg>
          <span className="sr-only">Info</span>
          <div className="ms-3 text-sm font-medium">{alertMessage}</div>
          <button
            onClick={() => setAlertMessage(null)}
            type="button"
            className="ms-auto -mx-1.5 -my-1.5 bg-red-50 text-red-500 rounded-lg focus:ring-2 focus:ring-red-400 p-1.5 hover:bg-red-200 inline-flex items-center justify-center h-8 w-8 dark:text-red-400"
            aria-label="Close"
          >
            <span className="sr-only">Close</span>
            <svg className="w-3 h-3" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 14 14">
              <path stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="m1 1 6 6m0 0 6 6M7 7l6-6M7 7l-6 6" />
            </svg>
          </button>
        </div>
      )}
      
      {signUp ? (
        <form className="flex w-[80%] flex-col items-center">
          <input
            type="text"
            value={lastName}
            onChange={(e) => setLastName(e.target.value)}
            name="lastName"
            className="mb-4 p-2 border border-[#8CB2DF] bg-[#181c2e] rounded-xl w-full text-white placeholder-gray-400"
            placeholder="Last Name"
            required
            style={inputStyles}
          />
          
          <input
            type="text"
            value={firstName}
            onChange={(e) => setFirstName(e.target.value)}
            name="firstName"
            className="mb-4 p-2 border border-[#8CB2DF] bg-[#181c2e] rounded-xl w-full text-white placeholder-gray-400"
            placeholder="First Name"
            required
            style={inputStyles}
          />
          
          <input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            name="mail"
            className="mb-4 p-2 border border-[#8CB2DF] bg-[#181c2e] rounded-xl w-full text-white placeholder-gray-400"
            placeholder="Email"
            required
            style={inputStyles}
          />
          
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            name="password"
            className="mb-4 p-2 border border-[#8CB2DF] bg-[#181c2e] rounded-xl w-full text-white placeholder-gray-400"
            placeholder="Password"
            required
            style={inputStyles}
          />
          
          <input
            type="password"
            value={confirmedPassword}
            onChange={(e) => setConfirmedPassword(e.target.value)}
            name="confirmedPassword"
            className="mb-4 p-2 border border-[#8CB2DF] bg-[#181c2e] rounded-xl w-full text-white placeholder-gray-400"
            placeholder="Confirm Password"
            required
            style={inputStyles}
          />
          
          <button
            onClick={handleSignUp}
            type="button"
            className="bg-[#FBD5BD] font-bold text-black py-3 px-5 rounded-lg text-base"
            disabled={loading}
          >
            {loading ? 'Creating...' : 'Create an Account'}
          </button>
        </form>
      ) : (
        <form className="flex w-[80%] flex-col items-center">
          <input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            name="mail"
            className="mb-4 p-2 border border-[#8CB2DF] bg-[#181c2e] rounded-xl w-full text-white placeholder-gray-400"
            placeholder="Email"
            required
            style={inputStyles}
          />
          
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            name="password"
            className="mb-4 p-2 border border-[#8CB2DF] bg-[#181c2e] rounded-xl w-full text-white placeholder-gray-400"
            placeholder="Password"
            required
            style={inputStyles}
          />
          
          <button
            onClick={handleLogin}
            type="button"
            className="bg-[#FBD5BD] font-bold text-black py-3 px-5 rounded-lg text-base"
            disabled={loading}
          >
            {loading ? 'Logging in...' : 'Log In'}
          </button>
        </form>
      )}
      
      <div className="flex w-full justify-center items-center gap-2 mt-6">
        <div className="border-t-2 border-black w-[35%] dark:border-white"></div>
        <span className="text-black dark:text-white">OR</span>
        <div className="border-t-2 border-black w-[35%] dark:border-white"></div>
      </div>
      
      <div className="flex items-center justify-center gap-4 my-4">
        <button className="bg-white border border-gray-300 p-2 rounded">
          <Image src="/images/Google.png" alt="Google" width={40} height={40} className="h-10" />
        </button>
        <button className="bg-white border border-gray-300 p-2 rounded">
          <Image src="/images/Microsoft.png" alt="Microsoft" width={40} height={40} className="h-10" />
        </button>
        <button className="bg-white border border-gray-300 p-2 rounded">
          <Image src="/images/Apple.png" alt="Apple" width={40} height={40} className="h-10" />
        </button>
      </div>
      
      <div>
        {signUp ? (
          <>
            <span className="text-sm text-center text-black dark:text-white">Already have an account? </span>
            <button className="text-[#8CB2DF] text-sm underline" onClick={clickLogin}>
              Log In
            </button>
          </>
        ) : (
          <button className="text-[#8CB2DF] text-sm underline">
            Forgot password?
          </button>
        )}
      </div>
      
      <div className="text-[#8CB2DF] bottom-0 text-xs lg:text-sm xl:text-lg text-center py-4 w-full pt-10">
        <Link href="/" className="mx-4">Terms of Use</Link> | <Link href="/" className="mx-4">Privacy Policy</Link>
      </div>
    </div>
  );
}
'use client';

import React, { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { useAuth } from '@/lib/auth/use-auth';
import { organizationApi } from '@/lib/api/api-client';

export default function CreateOrgForm() {
  const router = useRouter();
  const { user } = useAuth();
  
  // Form state
  const [name, setName] = useState('');
  const [alertMessage, setAlertMessage] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  
  // Check if user is already in organization
  useEffect(() => {
    if (user?.organization) {
      router.push('/dashboard');
    }
  }, [user, router]);
  
  // Handle organization creation
  const createOrganization = async (e: React.FormEvent) => {
    e.preventDefault();
    if (name === '') {
      setAlertMessage('Please enter a name');
      return;
    }
    
    try {
      setLoading(true);
      const res = await organizationApi.create({ name });
      
      if (res.success) {
        // Set cookie for middleware
        document.cookie = 'has_organization=true; path=/';
        router.push('/dashboard');
      }
    } catch (err: any) {
      if (err.message?.includes('exists')) {
        setAlertMessage('An organization with this name already exists. Please choose another name.');
      } else {
        setAlertMessage('An error occurred while creating your organization. Please try again later.');
      }
    } finally {
      setLoading(false);
    }
  };
  
  return (
    <div className="bg-white py-8 px-4 shadow sm:rounded-lg sm:px-10 dark:bg-gray-800">
      <form className="space-y-6" onSubmit={createOrganization}>
        <label
          htmlFor="name"
          className="block text-sm font-medium text-gray-700 dark:text-gray-300"
        >
          Organization name:
        </label>
        
        {alertMessage && (
          <div
            id="alert-2"
            className="flex items-center p-4 mb-4 text-red-800 rounded-lg bg-red-50 dark:text-red-400 dark:bg-red-900"
            role="alert"
          >
            <svg
              className="flex-shrink-0 w-4 h-4"
              aria-hidden="true"
              xmlns="http://www.w3.org/2000/svg"
              fill="currentColor"
              viewBox="0 0 20 20"
            >
              <path d="M10 .5a9.5 9.5 0 1 0 9.5 9.5A9.51 9.51 0 0 0 10 .5ZM9.5 4a1.5 1.5 0 1 1 0 3 1.5 1.5 0 0 1 0-3ZM12 15H8a1 1 0 0 1 0-2h1v-3H8a1 1 0 0 1 0-2h2a1 1 0 0 1 1 1v4h1a1 1 0 0 1 0 2Z" />
            </svg>
            <span className="sr-only">Info</span>
            <div className="ms-3 text-sm font-medium">
              {alertMessage}
            </div>
            <button
              onClick={() => setAlertMessage(null)}
              type="button"
              className="ms-auto -mx-1.5 -my-1.5 bg-red-50 text-red-500 rounded-lg focus:ring-2 focus:ring-red-400 p-1.5 hover:bg-red-200 inline-flex items-center justify-center h-8 w-8 dark:text-red-400 dark:bg-red-900"
              aria-label="Close"
            >
              <span className="sr-only">Close</span>
              <svg className="w-3 h-3" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 14 14">
                <path stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="m1 1 6 6m0 0 6 6M7 7l6-6M7 7l-6 6" />
              </svg>
            </button>
          </div>
        )}
        
        <input
          value={name}
          onChange={(e) => setName(e.target.value)}
          className="appearance-none rounded-md relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-indigo-500 dark:focus:border-indigo-500 dark:shadow-sm-light"
          placeholder="Enter the name of your organization"
        />
        
        <button
          type="submit"
          disabled={loading}
          className="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 dark:bg-indigo-500 dark:hover:bg-indigo-600 dark:focus:ring-indigo-500"
        >
          {loading ? 'Creating...' : 'Create'}
        </button>
      </form>
      
      <div className="mt-6">
        <div className="relative">
          <div className="absolute inset-0 flex items-center">
            <div className="w-full border-t border-gray-300 dark:border-gray-700"></div>
          </div>
          <div className="relative flex justify-center text-sm">
            <span className="px-2 bg-white text-gray-500 dark:text-gray-400 dark:bg-gray-800">
              OR
            </span>
          </div>
        </div>
        
        <Link
          href="/organization/join"
          className="mt-4 gap-3 text-blue-500 w-full flex items-center justify-center text-sm dark:text-blue-400"
        >
          Join an existing organization
        </Link>
      </div>
    </div>
  );
}
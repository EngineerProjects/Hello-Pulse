// API client utilities
'use client';

// Define the base API URL - adjust based on your environment
const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:5000';

// Reusable fetch function with proper error handling
async function fetchApi<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const url = `${API_URL}${endpoint}`;
  
  // Default options for all requests
  const defaultOptions: RequestInit = {
    credentials: 'include', // Required for cookies
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
  };
  
  const response = await fetch(url, { ...defaultOptions, ...options });
  
  // Parse JSON response
  const data = await response.json();
  
  // Check if the request was successful
  if (!response.ok) {
    // Create error object with API error message
    const error = new Error(data.error || 'An unknown error occurred');
    (error as any).status = response.status;
    (error as any).data = data;
    throw error;
  }
  
  return data as T;
}

// Auth API methods
export const authApi = {
  // Register a new user
  register: async (userData: {
    email: string;
    first_name: string;
    last_name: string;
    phone: string;
    address: string;
    password: string;
  }) => {
    return fetchApi<{ success: boolean; message: string; user: any }>('/register', {
      method: 'POST',
      body: JSON.stringify(userData),
    });
  },
  
  // Login user
  login: async (credentials: { email: string; password: string }) => {
    return fetchApi<{ success: boolean; message: string; user: any }>('/login', {
      method: 'POST',
      body: JSON.stringify(credentials),
    });
  },
  
  // Get current user
  getMe: async () => {
    return fetchApi<{ success: boolean; user: any }>('/me');
  },
  
  // Logout user
  logout: async () => {
    return fetchApi<{ success: boolean; message: string }>('/logout', {
      method: 'POST',
    });
  },
};

// Organization API methods
export const organizationApi = {
  // Create new organization
  create: async (data: { name: string }) => {
    return fetchApi<{ success: boolean; message: string; organization: any }>('/organizations', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  },
  
  // Join organization with invite code
  join: async (data: { code: string }) => {
    return fetchApi<{ success: boolean; message: string }>('/organizations/join', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  },
  
  // Get invite codes
  getInviteCodes: async () => {
    return fetchApi<{ success: boolean; codes: any[] }>('/organizations/invite-codes');
  },
};

export default {
  auth: authApi,
  organization: organizationApi,
};
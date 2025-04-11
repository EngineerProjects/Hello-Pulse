// User model
export interface User {
    id: string;
    email: string;
    firstName: string;
    lastName: string;
    phone: string;
    address: string;
    role: string;
    createdAt: string;
    organization: Organization | null;
  }
  
  // Organization model
  export interface Organization {
    id: string;
    name: string;
  }
  
  // Invite code model
  export interface InviteCode {
    id: string;
    code: string;
    expirationTimeMs: number;
  }
  
  // API response types
  export interface ApiResponse {
    success: boolean;
    message?: string;
    error?: string;
  }
  
  export interface LoginResponse extends ApiResponse {
    user?: User;
  }
  
  export interface RegisterResponse extends ApiResponse {
    user?: User;
  }
  
  export interface MeResponse extends ApiResponse {
    user?: User;
  }
  
  export interface CreateOrganizationResponse extends ApiResponse {
    organization?: Organization;
  }
  
  export interface InviteCodesResponse extends ApiResponse {
    codes: InviteCode[];
  }
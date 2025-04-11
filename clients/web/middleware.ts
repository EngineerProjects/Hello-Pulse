import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';

// This function can be marked `async` if using `await` inside
export function middleware(request: NextRequest) {
  // Get the path and check for authentication via cookie
  const path = request.nextUrl.pathname;
  const isAuthenticated = !!request.cookies.get('token')?.value;
  
  // Get current user organization status - could be in a session/cookie
  // We'll check the path of /organization/ in cookies as a workaround
  // In a real app, you might have a dedicated cookie or session for this
  const hasOrganization = request.cookies.get('has_organization')?.value === 'true';

  // Public paths that don't require authentication
  const isPublicPath = 
    path === '/login' || 
    path === '/register' || 
    path === '/' ||
    path.startsWith('/_next') || 
    path.startsWith('/api') ||
    path.includes('/organization');
  
  // Handle paths that require authentication (dashboard)
  if (path.startsWith('/dashboard')) {
    // If user is not authenticated, redirect to login
    if (!isAuthenticated) {
      return NextResponse.redirect(new URL('/login', request.url));
    }
    
    // If user is authenticated but has no organization, redirect to organization flow
    if (isAuthenticated && !hasOrganization) {
      return NextResponse.redirect(new URL('/organization/join', request.url));
    }
  }

  // If user is trying to access login/register but is already authenticated
  if ((path === '/login' || path === '/register') && isAuthenticated) {
    // If user has organization, go to dashboard, otherwise organization flow
    if (hasOrganization) {
      return NextResponse.redirect(new URL('/dashboard', request.url));
    } else {
      return NextResponse.redirect(new URL('/organization/join', request.url));
    }
  }

  // Organization join/create routes should only be accessible by authenticated users
  if ((path === '/organization/join' || path === '/organization/create') && !isAuthenticated) {
    return NextResponse.redirect(new URL('/login', request.url));
  }
  
  return NextResponse.next();
}

// See "Matching Paths" below to learn more
export const config = {
  matcher: [
    // Protected routes that need authentication
    '/dashboard/:path*',
    
    // Auth routes for redirecting authenticated users
    '/login',
    '/register',
    
    // Organization routes for authenticated users only
    '/organization/:path*',
  ],
};
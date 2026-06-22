// apps/analytics-dashboard/middleware.ts

import { NextRequest, NextResponse } from "next/server";

const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_GATEWAY_URL ??
  "http://localhost:8080";

export async function middleware(
  request: NextRequest
) {
  if (
    process.env.NODE_ENV === "development" &&
    request.nextUrl.searchParams.has("__preview")
  ) {
    return NextResponse.next();
  }

  const token =
    request.cookies.get("access_token")?.value;

  if (!token) {
    const loginUrl = new URL(
      "/login",
      request.url
    );

    loginUrl.searchParams.set(
      "redirect",
      request.nextUrl.pathname
    );

    return NextResponse.redirect(loginUrl);
  }

  try {
    const response = await fetch(
      `${API_BASE_URL}/v1/auth/me`,
      {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }
    );

    if (!response.ok) {
      throw new Error("Unauthorized");
    }

    const user = await response.json();

    const role =
      user.role ??
      user.data?.role ??
      user.user?.role;

    if (
      !["admin", "super_admin"].includes(role)
    ) {
      const loginUrl = new URL(
        "/login",
        request.url
      );

      loginUrl.searchParams.set(
        "error",
        "unauthorized"
      );

      return NextResponse.redirect(loginUrl);
    }

    return NextResponse.next();
  } catch {
    const loginUrl = new URL(
      "/login",
      request.url
    );

    loginUrl.searchParams.set(
      "redirect",
      request.nextUrl.pathname
    );

    return NextResponse.redirect(loginUrl);
  }
}

export const config = {
  matcher: [
    "/((?!login|forgot-password|reset-password|_next|favicon\\.png|logo-|icon-|api).*)",
  ],
};
//import { NextResponse, type NextRequest } from 'next/server';
//import { createMiddlewareClient } from '@/lib/api-client';

/**
 * Middleware that protects all admin routes.
 * Redirects to /login if:
 *  - No valid Supabase session
 *  - User does not have admin or super_admin role
 
export async function middleware(request: NextRequest) {
  // Dev-only escape hatch for design previews: /foo?__preview=1
  // Gated by NODE_ENV so it can never run in production builds.
  if (
    process.env.NODE_ENV === 'development' &&
    request.nextUrl.searchParams.has('__preview')
  ) {
    return NextResponse.next();
  }

  const { supabase, response } = createMiddlewareClient(request);

  // Check for valid session
  const { data: { user }, error } = await supabase.auth.getUser();

  if (error || !user) {
    const loginUrl = new URL('/login', request.url);
    loginUrl.searchParams.set('redirect', request.nextUrl.pathname);
    return NextResponse.redirect(loginUrl);
  }

  // Check admin role
  const { data: userData } = await supabase
    .from('users')
    .select('role')
    .eq('id', user.id)
    .single();

  if (!userData || !['admin', 'super_admin'].includes(userData.role)) {
    const loginUrl = new URL('/login', request.url);
    loginUrl.searchParams.set('error', 'unauthorized');
    return NextResponse.redirect(loginUrl);
  }

  return response;
}

export const config = {
  matcher: [
    /*
     * Match all routes except:
     * - /login (auth page)
     * - /_next (Next.js internals)
     * - /favicon.png, /logo-*, /icon-* (static assets)
     * - /api (API routes if any)
     *
    '/((?!login|forgot-password|reset-password|_next|favicon\\.png|logo-|icon-|api).*)',
  ],
};
*/

export function getJWT(): string | null {
  if (typeof window === 'undefined') return null; // Guard for server-side rendering
  return localStorage.getItem('token');
}

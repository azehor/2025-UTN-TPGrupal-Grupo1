// Helper para resolver la URL base de la API.
// Orden de preferencia:
// 1) runtime config servido por el backend: window.__ENV__.VITE_API_BASE_URL
// 2) build-time Vite env: import.meta.env.VITE_API_BASE_URL
// 3) fallback: same origin + /v1
export function getApiBase(): string {
  try {
    const runtime = (window as any).__ENV?.VITE_API_BASE_URL;
    if (runtime && typeof runtime === 'string' && runtime.length) return runtime;
  } catch (e) {
    // ignore
  }

  const build = (import.meta as any).env?.VITE_API_BASE_URL;
  if (build && typeof build === 'string' && build.length) return build;

  if (typeof window !== 'undefined' && window.location) {
    return `${window.location.origin}/v1`;
  }

  return '/v1';
}

export default getApiBase;

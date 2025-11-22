const TOKEN_KEY = 'timmach_token';

function safeLocalStorage(): Storage | null {
  if (typeof window === 'undefined') return null;
  try {
    return window.localStorage;
  } catch (_err) {
    return null;
  }
}

export function getToken(): string | null {
  const store = safeLocalStorage();
  return store ? store.getItem(TOKEN_KEY) : null;
}

export function setToken(token: string) {
  const store = safeLocalStorage();
  if (store) store.setItem(TOKEN_KEY, token);
}

export function clearToken() {
  const store = safeLocalStorage();
  if (store) store.removeItem(TOKEN_KEY);
}

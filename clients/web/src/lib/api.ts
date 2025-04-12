const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:3001/api";

export async function fetchJson<T>(
  input: RequestInfo,
  init?: RequestInit
): Promise<T> {
  const response = await fetch(input, init);
  
  if (!response.ok) {
    const error = new Error(response.statusText);
    try {
      const data = await response.json();
      if (data.error) {
        (error as any).message = data.error;
      }
    } catch {
      // If the response is not JSON, use the status text
    }
    throw error;
  }
  
  return response.json();
}

export async function login(email: string, password: string) {
  return fetchJson(`${API_BASE_URL}/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
  });
}

export async function register(firstName: string, lastName: string, email: string, password: string) {
  return fetchJson(`${API_BASE_URL}/register`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      first_name: firstName,
      last_name: lastName,
      email,
      password,
      phone: "0",
      address: "0",
    }),
  });
}
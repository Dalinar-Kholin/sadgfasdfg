/**
 * Fetch wrapper that adds Authorization header from local storage.
 * @param url The URL to send the request to.
 * @param options Fetch options.
 * @returns A promise that resolves to the response.
 */
export default async function fetchWithAuth(url: string, options: RequestInit = {}): Promise<Response> {
    const token = localStorage.getItem('accessToken');

    options.headers = options.headers || {};

    // Cast headers to the correct type
    const headers = new Headers(options.headers as HeadersInit);

    // Add the Authorization header if the token is available
    if (token) {
        headers.set('Authorization', `${token}`);
    }

    // Convert Headers back to the appropriate type for options
    const headersObject: Record<string, string> = {};
    headers.forEach((value, key) => {
        headersObject[key] = value;
    });

    // Set the headers back to options
    options.headers = headersObject;

    // Perform the fetch request
    const response = await fetch(url, options);

    // Handle errors or return the response
    if (!response.ok) {
        // Handle error responses (you can customize this as needed)
        const errorText = await response.text();
        throw new Error(`Fetch error: ${response.status} ${response.statusText} - ${errorText}`);
    }

    return response;
}
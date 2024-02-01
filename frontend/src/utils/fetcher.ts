type ModifiedResponse<T = Record<string, unknown>> = {
  statusCode: number;
  message: string;
  data?: T;
};

async function Post<T = Record<string, unknown>>(
  url: string,
  options?: {
    form?: object;
    token?: string | null;
    headers?: Headers;
  }
): Promise<ModifiedResponse<T>> {
  try {
    const headers = options?.headers
      ? options.headers
      : new Headers({ 'Content-Type': 'application/json' });
    const body = options && options.form ? options.form : {};

    if (options && options.token) {
      headers.append('Authorization', `Bearer ${options.token}`);
    }

    const response = await fetch(url, {
      method: 'POST',
      headers: headers,
      body: JSON.stringify(body),
    });

    return await handleResponse(response);
  } catch (error) {
    return (await handleError(error)) as ModifiedResponse<T>;
  }
}

async function Patch<T = Record<string, unknown>>(
  url: string,
  options?: {
    form?: object;
    token?: string | null;
    headers?: Headers;
  }
): Promise<ModifiedResponse<T>> {
  try {
    const headers = options?.headers
      ? options.headers
      : new Headers({ 'Content-Type': 'application/json' });
    const body = options && options.form ? options.form : {};

    if (options && options.token) {
      headers.append('Authorization', `Bearer ${options.token}`);
    }

    const response = await fetch(url, {
      method: 'PATCH',
      headers: headers,
      body: JSON.stringify(body),
    });

    return await handleResponse(response);
  } catch (error) {
    return (await handleError(error)) as ModifiedResponse<T>;
  }
}

async function Put<T = Record<string, unknown>>(
  url: string,
  options?: {
    form?: object;
    token?: string | null;
    headers?: Headers;
  }
): Promise<ModifiedResponse<T>> {
  try {
    const headers = options?.headers
      ? options.headers
      : new Headers({ 'Content-Type': 'application/json' });
    const body = options && options.form ? options.form : {};

    if (options && options.token) {
      headers.append('Authorization', `Bearer ${options.token}`);
    }

    const response = await fetch(url, {
      method: 'PUT',
      headers: headers,
      body: JSON.stringify(body),
    });

    return await handleResponse(response);
  } catch (error) {
    return (await handleError(error)) as ModifiedResponse<T>;
  }
}

async function Delete<T = Record<string, unknown>>(
  url: string,
  options?: {
    form?: object;
    token?: string | null;
    headers?: Headers;
  }
): Promise<ModifiedResponse<T>> {
  try {
    const headers = options?.headers
      ? options.headers
      : new Headers({ 'Content-Type': 'application/json' });
    const body = options && options.form ? options.form : {};

    if (options && options.token) {
      headers.append('Authorization', `Bearer ${options.token}`);
    }

    const response = await fetch(url, {
      method: 'DELETE',
      headers: headers,
      body: JSON.stringify(body),
    });

    return await handleResponse(response);
  } catch (error) {
    return (await handleError(error)) as ModifiedResponse<T>;
  }
}

async function Get<T = Record<string, unknown>>(
  url: string,
  options?: { token?: string | null; headers?: Headers }
): Promise<ModifiedResponse<T>> {
  try {
    const headers = options?.headers
      ? options.headers
      : new Headers({ 'Content-Type': 'application/json' });

    if (options?.token) {
      headers.append('Authorization', `Bearer ${options.token}`);
    }

    const response = await fetch(url, {
      method: 'GET',
      headers: headers,
    });

    return await handleResponse(response);
  } catch (error) {
    return (await handleError(error)) as ModifiedResponse<T>;
  }
}

async function Upload<T = Record<string, unknown>>(
  url: string,
  options: {
    token?: string;
    headers?: Headers;
    form: FormData;
    method?: 'POST' | 'PATCH';
  }
) {
  try {
    const headers = options?.headers
      ? options.headers
      : new Headers();

    if (options && options.token) {
      headers.append('Authorization', `Bearer ${options.token}`);
    }

    const response = await fetch(url, {
      method: options.method ?? 'POST',
      headers: headers,
      body: options.form,
    });

    return await handleResponse(response);
  } catch (error) {
    return (await handleError(error)) as ModifiedResponse<T>;
  }
}

async function Download<T = Record<string, unknown>>(
  url: string,
  options: {
    token?: string;
    headers?: Headers;
  }
) {
  try {
    const headers = options?.headers
      ? options.headers
      : new Headers();

    if (options && options.token) {
      headers.append('Authorization', `Bearer ${options.token}`);
    }

    const response = await fetch(url, {
      method: 'GET',
      headers: headers,
    });

    const fileName = response.headers
      .get('Content-Disposition')
      ?.split('filename=')
      .pop()
      ?.replace(/"/g, '');

    response.blob().then((blob) => {
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = fileName as string;
      document.body.appendChild(a);
      a.click();
      a.remove();
    });

    return response;
  } catch (error) {
    return (await handleError(error)) as ModifiedResponse<T>;
  }
}

async function handleResponse<T = Record<string, unknown>>(
  response: Response
): Promise<ModifiedResponse<T>> {
  if (response.ok) {
    try {
      const result = await response.json();

      return {
        statusCode: result.statusCode,
        message: result.message,
        data: result.data,
      };
    } catch {
      return { statusCode: 0, message: '' };
    }
  } else {
    throw response;
  }
}

async function handleError(
  error: unknown
): Promise<ModifiedResponse> {
  if (error instanceof Response) {
    const errorResponse = await error.json();
    try {
      return {
        statusCode: errorResponse.statusCode,
        data: errorResponse.data,
        message: errorResponse.message,
      };
    } catch {
      return { statusCode: error.status, message: '' };
    }
  } else {
    return { statusCode: 0, message: 'Failed to Fetch!' };
  }
}

const fetchers = {
  Post,
  Get,
  Patch,
  Put,
  Delete,
  Upload,
  Download,
};

export default fetchers;

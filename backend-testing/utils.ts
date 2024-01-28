import { BACKEND_ENDPOINT, httpClient } from './constants';
import axios from 'axios';

export async function getAuthenticatedClient() {
  const result = await httpClient.post('/auth/signin', {
    username: 'userA1234',
    password: '12345678',
  });

  return axios.create({
    baseURL: BACKEND_ENDPOINT,
    headers: {
      Authorization: `Bearer ${result.data.data.token}`,
    },
  });
}

export async function getAuthenticatedClientWithRefreshToken() {
  const result = await httpClient.post('/auth/signin', {
    username: 'userA1234',
    password: '12345678',
  });

  return {
    client: axios.create({
      baseURL: BACKEND_ENDPOINT,
      headers: {
        Authorization: `Bearer ${result.data.data.token}`,
      },
    }),
    refreshToken: result.data.data.refreshToken as string,
  };
}

import { PrismaClient } from '@prisma/client';
import { BACKEND_ENDPOINT, httpClient } from './constants';
import axios from 'axios';

export const prismaClient = new PrismaClient();

type GetAuthenticatedClientParam = {
  includeProjectId?: boolean;
};
export async function getAuthenticatedClient({
  includeProjectId = false,
}: GetAuthenticatedClientParam) {
  const result = await httpClient.post('/auth/signin', {
    username: 'userA1234',
    password: '12345678',
  });

  return axios.create({
    baseURL: BACKEND_ENDPOINT,
    headers: {
      Authorization: `Bearer ${result.data.data.token}`,
      ...(includeProjectId && {
        'x-project-id': '246bb085-8ccc-4def-ac78-dc2ad5c7760b',
      }),
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

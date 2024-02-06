import { PrismaClient } from '@prisma/client';
import { BACKEND_ENDPOINT, httpClient } from './constants';
import axios from 'axios';
import { config } from 'dotenv';
import { createClient } from '../graphql/generated';
config();

export const graphqlClient = createClient({
  url: process.env.BACKEND_ENDPOINT,
});

export const prismaClient = new PrismaClient();

type GetAuthenticatedClientParam = {
  includeProjectId?: boolean;
};
export async function getAuthenticatedClient({
  includeProjectId = false,
}: GetAuthenticatedClientParam) {
  const result = await graphqlClient.mutation({
    signin: {
      __scalar: true,
      __args: {
        username: 'userA1234',
        password: '12345678',
      },
    },
  });

  return createClient({
    url: process.env.BACKEND_ENDPOINT,
    headers: {
      Authorization: `Bearer ${result.signin.token}`,
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
      paramsSerializer: {
        indexes: null,
      },
      baseURL: BACKEND_ENDPOINT,
      headers: {
        Authorization: `Bearer ${result.data.data.token}`,
      },
    }),
    refreshToken: result.data.data.refreshToken as string,
  };
}

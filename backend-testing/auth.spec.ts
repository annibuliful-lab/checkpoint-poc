import { nanoid } from 'nanoid';
import { httpClient } from './constants';
import {
  getAuthenticatedClient,
  getAuthenticatedClientWithRefreshToken,
} from './utils';

describe('Auth', () => {
  it('calls signout with token', async () => {
    const client = await getAuthenticatedClient();

    const result = await client.post('/auth/signout');

    expect(result.data.message).toEqual('success');
  });

  it('calls signup', async () => {
    const result = await httpClient.post('/auth/signup', {
      username: nanoid(),
      password: nanoid(),
    });

    expect(result.data.data.id).toBeDefined();
    expect(result.data.data.createdAt).toBeDefined();
    expect(result.data.data.updatedAt).toBeNull();
  });

  it('gets token with refresh token', async () => {
    const { refreshToken } =
      await getAuthenticatedClientWithRefreshToken();

    const result = await httpClient.post(
      '/auth/refresh-token',
      null,
      {
        headers: {
          Authorization: `Bearer ${refreshToken}`,
        },
      }
    );
    const responseData = result.data.data;

    expect(responseData.userId).toBeDefined();
    expect(responseData.token).toBeDefined();
    expect(responseData.refreshToken).toBeDefined();
    expect(result.data.message).toEqual('signed');
  });

  it('throws error when call sigin API with random user and password', (done) => {
    httpClient
      .post('/auth/signin', {
        username: nanoid(),
        password: nanoid(),
      })
      .catch((error) => {
        expect(error.response.data.message).toEqual(
          'username or password is incorrect'
        );
        done();
      });
  });

  it('gets authentication response when sigin with correct username and password', async () => {
    const result = await httpClient.post('/auth/signin', {
      username: 'userA1234',
      password: '12345678',
    });

    expect(result.data.data.userId).toBeDefined();
    expect(result.data.data.token).toBeDefined();
    expect(result.data.data.refreshToken).toBeDefined();
  });
});

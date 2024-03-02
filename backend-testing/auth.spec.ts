import { nanoid } from 'nanoid';
import { graphqlClient } from './utils/utils';
import { Client } from './graphql/generated';

describe('Authentication', () => {
  let client: Client;
  beforeAll(() => {
    client = graphqlClient;
  });

  it('signup', async () => {
    const username = nanoid();
    const response = await client.mutation({
      signup: {
        id: true,
        username: true,
        __args: {
          username,
          password: '12345678',
        },
      },
    });
    expect(response.signup.id).toBeDefined();
    expect(response.signup.username).toEqual(username);
  });

  it('signin', async () => {
    const response = await client.mutation({
      signin: {
        __scalar: true,
        __args: {
          username: 'userA1234',
          password: '12345678',
        },
      },
    });

    expect(response.signin.refreshToken).toBeDefined();
    expect(response.signin.token).toBeDefined();
    expect(response.signin.userId).toEqual(
      'af466ea9-04ba-432e-98e0-3b8787dcda41'
    );
  });
});

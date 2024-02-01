import { nanoid } from 'nanoid';
import { getAuthenticatedClient } from './utils';

export async function createProject() {
  const client = await getAuthenticatedClient({});
  const response = await client.post('/projects', {
    title: nanoid(),
  });

  const responseData = response.data;

  return {
    id: responseData.data.id as string,
    title: responseData.data.title as string,
    message: responseData.message as string,
  };
}

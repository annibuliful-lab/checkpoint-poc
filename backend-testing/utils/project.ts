import { nanoid } from 'nanoid';
import { getAuthenticatedClient } from './utils';

export async function createProject() {
  const client = await getAuthenticatedClient({});
  const response = await client.mutation({
    createProject: {
      __scalar: true,
      __args: {
        title: nanoid(),
      },
    },
  });

  return {
    id: response.createProject.id,
    title: response.createProject.title,
  };
}

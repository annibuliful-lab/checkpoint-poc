import { AxiosInstance, AxiosResponse } from 'axios';
import { getAuthenticatedClient } from './utils/utils';
import { nanoid } from 'nanoid';
import { createProject } from './utils/project';
import { v4 } from 'uuid';

describe('Project', () => {
  let client: AxiosInstance;

  beforeAll(async () => {
    client = await getAuthenticatedClient({});
  });

  it('creates new', async () => {
    const title = nanoid();
    const response = await client.post('/projects', {
      title,
    });

    const responseData = response.data;
    expect(responseData.data.id).toBeDefined();
    expect(responseData.data.title).toEqual(title);
    expect(responseData.message).toEqual('created');
  });

  it('updates', async () => {
    const project = await createProject();
    const title = nanoid();
    const updatedResponse = await client.patch(
      `/projects/${project.id}`,
      {
        title,
      }
    );

    const updatedResponseData = updatedResponse.data;
    expect(updatedResponseData.data.id).toEqual(project.id);
    expect(updatedResponseData.data.title).toEqual(title);
  });

  it('throws error when update by wrong id', async () => {
    try {
      await client.patch(`/projects/${v4()}`, {
        title: nanoid(),
      });
    } catch (error: any) {
      expect(error.response.status).toEqual(403);
      expect(error.response.data.message).toEqual(
        'forbidden operation'
      );
    }
  });

  it('gets by id', async () => {
    const project = await createProject();
    const getResponse = await client.get(`/projects/${project.id}`);

    const responseData = getResponse.data.data;

    expect(responseData.id).toEqual(project.id);
    expect(responseData.title).toEqual(project.title);
  });

  it('throws when get by wrong id', async () => {
    try {
      await client.get(`/projects/${v4()}`);
    } catch (error: any) {
      expect(error.response.status).toEqual(403);
      expect(error.response.data.message).toEqual(
        'forbidden operation'
      );
    }
  });
});

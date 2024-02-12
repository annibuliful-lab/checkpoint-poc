import { AxiosInstance, AxiosResponse } from 'axios';
import { getAuthenticatedClient } from './utils/utils';
import { nanoid } from 'nanoid';
import { createProject } from './utils/project';
import { v4 } from 'uuid';
import { Client } from './graphql/generated';

describe('Project', () => {
  let client: Client;

  beforeAll(async () => {
    client = await getAuthenticatedClient({});
  });

  it('creates new', async () => {
    const title = nanoid();
    const response = await client.mutation({
      createProject: {
        __scalar: true,
        __args: {
          title,
        },
      },
    });

    expect(response.createProject.id).toBeDefined();
    expect(response.createProject.title).toEqual(title);
  });

  it('updates', async () => {
    const project = await createProject();
    const title = nanoid();
    const updatedResponse = await client.mutation({
      updateProject: {
        __scalar: true,
        __args: {
          id: project.id,
          title,
        },
      },
    });

    expect(updatedResponse.updateProject.id).toEqual(project.id);
    expect(updatedResponse.updateProject.title).toEqual(title);
  });

  it('throws error when update by wrong id', async () => {
    try {
      await client.mutation({
        updateProject: {
          __scalar: true,
          __args: {
            id: v4(),
            title: nanoid(),
          },
        },
      });
    } catch (error: any) {
      expect(error.errors[0].extensions.code).toEqual(
        'forbidden operation'
      );
    }
  });

  it('gets by id', async () => {
    const project = await createProject();
    const getResponse = await client.query({
      getProjectById: {
        __scalar: true,
        __args: {
          id: project.id,
        },
      },
    });

    expect(getResponse.getProjectById.id).toEqual(project.id);
    expect(getResponse.getProjectById.title).toEqual(project.title);
  });

  it('throws when get by wrong id', async () => {
    expect(
      client.query({
        getProjectById: {
          __scalar: true,
          __args: {
            id: v4(),
          },
        },
      })
    ).rejects.toThrow();
  });
});

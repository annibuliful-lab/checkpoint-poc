import { AxiosInstance } from 'axios';
import { createProject } from './utils/project';
import { getAuthenticatedClient, prismaClient } from './utils/utils';
import { nanoid } from 'nanoid';
import { v4 } from 'uuid';
import { PROJECT_ID } from './utils/constants';

describe('Project role', () => {
  let client: AxiosInstance;

  beforeAll(async () => {
    client = await getAuthenticatedClient({ includeProjectId: true });
  });

  it('creates', async () => {
    const permission = await prismaClient.permission.create({
      data: {
        id: v4(),
        subject: nanoid(),
        action: 'CREATE',
      },
    });
    const title = nanoid();
    const response = await client.post(`/projects/roles`, {
      title,
      permissionIds: [permission.id],
    });

    const responseData = response.data;
    expect(responseData.data.permissions).toHaveLength(1);
    expect(responseData.data.permissions[0]).toEqual(
      expect.objectContaining(permission)
    );
    expect(responseData.data.id).toBeDefined();
    expect(responseData.data.title).toEqual(title);
    expect(responseData.message).toEqual('created');
  });

  it('updates', async () => {
    const permission = await prismaClient.permission.create({
      data: {
        id: v4(),
        subject: nanoid(),
        action: 'CREATE',
      },
    });
    const title = nanoid();
    const response = await client.post(`/projects/roles`, {
      title,
      permissionIds: [permission.id],
    });
    const responseData = response.data;
    const newTitle = nanoid();
    const newPermission = await prismaClient.permission.create({
      data: {
        id: v4(),
        subject: nanoid(),
        action: 'CREATE',
      },
    });

    const updatedResponse = await client.patch(
      `/projects/roles/${responseData.data.id}`,
      {
        title: newTitle,
        permissionIds: [newPermission.id],
      }
    );
    const updatedResponseData = updatedResponse.data;
    expect(updatedResponseData.data.title).toEqual(newTitle);
    expect(updatedResponseData.data.permissions).toHaveLength(1);
    expect(newPermission).toEqual(
      updatedResponseData.data.permissions[0]
    );
  });

  it('gets by search', async () => {
    const permission = await prismaClient.permission.create({
      data: {
        id: v4(),
        subject: nanoid(),
        action: 'CREATE',
      },
    });
    const title = nanoid();
    await client.post(`/projects/roles`, {
      title,
      permissionIds: [permission.id],
    });

    const getRoleResponse = await client.get(
      `/projects/roles?search=${title}`
    );
    const responseData = getRoleResponse.data;

    expect(responseData.data.length).toEqual(1);
    expect(responseData.data.every((el: any) => el.title === title));
    expect(
      responseData.data.every(
        (el: any) => el.projectId === PROJECT_ID
      )
    );
  });

  it('gets by project id', async () => {
    const permission = await prismaClient.permission.create({
      data: {
        id: v4(),
        subject: nanoid(),
        action: 'CREATE',
      },
    });
    const title = nanoid();
    await client.post(`/projects/roles`, {
      title,
      permissionIds: [permission.id],
    });

    const getRoleResponse = await client.get(
      `/projects/roles?search=Admin`
    );
    const responseData = getRoleResponse.data;

    expect(responseData.data.length).toBeGreaterThanOrEqual(1);
    expect(responseData.data.some((el: any) => el.title === title));
    expect(
      responseData.data.every(
        (el: any) => el.projectId === PROJECT_ID
      )
    );
  });

  it('gets by id', async () => {
    const permission = await prismaClient.permission.create({
      data: {
        id: v4(),
        subject: nanoid(),
        action: 'CREATE',
      },
    });
    const title = nanoid();
    const createdRoleResponse = await client.post(`/projects/roles`, {
      title,
      permissionIds: [permission.id],
    });

    const createdRoleResponseData = createdRoleResponse.data;

    const getRoleResponse = await client.get(
      `/projects/roles/${createdRoleResponseData.data.id}`
    );
    const responseData = getRoleResponse.data;

    expect(responseData.data.title).toEqual(title);
    expect(responseData.data.permissions).toHaveLength(1);
    expect(responseData.data.permissions[0].id).toEqual(
      permission.id
    );
  });

  it('throws err when get by wrong id', async () => {
    try {
      const getRoleResponse = await client.get(
        `/projects/roles/${v4()}`
      );
    } catch (error: any) {
      expect(error.response.status).toEqual(404);
      expect(error.response.data.message).toEqual('id not found');
    }
  });
});

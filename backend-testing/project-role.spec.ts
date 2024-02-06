import { AxiosInstance } from 'axios';
import { createProject } from './utils/project';
import { getAuthenticatedClient, prismaClient } from './utils/utils';
import { nanoid } from 'nanoid';
import { v4 } from 'uuid';
import { PROJECT_ID } from './utils/constants';
import { Client } from './graphql/generated';

describe('Project role', () => {
  let client: Client;

  beforeAll(async () => {
    client = await getAuthenticatedClient({ includeProjectId: true });
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
    const created = await client.mutation({
      createProjectRole: {
        __scalar: true,
        __args: {
          title,
          permissionIds: [permission.id],
        },
      },
    });

    const response = await client.query({
      getProjectRoles: {
        __scalar: true,
        __args: {
          search: title,
          limit: 20,
          skip: 0,
        },
      },
    });
    expect(response.getProjectRoles).toHaveLength(1);
    expect(response.getProjectRoles[0]).toEqual(
      created.createProjectRole
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
    const created = await client.mutation({
      createProjectRole: {
        __scalar: true,
        __args: {
          title,
          permissionIds: [permission.id],
        },
      },
    });

    const response = await client.query({
      getProjectRoleById: {
        __scalar: true,
        __args: {
          id: created.createProjectRole.id,
        },
      },
    });
    expect(response.getProjectRoleById).toEqual(
      created.createProjectRole
    );
  });

  it('deletes an existing', async () => {
    const permission = await prismaClient.permission.create({
      data: {
        id: v4(),
        subject: nanoid(),
        action: 'CREATE',
      },
    });
    const title = nanoid();
    const created = await client.mutation({
      createProjectRole: {
        permissions: {
          __scalar: true,
        },
        __scalar: true,
        __args: {
          title,
          permissionIds: [permission.id],
        },
      },
    });

    const response = await client.mutation({
      deleteProjectRole: {
        __scalar: true,
        __args: {
          id: created.createProjectRole.id,
        },
      },
    });

    expect(response.deleteProjectRole.success).toBeTruthy();
  });

  it('updates an existing', async () => {
    const permission = await prismaClient.permission.create({
      data: {
        id: v4(),
        subject: nanoid(),
        action: 'CREATE',
      },
    });
    const title = nanoid();
    const created = await client.mutation({
      createProjectRole: {
        permissions: {
          __scalar: true,
        },
        __scalar: true,
        __args: {
          title,
          permissionIds: [permission.id],
        },
      },
    });
    const newTitle = nanoid();
    const response = await client.mutation({
      updateProjectRole: {
        __scalar: true,
        permissions: {
          __scalar: true,
        },
        __args: {
          id: created.createProjectRole.id,
          title: newTitle,
          permissionIds: [permission.id],
        },
      },
    });

    expect(response.updateProjectRole.permissions).toHaveLength(1);
    expect(response.updateProjectRole.permissions[0]).toEqual(
      expect.objectContaining(permission)
    );
    expect(response.updateProjectRole.id).toEqual(
      created.createProjectRole.id
    );
    expect(response.updateProjectRole.title).toEqual(newTitle);
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
    const response = await client.mutation({
      createProjectRole: {
        permissions: {
          __scalar: true,
        },
        __scalar: true,
        __args: {
          title,
          permissionIds: [permission.id],
        },
      },
    });

    expect(response.createProjectRole.permissions).toHaveLength(1);
    expect(response.createProjectRole.permissions[0]).toEqual(
      expect.objectContaining(permission)
    );
    expect(response.createProjectRole.id).toBeDefined();
    expect(response.createProjectRole.title).toEqual(title);
  });
});

// describe.skip('Project role', () => {
//   let client: AxiosInstance;

//   beforeAll(async () => {
//     client = await getAuthenticatedClient({ includeProjectId: true });
//   });

//   it('updates', async () => {
//     const permission = await prismaClient.permission.create({
//       data: {
//         id: v4(),
//         subject: nanoid(),
//         action: 'CREATE',
//       },
//     });
//     const title = nanoid();
//     const response = await client.post(`/projects/roles`, {
//       title,
//       permissionIds: [permission.id],
//     });
//     const responseData = response.data;
//     const newTitle = nanoid();
//     const newPermission = await prismaClient.permission.create({
//       data: {
//         id: v4(),
//         subject: nanoid(),
//         action: 'CREATE',
//       },
//     });

//     const updatedResponse = await client.patch(
//       `/projects/roles/${responseData.data.id}`,
//       {
//         title: newTitle,
//         permissionIds: [newPermission.id],
//       }
//     );
//     const updatedResponseData = updatedResponse.data;
//     expect(updatedResponseData.data.title).toEqual(newTitle);
//     expect(updatedResponseData.data.permissions).toHaveLength(1);
//     expect(newPermission).toEqual(
//       updatedResponseData.data.permissions[0]
//     );
//   });

//   it('gets by search', async () => {
//     const permission = await prismaClient.permission.create({
//       data: {
//         id: v4(),
//         subject: nanoid(),
//         action: 'CREATE',
//       },
//     });
//     const title = nanoid();
//     await client.post(`/projects/roles`, {
//       title,
//       permissionIds: [permission.id],
//     });

//     const getRoleResponse = await client.get(
//       `/projects/roles?search=${title}`
//     );
//     const responseData = getRoleResponse.data;

//     expect(responseData.data.length).toEqual(1);
//     expect(responseData.data.every((el: any) => el.title === title));
//     expect(
//       responseData.data.every(
//         (el: any) => el.projectId === PROJECT_ID
//       )
//     );
//   });

//   it('gets by project id', async () => {
//     const permission = await prismaClient.permission.create({
//       data: {
//         id: v4(),
//         subject: nanoid(),
//         action: 'CREATE',
//       },
//     });
//     const title = nanoid();
//     await client.post(`/projects/roles`, {
//       title,
//       permissionIds: [permission.id],
//     });

//     const getRoleResponse = await client.get(
//       `/projects/roles?search=Admin`
//     );
//     const responseData = getRoleResponse.data;

//     expect(responseData.data.length).toBeGreaterThanOrEqual(1);
//     expect(responseData.data.some((el: any) => el.title === title));
//     expect(
//       responseData.data.every(
//         (el: any) => el.projectId === PROJECT_ID
//       )
//     );
//   });

//   it('gets by id', async () => {
//     const permission = await prismaClient.permission.create({
//       data: {
//         id: v4(),
//         subject: nanoid(),
//         action: 'CREATE',
//       },
//     });
//     const title = nanoid();
//     const createdRoleResponse = await client.post(`/projects/roles`, {
//       title,
//       permissionIds: [permission.id],
//     });

//     const createdRoleResponseData = createdRoleResponse.data;

//     const getRoleResponse = await client.get(
//       `/projects/roles/${createdRoleResponseData.data.id}`
//     );
//     const responseData = getRoleResponse.data;

//     expect(responseData.data.title).toEqual(title);
//     expect(responseData.data.permissions).toHaveLength(1);
//     expect(responseData.data.permissions[0].id).toEqual(
//       permission.id
//     );
//   });

//   it('throws err when get by wrong id', async () => {
//     try {
//       const getRoleResponse = await client.get(
//         `/projects/roles/${v4()}`
//       );
//     } catch (error: any) {
//       expect(error.response.status).toEqual(404);
//       expect(error.response.data.message).toEqual('id not found');
//     }
//   });
// });

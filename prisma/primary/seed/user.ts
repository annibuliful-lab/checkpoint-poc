import { v4 } from 'uuid';
import { client } from './client';
import argon from 'argon2';

export async function seedUsers() {
  const userA = await client.account.upsert({
    include: {
      accountConfiguration: true,
      projectAccounts: {
        include: {
          role: {
            include: {
              project: true,
            },
          },
        },
      },
    },
    where: {
      username: 'userA1234',
    },
    update: {},
    create: {
      id: 'af466ea9-04ba-432e-98e0-3b8787dcda41',
      username: 'userA1234',
      password: await argon.hash('12345678'),
      createdBy: 'SEED',
      accountConfiguration: {
        create: {
          isActive: true,
        },
      },
      projectAccounts: {
        create: {
          id: v4(),
          createdBy: 'SEED',
          role: {
            create: {
              id: v4(),
              projectId: '246bb085-8ccc-4def-ac78-dc2ad5c7760b',
              title: 'Owner',
            },
          },
          project: {
            create: {
              id: '246bb085-8ccc-4def-ac78-dc2ad5c7760b',
              title: 'userA Project',
              createdBy: 'SEED',
              stationLocations: {
                create: {
                  id: 'e1c6783c-e09c-43dd-b1e5-8041dcd2816e',
                  title: 'SEED_STATION',
                  department: 'SEED_DEPARTMENT',
                  createdBy: 'SEED',
                  latitude: 0,
                  longitude: 0,
                },
              },
            },
          },
        },
      },
    },
  });

  const permissionIds = (
    await client.permission.findMany({
      select: {
        id: true,
      },
    })
  ).map((p) => p.id);

  for (const permissionId of permissionIds) {
    await client.projectRolePermission.upsert({
      update: {},
      create: {
        id: v4(),
        roleId: userA.projectAccounts[0].roleId,
        permissionId,
        projectId: userA.projectAccounts[0].projectId,
      },
      where: {
        roleId_permissionId: {
          roleId: userA.projectAccounts[0].roleId,
          permissionId,
        },
      },
    });
  }
}

import { nanoid } from 'nanoid';
import { getAuthenticatedClient, prismaClient } from './utils';
import { client } from '../../prisma/primary/seed/client';
import { v4 } from 'uuid';
import { PROJECT_ID } from './constants';

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

export async function createStationLocation() {
  return prismaClient.stationLocation.create({
    data: {
      id: v4(),
      title: nanoid(),
      department: nanoid(),
      latitude: 0,
      longitude: 0,
      projectId: PROJECT_ID,
      createdBy: 'SYSTEM',
    },
  });
}

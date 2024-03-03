import { PrismaClient } from '@prisma/client';
import {
  BACKEND_ENDPOINT,
  STATION_LOCATION_ID,
  httpClient,
} from './constants';
import axios from 'axios';
import { config } from 'dotenv';
import { createClient } from '../graphql/generated';
import { nanoid } from 'nanoid';
import { fetch } from 'undici';
config();

export const graphqlClient = createClient({
  url: process.env.BACKEND_ENDPOINT,
  fetch,
});

export const prismaClient = new PrismaClient();

type GetAuthenticatedClientParam = {
  includeProjectId?: boolean;
};
export async function getAuthenticatedClient({
  includeProjectId = false,
}: GetAuthenticatedClientParam) {
  const result = await graphqlClient.mutation({
    signin: {
      __scalar: true,
      __args: {
        username: 'userA1234',
        password: '12345678',
      },
    },
  });

  return createClient({
    url: process.env.BACKEND_ENDPOINT,
    headers: {
      Authorization: `Bearer ${result.signin.token}`,
      ...(includeProjectId && {
        'x-project-id': '246bb085-8ccc-4def-ac78-dc2ad5c7760b',
      }),
    },
  });
}

export async function getAuthenticatedClientWithRefreshToken() {
  const result = await httpClient.post('/auth/signin', {
    username: 'userA1234',
    password: '12345678',
  });

  return {
    client: axios.create({
      paramsSerializer: {
        indexes: null,
      },
      baseURL: BACKEND_ENDPOINT,
      headers: {
        Authorization: `Bearer ${result.data.data.token}`,
      },
    }),
    refreshToken: result.data.data.refreshToken as string,
  };
}

export async function createStationOfficer(
  stationLocationId: string
) {
  const firstname = nanoid();
  const lastname = nanoid();
  const msisdn = nanoid(10);
  const client = await getAuthenticatedClient({
    includeProjectId: true,
  });
  const stationOfficerResponse = await client.mutation({
    createStationOfficer: {
      __scalar: true,
      __args: {
        stationLocationId,
        firstname,
        lastname,
        msisdn,
      },
    },
  });

  return stationOfficerResponse.createStationOfficer;
}

export async function createStationLocation() {
  const title = nanoid();
  const department = nanoid();
  const latitude = Number(Math.random().toFixed(6));
  const longitude = Number(Math.random().toFixed(6));
  const tag = nanoid();
  const client = await getAuthenticatedClient({
    includeProjectId: true,
  });

  return (
    await client.mutation({
      createStationLocation: {
        __scalar: true,
        tags: {
          __scalar: true,
        },
        __args: {
          title,
          department,
          latitude,
          longitude,
          tags: [tag],
        },
      },
    })
  ).createStationLocation;
}
export async function createImeiConfiguration(
  stationLocationId: string
) {
  const tag = nanoid();
  const imei = nanoid();
  const client = await getAuthenticatedClient({
    includeProjectId: true,
  });
  const imeiResponse = await client.mutation({
    createImeiConfiguration: {
      __scalar: true,
      tags: {
        __scalar: true,
      },
      __args: {
        imei,
        stationLocationId,
        permittedLabel: 'NONE',
        blacklistPriority: 'NORMAL',
        tags: ['A', tag],
      },
    },
  });

  return imeiResponse.createImeiConfiguration;
}

export async function createImsiConfiguration() {
  const client = await getAuthenticatedClient({
    includeProjectId: true,
  });

  const imsi = nanoid();
  const created = await client.mutation({
    createImsiConfiguration: {
      __scalar: true,
      tags: {
        __scalar: true,
      },
      __args: {
        stationLocationId: STATION_LOCATION_ID,
        imsi,
        permittedLabel: 'WHITELIST',
        blacklistPriority: 'NORMAL',
        tags: ['A'],
      },
    },
  });

  return created.createImsiConfiguration;
}

export async function createMobileDeviceConfiguration() {
  const imei = await createImeiConfiguration(STATION_LOCATION_ID);
  const imsi = await createImsiConfiguration();
  const title = nanoid();
  const client = await getAuthenticatedClient({
    includeProjectId: true,
  });

  const mobileResponse = await client.mutation({
    createMobileDeviceConfiguration: {
      __scalar: true,
      referenceImeiConfiguration: {
        __scalar: true,
      },
      referenceImsiConfiguration: {
        __scalar: true,
      },
      __args: {
        stationLocationId: STATION_LOCATION_ID,
        referenceImeiConfigurationId: imei.id,
        referenceImsiConfigurationId: imsi.id,
        title,
        permittedLabel: 'NONE',
        blacklistPriority: 'NORMAL',
      },
    },
  });

  return mobileResponse.createMobileDeviceConfiguration;
}

export async function createStationDevice(stationLocationId: string) {
  const title = nanoid();
  const softwareVersion = nanoid(3);
  const hardwareVersion = nanoid(3);
  const client = await getAuthenticatedClient({
    includeProjectId: true,
  });
  const stationDeviceResponse = await client.mutation({
    createStationDevice: {
      __scalar: true,
      __args: {
        stationLocationId,
        title,
        softwareVersion,
        hardwareVersion,
      },
    },
  });

  return stationDeviceResponse.createStationDevice;
}

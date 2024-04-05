import { PrismaClient } from '@prisma/client';
import {
  BACKEND_ENDPOINT,
  STATION_LOCATION_ID,
  httpClient,
} from './constants';
import axios from 'axios';
import { config } from 'dotenv';
import { StationStatus, createClient } from '../graphql/generated';
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
  includeStationId?: boolean;
};

export async function getAuthenticatedClient({
  includeProjectId = false,
  includeStationId = false,
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
      ...(includeStationId && {
        'x-station-id': 'e1c6783c-e09c-43dd-b1e5-8041dcd2816e',
      }),
    },
  });
}

export function getStationAuthenticated() {
  return createClient({
    url: process.env.BACKEND_ENDPOINT,
    headers: {
      'x-api-key': 'V1StGXR8_Z5jdHi6B-myT',
      'x-device-id': '81bbfd00-f9f2-4145-b467-9423390f139d',
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
  const imei = generateIMEI();
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
    includeStationId: true,
  });

  const imsi = generateIMSI();
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
        imei: imei.imei,
        imsi: imsi.imsi,
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

export async function createStationLocationHealthCheckActivity(
  stationStatus: StationStatus,
  stationId: string
) {
  const client = await getAuthenticatedClient({
    includeProjectId: true,
  });

  const now = new Date();
  const createdActivityResponse = await client.mutation({
    createStationLocationHealthCheckActivity: {
      __scalar: true,
      __args: {
        stationId,
        stationStatus,
        startDatetime: now,
      },
    },
  });

  return createdActivityResponse.createStationLocationHealthCheckActivity;
}

export async function createVehicleTargetConfiguration() {
  const prefix = nanoid(2);
  const number = nanoid(6);
  const province = nanoid(6);
  const client = await getAuthenticatedClient({
    includeProjectId: true,
  });
  const vehicleTargetResponse = await client.mutation({
    createVehicleTargetConfiguration: {
      __scalar: true,
      __args: {
        stationLocationId: STATION_LOCATION_ID,
        color: 'BLACK',
        brand: 'TOYOTA',
        prefix,
        number,
        province,
        type: 'VIP',
        permittedLabel: 'WHITELIST',
        blacklistPriority: 'WARNING',
      },
    },
  });

  const vehicleTarget =
    vehicleTargetResponse.createVehicleTargetConfiguration;

  return vehicleTarget;
}

export function generateIMSI(): string {
  // Generate MCC (Mobile Country Code), 3 digits
  const mcc = String(Math.floor(Math.random() * 1000)).padStart(
    3,
    '0'
  );

  // Generate MNC (Mobile Network Code), 2 digits
  const mnc = String(Math.floor(Math.random() * 100)).padStart(
    2,
    '0'
  );

  // Generate MSIN (Mobile Subscription Identification Number), 10 digits
  const msin = String(
    Math.floor(Math.random() * 10000000000)
  ).padStart(10, '0');

  // Concatenate MCC, MNC, and MSIN to form IMSI
  const imsi = mcc + mnc + msin;

  return imsi;
}

// Generate a random IMEI (International Mobile Equipment Identity) number
export function generateIMEI(): string {
  // Function to seed the random number generator
  function seedRandom() {
    return Math.floor(Math.random() * 1000000);
  }

  // Seed the random number generator
  let seed = seedRandom();
  Math.random = function () {
    seed = (seed * 9301 + 49297) % 233280;
    return seed / 233280;
  };

  // Generate the first 14 digits randomly
  const imeiDigits: number[] = [];
  for (let i = 0; i < 14; i++) {
    imeiDigits.push(Math.floor(Math.random() * 10));
  }

  // Calculate the checksum
  let sum = 0;
  for (let i = 0; i < imeiDigits.length; i++) {
    let digit = imeiDigits[i];
    if (i % 2 === 0) {
      digit *= 2;
      if (digit >= 10) {
        digit = Math.floor(digit / 10) + (digit % 10);
      }
    }
    sum += digit;
  }
  const checksum = (10 - (sum % 10)) % 10;

  // Append the checksum to the IMEI digits
  imeiDigits.push(checksum);

  // Convert the digits to a string
  const imei = imeiDigits.join('');

  return imei;
}

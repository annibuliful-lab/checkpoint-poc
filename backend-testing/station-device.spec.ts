import { nanoid } from 'nanoid';
import { Client } from './graphql/generated';
import {
  createStationDevice,
  createStationLocation,
  getAuthenticatedClient,
} from './utils/utils';
import { v4 } from 'uuid';

describe('Station device', () => {
  let client: Client;
  let stationLocationId: string;

  beforeAll(async () => {
    client = await getAuthenticatedClient({ includeProjectId: true });
    stationLocationId = (await createStationLocation()).id;
  });

  it('gets by search', async () => {
    const createdStationDevice = await createStationDevice(
      stationLocationId
    );

    const stationDevicesResponse = await client.query({
      getStationDevices: {
        __scalar: true,
        __args: {
          stationLocationId,
          search: createdStationDevice.title,
          limit: 20,
          skip: 0,
        },
      },
    });

    const stationDevices = stationDevicesResponse.getStationDevices;

    expect(stationDevices.length).toEqual(1);
    expect(stationDevices[0]).toEqual(createdStationDevice);
  });

  it('throws error when get by wrong id', async () => {
    expect(
      client.query({
        getStationDeviceById: {
          __scalar: true,
          __args: {
            id: v4(),
          },
        },
      })
    ).rejects.toThrow();
  });

  it('gets by id', async () => {
    const createdStationDevice = await createStationDevice(
      stationLocationId
    );

    const stationDeviceResponse = await client.query({
      getStationDeviceById: {
        __scalar: true,
        __args: {
          id: createdStationDevice.id,
        },
      },
    });

    expect(stationDeviceResponse.getStationDeviceById).toEqual(
      createdStationDevice
    );
  });

  it('throws error when delete by wrong id', () => {
    expect(
      client.mutation({
        deleteStationDevice: {
          __scalar: true,
          __args: {
            id: v4(),
          },
        },
      })
    ).rejects.toThrow();
  });

  it('deletes an existing', async () => {
    const createdStationDevice = await createStationDevice(
      stationLocationId
    );

    const deleted = await client.mutation({
      deleteStationDevice: {
        __scalar: true,
        __args: {
          id: createdStationDevice.id,
        },
      },
    });

    expect(deleted.deleteStationDevice.success).toBeTruthy();
  });

  it('throws error when update by wrong id', () => {
    expect(
      client.mutation({
        updateStationDevice: {
          __scalar: true,
          __args: {
            id: v4(),
          },
        },
      })
    ).rejects.toThrow();
  });

  it('updates an existing', async () => {
    const createdStationDevice = await createStationDevice(
      stationLocationId
    );
    const title = nanoid();
    const softwareVersion = nanoid(3);
    const hardwareVersion = nanoid(3);
    const updated = await client.mutation({
      updateStationDevice: {
        __scalar: true,
        __args: {
          id: createdStationDevice.id,
          softwareVersion,
          hardwareVersion,
          title,
        },
      },
    });

    const stationDevice = updated.updateStationDevice;

    expect(stationDevice.title).toEqual(title);
    expect(stationDevice.hardwareVersion).toEqual(hardwareVersion);
    expect(stationDevice.softwareVersion).toEqual(softwareVersion);
  });

  it('creates', async () => {
    const stationLocation = await createStationLocation();
    const title = nanoid();
    const softwareVersion = nanoid(3);
    const hardwareVersion = nanoid(3);
    const stationDeviceResponse = await client.mutation({
      createStationDevice: {
        __scalar: true,
        __args: {
          stationLocationId: stationLocation.id,
          title,
          softwareVersion,
          hardwareVersion,
        },
      },
    });

    const stationDevice = stationDeviceResponse.createStationDevice;

    expect(stationDevice.title).toEqual(title);
    expect(stationDevice.hardwareVersion).toEqual(hardwareVersion);
    expect(stationDevice.softwareVersion).toEqual(softwareVersion);
  });
});

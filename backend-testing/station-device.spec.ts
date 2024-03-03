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

  it('gets activities by field resolver with status', async () => {
    const createdStationDevice = await createStationDevice(
      stationLocationId
    );
    const now = new Date();

    await Promise.all([
      client.mutation({
        createStationDeviceHealthCheckActivity: {
          __scalar: true,
          __args: {
            stationDeviceId: createdStationDevice.id,
            status: 'OFFLINE',
            activityTime: now,
          },
        },
      }),
      client.mutation({
        createStationDeviceHealthCheckActivity: {
          __scalar: true,
          __args: {
            stationDeviceId: createdStationDevice.id,
            status: 'ONLINE',
            activityTime: now,
          },
        },
      }),
    ]);

    const activitiesResponse = await client.query({
      getStationDeviceById: {
        __scalar: true,
        healthActivities: {
          __scalar: true,
          __args: {
            status: 'OFFLINE',
            limit: 20,
            skip: 0,
          },
        },
        __args: {
          id: createdStationDevice.id,
        },
      },
    });

    const activities =
      activitiesResponse.getStationDeviceById.healthActivities;

    expect(activities.length).toEqual(1);
    expect(activities[0].status).toEqual('OFFLINE');
  });

  it('gets activities by status', async () => {
    const createdStationDevice = await createStationDevice(
      stationLocationId
    );
    const now = new Date();

    await Promise.all([
      client.mutation({
        createStationDeviceHealthCheckActivity: {
          __scalar: true,
          __args: {
            stationDeviceId: createdStationDevice.id,
            status: 'OFFLINE',
            activityTime: now,
          },
        },
      }),
      client.mutation({
        createStationDeviceHealthCheckActivity: {
          __scalar: true,
          __args: {
            stationDeviceId: createdStationDevice.id,
            status: 'ONLINE',
            activityTime: now,
          },
        },
      }),
    ]);

    const activitiesResponse = await client.query({
      getStationDeviceHealthCheckActivities: {
        __scalar: true,
        __args: {
          stationDeviceId: createdStationDevice.id,
          limit: 20,
          skip: 0,
          status: 'OFFLINE',
        },
      },
    });

    const activities =
      activitiesResponse.getStationDeviceHealthCheckActivities;

    expect(activities.length).toEqual(1);
    expect(activities[0].status).toEqual('OFFLINE');
  });

  it('gets activities by dates', async () => {
    const createdStationDevice = await createStationDevice(
      stationLocationId
    );
    const now = new Date();

    await Promise.all([
      client.mutation({
        createStationDeviceHealthCheckActivity: {
          __scalar: true,
          __args: {
            stationDeviceId: createdStationDevice.id,
            status: 'OFFLINE',
            activityTime: now,
          },
        },
      }),
      client.mutation({
        createStationDeviceHealthCheckActivity: {
          __scalar: true,
          __args: {
            stationDeviceId: createdStationDevice.id,
            status: 'ONLINE',
            activityTime: now,
          },
        },
      }),
    ]);

    const activitiesResponse = await client.query({
      getStationDeviceHealthCheckActivities: {
        __scalar: true,
        __args: {
          stationDeviceId: createdStationDevice.id,
          limit: 20,
          skip: 0,
          startDate: now,
          endDate: now,
        },
      },
    });

    const activities =
      activitiesResponse.getStationDeviceHealthCheckActivities;

    expect(
      activities.every(
        (t) => new Date(t.activityTime).getTime() === now.getTime()
      )
    ).toBeTruthy();
  });

  it('creates new activity', async () => {
    const createdStationDevice = await createStationDevice(
      stationLocationId
    );
    const now = new Date();
    const activityResponse = await client.mutation({
      createStationDeviceHealthCheckActivity: {
        __scalar: true,
        __args: {
          stationDeviceId: createdStationDevice.id,
          status: 'OFFLINE',
          activityTime: now,
        },
      },
    });

    const activity =
      activityResponse.createStationDeviceHealthCheckActivity;
    expect(activity.stationDeviceId).toEqual(createdStationDevice.id);
    expect(activity.status).toEqual('OFFLINE');
    expect(new Date(activity.activityTime)).toEqual(now);
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

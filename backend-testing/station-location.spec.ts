import { nanoid } from 'nanoid';
import { Client } from './graphql/generated';
import {
  createStationLocation,
  createStationLocationHealthCheckActivity,
  createStationOfficer,
  getAuthenticatedClient,
} from './utils/utils';
import omit from 'lodash.omit';
import { v4 } from 'uuid';

describe('Station location', () => {
  let client: Client;

  beforeAll(async () => {
    client = await getAuthenticatedClient({ includeProjectId: true });
  });

  it('gets health check activities by dates', async () => {
    const createdStationLocation = await createStationLocation();

    const created = await Promise.all([
      createStationLocationHealthCheckActivity(
        'CLOSED',
        createdStationLocation.id
      ),
      createStationLocationHealthCheckActivity(
        'ONLINE',
        createdStationLocation.id
      ),
    ]);

    const result = await client.query({
      getStationLocationHealthCheckActivities: {
        __scalar: true,
        __args: {
          stationId: createdStationLocation.id,
          startDatetime: created[0].startDatetime,
          limit: 20,
          skip: 0,
        },
      },
    });

    expect(
      result.getStationLocationHealthCheckActivities.length
    ).toEqual(2);
  });

  it('gets health check activities by status', async () => {
    const createdStationLocation = await createStationLocation();

    await Promise.all([
      createStationLocationHealthCheckActivity(
        'CLOSED',
        createdStationLocation.id
      ),
      createStationLocationHealthCheckActivity(
        'ONLINE',
        createdStationLocation.id
      ),
    ]);

    const result = await client.query({
      getStationLocationHealthCheckActivities: {
        __scalar: true,
        __args: {
          stationId: createdStationLocation.id,
          stationStatus: 'ONLINE',
          limit: 20,
          skip: 0,
        },
      },
    });
    expect(
      result.getStationLocationHealthCheckActivities.length
    ).toEqual(1);
    expect(
      result.getStationLocationHealthCheckActivities[0].stationStatus
    ).toEqual('ONLINE');
  });

  it('updates an existing health check activity', async () => {
    const createdStationLocation = await createStationLocation();

    const createdActivity =
      await createStationLocationHealthCheckActivity(
        'CLOSED',
        createdStationLocation.id
      );

    const now = new Date();
    const updatedResponse = await client.mutation({
      updateStationLocationHealthCheckActivity: {
        __scalar: true,
        __args: {
          id: createdActivity.id,
          endDatetime: now,
          stationStatus: 'ONLINE',
        },
      },
    });

    const updatedActivity =
      updatedResponse.updateStationLocationHealthCheckActivity;

    expect(updatedActivity.id).toEqual(createdActivity.id);
    expect(updatedActivity.stationStatus).toEqual('ONLINE');
    expect(new Date(updatedActivity.endDatetime)).toEqual(now);
  });

  it('creates an health check activity', async () => {
    const createdStationLocation = await createStationLocation();
    const now = new Date();
    const createdActivityResponse = await client.mutation({
      createStationLocationHealthCheckActivity: {
        __scalar: true,
        __args: {
          stationId: createdStationLocation.id,
          stationStatus: 'MAINTENANCE',
          startDatetime: now,
        },
      },
    });

    const activity =
      createdActivityResponse.createStationLocationHealthCheckActivity;

    expect(activity.stationId).toEqual(createdStationLocation.id);
    expect(activity.stationStatus).toEqual('MAINTENANCE');
    expect(activity.endDatetime).toBeNull();
    expect(new Date(activity.startDatetime)).toEqual(now);
  });

  it('gets by tag', async () => {
    const createdStationLocation = await createStationLocation();
    const stationLocations = await client.query({
      getStationLocations: {
        __scalar: true,
        __args: {
          limit: 20,
          skip: 0,
          tags: createdStationLocation.tags?.map((t) => t.title),
        },
      },
    });

    expect(stationLocations.getStationLocations.length).toEqual(1);
    expect(stationLocations.getStationLocations[0]).toEqual(
      omit(createdStationLocation, 'tags')
    );
  });

  it('gets by search', async () => {
    const createdStationLocation = await createStationLocation();
    const stationLocations = await client.query({
      getStationLocations: {
        __scalar: true,
        __args: {
          search: createdStationLocation.title,
          limit: 20,
          skip: 0,
        },
      },
    });
    expect(stationLocations.getStationLocations.length).toEqual(1);
    expect(stationLocations.getStationLocations[0]).toEqual(
      omit(createdStationLocation, 'tags')
    );
  });

  it('throws error when delete by wrong id', () => {
    expect(
      client.mutation({
        deleteStationLocation: {
          __scalar: true,
          __args: {
            id: v4(),
          },
        },
      })
    ).rejects.toThrow();
  });

  it('deletes an existing', async () => {
    const createdStationLocation = await createStationLocation();
    const result = await client.mutation({
      deleteStationLocation: {
        __scalar: true,
        __args: {
          id: createdStationLocation.id,
        },
      },
    });

    expect(result.deleteStationLocation.success).toBeTruthy();
    expect(
      client.query({
        getStationLocationById: {
          __scalar: true,
          tags: {
            __scalar: true,
          },
          __args: {
            id: createdStationLocation.id,
          },
        },
      })
    ).rejects.toThrow();
  });

  it('throws error when get by wrong id', () => {
    expect(
      client.query({
        getStationLocationById: {
          __scalar: true,
          tags: {
            __scalar: true,
          },
          __args: {
            id: v4(),
          },
        },
      })
    ).rejects.toThrow();
  });

  it('gets an existing by id', async () => {
    const createdStationLocation = await createStationLocation();
    const stationOfficer = await createStationOfficer(
      createdStationLocation.id
    );
    const result = await client.query({
      getStationLocationById: {
        officers: {
          __scalar: true,
        },
        __scalar: true,
        tags: {
          __scalar: true,
        },
        __args: {
          id: createdStationLocation.id,
        },
      },
    });

    expect(result.getStationLocationById.officers?.length).toEqual(1);
    expect(result.getStationLocationById.officers?.[0]).toEqual(
      stationOfficer
    );
    expect(omit(result.getStationLocationById, 'officers')).toEqual(
      createdStationLocation
    );
  });

  it('throws error when update wrong id', async () => {
    const newTitle = nanoid();
    const tag = nanoid();
    expect(
      client.mutation({
        updateStationLocation: {
          __scalar: true,
          tags: {
            __scalar: true,
          },
          __args: {
            id: v4(),
            title: newTitle,
            tags: [tag],
          },
        },
      })
    ).rejects.toThrow();
  });

  it('updates an existing', async () => {
    const createdStationLocation = await createStationLocation();
    const newTitle = nanoid();
    const tag = nanoid();
    const updated = await client.mutation({
      updateStationLocation: {
        __scalar: true,
        tags: {
          __scalar: true,
        },
        __args: {
          id: createdStationLocation.id,
          title: newTitle,
          tags: [tag],
        },
      },
    });

    expect(updated.updateStationLocation.tags?.length).toEqual(1);
    expect(updated.updateStationLocation.tags?.[0].title).toEqual(
      tag
    );
    expect(omit(updated.updateStationLocation, 'tags')).toEqual(
      omit(
        {
          ...createdStationLocation,
          title: newTitle,
        },
        'tags'
      )
    );
  });

  it('creates', async () => {
    const title = nanoid();
    const department = nanoid();
    const latitude = Number(Math.random().toFixed(6));
    const longitude = Number(Math.random().toFixed(6));
    const tag = nanoid();
    const stationLocation = await client.mutation({
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
    });

    expect(
      stationLocation.createStationLocation.tags?.length
    ).toEqual(1);
    expect(
      stationLocation.createStationLocation.tags?.[0].title
    ).toEqual(tag);
    expect(stationLocation.createStationLocation.latitude).toEqual(
      latitude
    );
    expect(stationLocation.createStationLocation.longitude).toEqual(
      longitude
    );
    expect(stationLocation.createStationLocation.department).toEqual(
      department
    );
    expect(stationLocation.createStationLocation.title).toEqual(
      title
    );
  });
});

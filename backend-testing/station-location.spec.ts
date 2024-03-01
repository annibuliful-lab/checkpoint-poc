import { nanoid } from 'nanoid';
import { Client } from './graphql/generated';
import {
  createStationLocation,
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

import { nanoid } from 'nanoid';
import { Client } from './graphql/generated';
import {
  createImeiConfiguration,
  getAuthenticatedClient,
} from './utils/utils';
import { createStationLocation } from './utils/project';
import { STATION_LOCATION_ID } from './utils/constants';
import omit from 'lodash.omit';

describe('Imei configuration', () => {
  let client: Client;

  beforeAll(async () => {
    client = await getAuthenticatedClient({ includeProjectId: true });
  });

  it('throws error when create new without provide project id', async () => {
    const client = await getAuthenticatedClient({});
    const tag = nanoid();
    const imei = nanoid();
    try {
      await client.mutation({
        createImeiConfiguration: {
          __scalar: true,
          tags: {
            __scalar: true,
          },
          __args: {
            imei,
            stationLocationId: STATION_LOCATION_ID,
            permittedLabel: 'NONE',
            blacklistPriority: 'NORMAL',
            tags: ['A', tag],
          },
        },
      });
    } catch (error: any) {
      expect(error.errors[0].extensions.message).toEqual(
        'Project id is required'
      );
    }
  });
  it('gets by tags', async () => {
    const createdImeiConfiguration = await createImeiConfiguration(
      STATION_LOCATION_ID
    );
    const tags = createdImeiConfiguration.tags?.map((t) => t.title);
    const imeiConfigurations = await client.query({
      getImeiConfigurations: {
        __scalar: true,
        tags: {
          __scalar: true,
        },
        __args: {
          stationLocationId:
            createdImeiConfiguration.stationLocationId,
          limit: 20,
          skip: 0,
          tags: ['A'],
        },
      },
    });

    expect(
      imeiConfigurations.getImeiConfigurations.every((t) =>
        t.tags?.some((el) => el.title === 'A')
      )
    ).toBeTruthy();
  });

  it('gets by id', async () => {
    const createdImeiConfiguration = await createImeiConfiguration(
      STATION_LOCATION_ID
    );
    const imeiConfiguration = await client.query({
      getImeiConfigurationById: {
        __scalar: true,
        tags: {
          __scalar: true,
        },
        __args: {
          id: createdImeiConfiguration.id,
        },
      },
    });

    expect(
      omit(imeiConfiguration.getImeiConfigurationById, 'tags')
    ).toEqual(omit(createdImeiConfiguration, 'tags'));
  });

  it('deletes an existing', async () => {
    const createdImeiConfiguration = await createImeiConfiguration(
      STATION_LOCATION_ID
    );
    const deleted = await client.mutation({
      deleteImeiConfiguration: {
        __scalar: true,
        __args: {
          id: createdImeiConfiguration.id,
        },
      },
    });
    expect(deleted.deleteImeiConfiguration.success).toBeTruthy();
  });

  it('updates an existing', async () => {
    const createdImeiConfiguration = await createImeiConfiguration(
      STATION_LOCATION_ID
    );
    const newImei = nanoid();
    const newTag = nanoid();
    const updated = await client.mutation({
      updateImeiConfiguration: {
        __scalar: true,
        tags: {
          __scalar: true,
        },
        __args: {
          id: createdImeiConfiguration.id,
          imei: newImei,
          permittedLabel: 'BLACKLIST',
          blacklistPriority: 'WARNING',
          tags: [newTag],
        },
      },
    });

    const imeiConfiguration = updated.updateImeiConfiguration;
    expect(imeiConfiguration.imei).toEqual(newImei);
    expect(imeiConfiguration.blacklistPriority).toEqual('WARNING');
    expect(imeiConfiguration.permittedLabel).toEqual('BLACKLIST');
    expect(imeiConfiguration.tags?.length).toEqual(1);
    expect(imeiConfiguration.tags?.[0].title).toEqual(newTag);
  });

  it('creates', async () => {
    const stationLocation = await createStationLocation();
    const tag = nanoid();
    const imei = nanoid();
    const imeiResponse = await client.mutation({
      createImeiConfiguration: {
        __scalar: true,
        tags: {
          __scalar: true,
        },
        __args: {
          imei,
          stationLocationId: stationLocation.id,
          permittedLabel: 'NONE',
          blacklistPriority: 'NORMAL',
          tags: ['A', tag],
        },
      },
    });

    const imeiConfiguration = imeiResponse.createImeiConfiguration;
    expect(imeiConfiguration.imei).toEqual(imei);
    expect(imeiConfiguration.blacklistPriority).toEqual('NORMAL');
    expect(imeiConfiguration.permittedLabel).toEqual('NONE');
    expect(imeiConfiguration.tags?.length).toEqual(2);

    expect(
      ['A', tag].includes(imeiConfiguration.tags?.[0].title as string)
    ).toBeTruthy();
    expect(
      ['A', tag].includes(imeiConfiguration.tags?.[1].title as string)
    ).toBeTruthy();
  });
});

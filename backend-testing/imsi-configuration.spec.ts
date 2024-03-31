import { nanoid } from 'nanoid';
import { Client } from './graphql/generated';
import { createStationLocation } from './utils/project';
import {
  createImsiConfiguration,
  getAuthenticatedClient,
} from './utils/utils';
import { v4 } from 'uuid';
import omit from 'lodash.omit';
import { STATION_LOCATION_ID } from './utils/constants';

describe('Imsi configuration', () => {
  let client: Client;

  beforeAll(async () => {
    client = await getAuthenticatedClient({
      includeProjectId: true,
      includeStationId: true,
    });
  });

  it('upsert imsi', async () => {
    const { id, imsi } = await createImsiConfiguration();
    const upsertResponse = await client.mutation({
      upsertImsiConfiguration: {
        __scalar: true,
        __args: {
          imsi,
          permittedLabel: 'NONE',
          blacklistPriority: 'NORMAL',
        },
      },
    });

    const upsert = upsertResponse.upsertImsiConfiguration;

    expect(upsert.imsi).toEqual(imsi);
    expect(upsert.id).toEqual(id);
  });

  it('gets by tags', async () => {
    const { stationLocationId } = await createImsiConfiguration();
    const imsiConfigurationsResponse = await client.query({
      getImsiConfigurations: {
        __scalar: true,
        tags: {
          title: true,
        },
        __args: {
          stationLocationId: stationLocationId,
          limit: 20,
          skip: 0,
          tags: ['A'],
        },
      },
    });

    const imsiConfigurations =
      imsiConfigurationsResponse.getImsiConfigurations;

    expect(imsiConfigurations.length).toBeGreaterThanOrEqual(1);
    expect(
      imsiConfigurations.every((imsi) =>
        imsi.tags?.some((tag) => tag.title === 'A')
      )
    ).toBeTruthy();
  });

  it('deletes an existing', async () => {
    const created = await createImsiConfiguration();
    const deleted = await client.mutation({
      deleteImsiConfiguration: {
        __scalar: true,
        __args: {
          id: created.id,
        },
      },
    });

    expect(deleted.deleteImsiConfiguration.success).toBeTruthy();
  });

  it('throws error when invalid id', async () => {
    expect(
      client.query({
        getImsiConfigurationById: {
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

  it('gets by id', async () => {
    const created = await createImsiConfiguration();
    const imsiConfiguration = await client.query({
      getImsiConfigurationById: {
        __scalar: true,
        tags: {
          __scalar: true,
        },
        __args: {
          id: created.id,
        },
      },
    });

    expect(omit(created, 'tags')).toEqual(
      omit(imsiConfiguration.getImsiConfigurationById, 'tags')
    );
  });

  it('updates an existing', async () => {
    const created = await createImsiConfiguration();
    const imsi = created.imsi;
    const tag = nanoid();
    const updated = await client.mutation({
      updateImsiConfiguration: {
        __scalar: true,
        tags: {
          __scalar: true,
        },
        __args: {
          id: created.id,
          imsi: created.imsi,
          blacklistPriority: 'WARNING',
          permittedLabel: 'BLACKLIST',
          tags: [tag, 'A'],
        },
      },
    });
    const imsiConfiguration = updated.updateImsiConfiguration;
    expect(imsiConfiguration.tags?.length).toEqual(2);
    expect(
      [tag, 'A'].includes(imsiConfiguration.tags?.[0].title as string)
    ).toBeTruthy();
    expect(
      [tag, 'A'].includes(imsiConfiguration.tags?.[1].title as string)
    ).toBeTruthy();
    expect(imsiConfiguration.imsi).toEqual(imsi);
    expect(imsiConfiguration.mcc).toEqual(
      imsi[0] + imsi[1] + imsi[2]
    );
    expect(imsiConfiguration.mnc).toEqual(
      imsi[3] + imsi[4] + imsi[5]
    );
    expect(imsiConfiguration.stationLocationId).toEqual(
      created.stationLocationId
    );
    expect(imsiConfiguration.permittedLabel).toEqual('BLACKLIST');
    expect(imsiConfiguration.blacklistPriority).toEqual('WARNING');
  });

  it('throws error when create without providing project id', async () => {
    const imsi = nanoid();
    const client = await getAuthenticatedClient({});

    try {
      await client.mutation({
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
    } catch (error: any) {
      expect(error.errors[0].extensions.message).toEqual(
        'Project id is required'
      );
    }
  });
  it('creates', async () => {
    const stationLocation = await createStationLocation();
    const imsi = nanoid();
    const created = await client.mutation({
      createImsiConfiguration: {
        __scalar: true,
        tags: {
          __scalar: true,
        },
        __args: {
          stationLocationId: stationLocation.id,
          imsi,
          permittedLabel: 'WHITELIST',
          blacklistPriority: 'NORMAL',
          tags: ['A'],
        },
      },
    });
    const imsiConfiguration = created.createImsiConfiguration;

    expect(imsiConfiguration.tags?.length).toEqual(1);
    expect(imsiConfiguration.tags?.[0].title).toEqual('A');
    expect(imsiConfiguration.imsi).toEqual(imsi);
    expect(imsiConfiguration.mcc).toEqual(
      imsi[0] + imsi[1] + imsi[2]
    );
    expect(imsiConfiguration.mnc).toEqual(
      imsi[3] + imsi[4] + imsi[5]
    );
    expect(imsiConfiguration.stationLocationId).toEqual(
      stationLocation.id
    );
    expect(imsiConfiguration.permittedLabel).toEqual('WHITELIST');
    expect(imsiConfiguration.blacklistPriority).toEqual('NORMAL');
  });
});

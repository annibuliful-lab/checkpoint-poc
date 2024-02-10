import { nanoid } from 'nanoid';
import { Client } from './graphql/generated';
import {
  createImeiConfiguration,
  getAuthenticatedClient,
} from './utils/utils';
import { createStationLocation } from './utils/project';
import { title } from 'process';

describe('Imei configuration', () => {
  let client: Client;

  beforeAll(async () => {
    client = await getAuthenticatedClient({ includeProjectId: true });
  });

  it('gets by tags', async () => {
    const createdImeiConfiguration = await createImeiConfiguration();
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
          tags,
        },
      },
    });

    expect(
      imeiConfigurations.getImeiConfigurations.every((t) =>
        t.tags?.every((el) => tags?.includes(el.title))
      )
    ).toBeTruthy();
  });
  it('gets by id', async () => {
    const createdImeiConfiguration = await createImeiConfiguration();
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

    expect(imeiConfiguration.getImeiConfigurationById).toEqual(
      createdImeiConfiguration
    );
  });

  it('deletes an existing', async () => {
    const createdImeiConfiguration = await createImeiConfiguration();
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
    const createdImeiConfiguration = await createImeiConfiguration();
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
          priority: 'WARNING',
          tags: [newTag],
        },
      },
    });

    const imeiConfiguration = updated.updateImeiConfiguration;
    expect(imeiConfiguration.imei).toEqual(newImei);
    expect(imeiConfiguration.priority).toEqual('WARNING');
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
          priority: 'NORMAL',
          tags: ['A', tag],
        },
      },
    });

    const imeiConfiguration = imeiResponse.createImeiConfiguration;
    expect(imeiConfiguration.imei).toEqual(imei);
    expect(imeiConfiguration.priority).toEqual('NORMAL');
    expect(imeiConfiguration.permittedLabel).toEqual('NONE');
    expect(imeiConfiguration.tags?.length).toEqual(2);
    expect(imeiConfiguration.tags?.[0].title).toEqual('A');
    expect(imeiConfiguration.tags?.[1].title).toEqual(tag);
  });
});

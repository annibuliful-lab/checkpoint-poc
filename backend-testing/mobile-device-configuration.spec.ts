import { nanoid } from 'nanoid';
import { Client } from './graphql/generated';
import {
  createImeiConfiguration,
  createImsiConfiguration,
  createMobileDeviceConfiguration,
  getAuthenticatedClient,
} from './utils/utils';
import { v4 } from 'uuid';
import omit from 'lodash.omit';
import { STATION_LOCATION_ID } from './utils/constants';

describe('Mobile device configuration', () => {
  let client: Client;

  beforeAll(async () => {
    client = await getAuthenticatedClient({ includeProjectId: true });
  });

  it('gets by tags', async () => {
    await createMobileDeviceConfiguration();
    const response = await client.query({
      getMobileDeviceConfigurations: {
        __scalar: true,
        tags: {
          __scalar: true,
        },
        __args: {
          stationLocationId: STATION_LOCATION_ID,
          tags: ['A'],
          limit: 20,
          skip: 0,
        },
      },
    });

    const mobileDevices = response.getMobileDeviceConfigurations;
    expect(
      mobileDevices.every(
        (m) => m.stationLocationId === STATION_LOCATION_ID
      )
    ).toBeTruthy();
    expect(
      mobileDevices.every((m) => m.tags?.some((t) => t.title === 'A'))
    ).toBeTruthy();
  });

  it('throws error when update with wrong imsi id', async () => {
    const mobileDevice = await createMobileDeviceConfiguration();
    const newTags = nanoid();
    const msisdn = nanoid();
    expect(
      client.mutation({
        updateMobileDeviceConfiguration: {
          __scalar: true,
          __args: {
            id: mobileDevice.id,
            referenceImsiConfigurationId: v4(),
            permittedLabel: 'BLACKLIST',
            blacklistPriority: 'DANGER',
            msisdn,
            tags: [newTags],
          },
        },
      })
    ).rejects.toThrow();
  });

  it('throws error when update with wrong imei id', async () => {
    const mobileDevice = await createMobileDeviceConfiguration();
    const newTags = nanoid();
    const msisdn = nanoid();
    expect(
      client.mutation({
        updateMobileDeviceConfiguration: {
          __scalar: true,
          __args: {
            id: mobileDevice.id,
            referenceImeiConfigurationId: v4(),
            permittedLabel: 'BLACKLIST',
            blacklistPriority: 'DANGER',
            msisdn,
            tags: [newTags],
          },
        },
      })
    ).rejects.toThrow();
  });

  it('updates an existing', async () => {
    const mobileDevice = await createMobileDeviceConfiguration();
    const newTags = nanoid();
    const msisdn = nanoid();
    const updated = await client.mutation({
      updateMobileDeviceConfiguration: {
        __scalar: true,
        __args: {
          id: mobileDevice.id,
          permittedLabel: 'BLACKLIST',
          blacklistPriority: 'DANGER',
          msisdn,
          tags: [newTags],
        },
      },
    });

    const updatedMobileDevice =
      updated.updateMobileDeviceConfiguration;

    expect(updatedMobileDevice.blacklistPriority).toEqual('DANGER');
    expect(updatedMobileDevice.id).toEqual(mobileDevice.id);
    expect(updatedMobileDevice.permittedLabel).toEqual('BLACKLIST');
    expect(updatedMobileDevice.msisdn).toEqual(msisdn);
  });

  it('throws error when delete by wrong id', () => {
    expect(
      client.mutation({
        deleteMobileDeviceConfiguration: {
          __scalar: true,
          __args: {
            id: v4(),
          },
        },
      })
    ).rejects.toThrow();
  });

  it('deletes an existing', async () => {
    const mobileDevice = await createMobileDeviceConfiguration();
    const deleted = await client.mutation({
      deleteMobileDeviceConfiguration: {
        __scalar: true,
        __args: {
          id: mobileDevice.id,
        },
      },
    });

    expect(
      deleted.deleteMobileDeviceConfiguration.success
    ).toBeTruthy();
  });

  it('throws error when create with invalid imsi id', async () => {
    const imei = await createImeiConfiguration(STATION_LOCATION_ID);
    expect(
      client.mutation({
        createMobileDeviceConfiguration: {
          __scalar: true,
          __args: {
            stationLocationId: STATION_LOCATION_ID,
            referenceImeiConfigurationId: imei.id,
            referenceImsiConfigurationId: v4(),
            title: nanoid(),
            permittedLabel: 'NONE',
            blacklistPriority: 'NORMAL',
          },
        },
      })
    ).rejects.toThrow();
  });

  it('throws error when create with invalid imei id', async () => {
    const imsi = await createImsiConfiguration();
    expect(
      client.mutation({
        createMobileDeviceConfiguration: {
          __scalar: true,
          __args: {
            stationLocationId: STATION_LOCATION_ID,
            referenceImeiConfigurationId: v4(),
            referenceImsiConfigurationId: imsi.id,
            title: nanoid(),
            permittedLabel: 'NONE',
            blacklistPriority: 'NORMAL',
          },
        },
      })
    ).rejects.toThrow();
  });

  it('creates with msisdn', async () => {
    const imei = await createImeiConfiguration(STATION_LOCATION_ID);
    const imsi = await createImsiConfiguration();
    const title = nanoid();
    const msisdn = nanoid();
    const mobileResponse = await client.mutation({
      createMobileDeviceConfiguration: {
        __scalar: true,
        __args: {
          stationLocationId: STATION_LOCATION_ID,
          referenceImeiConfigurationId: imei.id,
          referenceImsiConfigurationId: imsi.id,
          title,
          permittedLabel: 'NONE',
          blacklistPriority: 'NORMAL',
          msisdn,
        },
      },
    });

    const mobileDevice =
      mobileResponse.createMobileDeviceConfiguration;
    expect(mobileDevice.title).toEqual(title);
    expect(mobileDevice.msisdn).toEqual(msisdn);
    expect(mobileDevice.permittedLabel).toEqual('NONE');
    expect(mobileDevice.blacklistPriority).toEqual('NORMAL');
  });

  it('creates', async () => {
    const imei = await createImeiConfiguration(STATION_LOCATION_ID);
    const imsi = await createImsiConfiguration();
    const title = nanoid();
    const tag = nanoid();
    const mobileResponse = await client.mutation({
      createMobileDeviceConfiguration: {
        __scalar: true,
        tags: {
          __scalar: true,
        },
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
          tags: [tag],
        },
      },
    });

    const mobileDevice =
      mobileResponse.createMobileDeviceConfiguration;

    expect(mobileDevice.referenceImeiConfiguration).toEqual(
      omit(imei, 'tags')
    );
    expect(mobileDevice.referenceImsiConfiguration).toEqual(
      omit(imsi, 'tags')
    );
    expect(mobileDevice.tags?.[0].title).toEqual(tag);
    expect(mobileDevice.title).toEqual(title);
    expect(mobileDevice.permittedLabel).toEqual('NONE');
    expect(mobileDevice.blacklistPriority).toEqual('NORMAL');
  });
});

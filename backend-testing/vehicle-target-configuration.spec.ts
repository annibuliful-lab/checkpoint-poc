import { nanoid } from 'nanoid';
import { Client } from './graphql/generated';
import { STATION_LOCATION_ID } from './utils/constants';
import {
  createVehicleTargetConfiguration,
  getAuthenticatedClient,
  prismaClient,
} from './utils/utils';

describe('Vehicle target configuration', () => {
  let client: Client;

  beforeAll(async () => {
    client = await getAuthenticatedClient({ includeProjectId: true });
  });

  it('creates with tags', async () => {
    const tags = [nanoid()];
    const prefix = nanoid(2);
    const number = nanoid(6);
    const province = nanoid(6);
    const vehicleTargetResponse = await client.mutation({
      createVehicleTargetConfiguration: {
        __scalar: true,
        tags: {
          __scalar: true,
        },
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
          tags,
        },
      },
    });

    const vehicleTarget =
      vehicleTargetResponse.createVehicleTargetConfiguration;

    const vehicleTargetTags = vehicleTarget.tags;

    expect(vehicleTargetTags?.length).toEqual(1);
    expect(vehicleTargetTags?.[0].title).toEqual(tags[0]);

    expect(vehicleTarget.number).toEqual(number);
    expect(vehicleTarget.color).toEqual('BLACK');
    expect(vehicleTarget.prefix).toEqual(prefix);
    expect(vehicleTarget.province).toEqual(province);
    expect(vehicleTarget.type).toEqual('VIP');
  });

  it('deletes an image existing', async () => {
    const prefix = nanoid(2);
    const number = nanoid(6);
    const province = nanoid(6);
    const image = nanoid();
    const vehicleTargetResponse = await client.mutation({
      createVehicleTargetConfiguration: {
        __scalar: true,
        images: {
          __scalar: true,
        },
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
          imageS3Keys: [
            {
              s3Key: image,
              type: 'CONFIG',
            },
          ],
        },
      },
    });

    const vehicleTarget =
      vehicleTargetResponse.createVehicleTargetConfiguration;

    const deleteVehicleTargetImage = await client.mutation({
      deleteVehicleTargetConfigurationImage: {
        __scalar: true,
        __args: {
          id: vehicleTarget.images?.[0].id as string,
        },
      },
    });

    expect(
      deleteVehicleTargetImage.deleteVehicleTargetConfigurationImage
        .success
    ).toBeTruthy();

    const getVehicleTargetResponse = await client.query({
      getVehicleTargetConfigurationById: {
        images: {
          __scalar: true,
        },
        __args: {
          id: vehicleTarget.id,
        },
      },
    });

    expect(
      getVehicleTargetResponse.getVehicleTargetConfigurationById
        .images?.length
    ).toEqual(0);
  });

  it('creates with images', async () => {
    const prefix = nanoid(2);
    const number = nanoid(6);
    const province = nanoid(6);
    const image = nanoid();
    const vehicleTargetResponse = await client.mutation({
      createVehicleTargetConfiguration: {
        __scalar: true,
        images: {
          __scalar: true,
        },
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
          imageS3Keys: [
            {
              s3Key: image,
              type: 'CONFIG',
            },
          ],
        },
      },
    });

    const vehicleTarget =
      vehicleTargetResponse.createVehicleTargetConfiguration;

    expect(vehicleTarget.images?.length).toEqual(1);
    expect(vehicleTarget.images?.[0].s3Key).toEqual(image);
    expect(vehicleTarget.images?.[0].type).toEqual('CONFIG');
  });

  it('creates', async () => {
    const prefix = nanoid(2);
    const number = nanoid(6);
    const province = nanoid(6);
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

    expect(vehicleTarget.number).toEqual(number);
    expect(vehicleTarget.color).toEqual('BLACK');
    expect(vehicleTarget.prefix).toEqual(prefix);
    expect(vehicleTarget.province).toEqual(province);
    expect(vehicleTarget.type).toEqual('VIP');
  });

  it('updates with tags', async () => {
    const vehicleTarget = await createVehicleTargetConfiguration();
    const prefix = nanoid(2);
    const number = nanoid(6);
    const province = nanoid(6);
    const type = nanoid(3);
    const tag = nanoid();
    const updatedResponse = await client.mutation({
      updateVehicleTargetConfiguration: {
        __scalar: true,
        tags: { __scalar: true },
        __args: {
          id: vehicleTarget.id,
          prefix,
          number,
          province,
          type,
          color: 'WHITE',
          tags: [tag],
        },
      },
    });

    const updated = updatedResponse.updateVehicleTargetConfiguration;
    const vehicleTargetTags = updated.tags;

    expect(vehicleTargetTags?.length).toEqual(1);
    expect(vehicleTargetTags?.[0].title).toEqual(tag);

    expect(updated.number).toEqual(number);
    expect(updated.color).toEqual('WHITE');
    expect(updated.prefix).toEqual(prefix);
    expect(updated.province).toEqual(province);
    expect(updated.type).toEqual(type);
  });

  it('updates an existing', async () => {
    const vehicleTarget = await createVehicleTargetConfiguration();
    const prefix = nanoid(2);
    const number = nanoid(6);
    const province = nanoid(6);
    const type = nanoid(3);
    const updatedResponse = await client.mutation({
      updateVehicleTargetConfiguration: {
        __scalar: true,
        __args: {
          id: vehicleTarget.id,
          prefix,
          number,
          province,
          type,
          color: 'WHITE',
        },
      },
    });

    const updated = updatedResponse.updateVehicleTargetConfiguration;
    expect(updated.number).toEqual(number);
    expect(updated.color).toEqual('WHITE');
    expect(updated.prefix).toEqual(prefix);
    expect(updated.province).toEqual(province);
    expect(updated.type).toEqual(type);
  });

  it('deletes an existing', async () => {
    const vehicleTarget = await createVehicleTargetConfiguration();
    const deletedResponse = await client.mutation({
      deleteVehicleTargetConfiguration: {
        __scalar: true,
        __args: {
          id: vehicleTarget.id,
        },
      },
    });

    const recheck =
      await prismaClient.vehicleTargetConfiguration.findUnique({
        where: {
          id: vehicleTarget.id,
        },
      });

    expect(recheck?.deletedAt).not.toBeNull();
    expect(recheck?.deletedBy).not.toBeNull();
    expect(
      deletedResponse.deleteVehicleTargetConfiguration.success
    ).toBeTruthy();
  });

  it('gets by id', async () => {
    const vehicleTarget = await createVehicleTargetConfiguration();
    const getResponse = await client.query({
      getVehicleTargetConfigurationById: {
        __scalar: true,
        __args: {
          id: vehicleTarget.id,
        },
      },
    });

    expect(getResponse.getVehicleTargetConfigurationById).toEqual(
      vehicleTarget
    );
  });
});

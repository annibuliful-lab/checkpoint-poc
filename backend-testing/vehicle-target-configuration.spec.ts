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

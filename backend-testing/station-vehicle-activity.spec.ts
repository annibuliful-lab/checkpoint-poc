import { nanoid } from 'nanoid';
import { Client } from './graphql/generated';
import { getStationAuthenticated } from './utils/utils';

describe('station-vehicle-activity', () => {
  let client: Client;
  beforeAll(() => {
    client = getStationAuthenticated();
  });

  it('creates with status', async () => {
    const model = nanoid();
    const brand = nanoid();
    const color = nanoid();
    const vehicleActivityResponse = await client.mutation({
      createStationVehicleActivity: {
        __scalar: true,
        __args: {
          model,
          brand,
          color,
          status: 'INVESTIGATING',
        },
      },
    });

    const vehicleActivity =
      vehicleActivityResponse.createStationVehicleActivity;
    expect(vehicleActivity.brand).toEqual(brand);
    expect(vehicleActivity.color).toEqual(color);
    expect(vehicleActivity.model).toEqual(model);
    expect(vehicleActivity.status).toEqual('INVESTIGATING');
  });

  it('creates', async () => {
    const model = nanoid();
    const brand = nanoid();
    const color = nanoid();
    const vehicleActivityResponse = await client.mutation({
      createStationVehicleActivity: {
        __scalar: true,
        __args: {
          model,
          brand,
          color,
        },
      },
    });
    const vehicleActivity =
      vehicleActivityResponse.createStationVehicleActivity;

    expect(vehicleActivity.brand).toEqual(brand);
    expect(vehicleActivity.color).toEqual(color);
    expect(vehicleActivity.model).toEqual(model);
    expect(vehicleActivity.status).toEqual('IN_QUEUE');
  });
});

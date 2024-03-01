import { nanoid } from 'nanoid';
import { Client } from './graphql/generated';
import {
  createStationLocation,
  createStationOfficer,
  getAuthenticatedClient,
} from './utils/utils';
import { v4 } from 'uuid';

describe('station-officer', () => {
  let client: Client;
  beforeAll(async () => {
    client = await getAuthenticatedClient({ includeProjectId: true });
  });

  it('gets by search', async () => {
    const stationLocation = await createStationLocation();
    const createdStationOfficer = await createStationOfficer(
      stationLocation.id
    );
    const stationOfficersResponse = await client.query({
      getStationOfficers: {
        __scalar: true,
        __args: {
          search: createdStationOfficer.firstname,
          stationlocationId: stationLocation.id,
          skip: 0,
          limit: 20,
        },
      },
    });

    const stationOfficers =
      stationOfficersResponse.getStationOfficers;
    expect(stationOfficers.length).toEqual(1);
    expect(stationOfficers[0]).toEqual(createdStationOfficer);
  });

  it('throws error when delete by wrong id', () => {
    expect(
      client.mutation({
        deleteStationOfficer: {
          __scalar: true,
          __args: {
            id: v4(),
          },
        },
      })
    ).rejects.toThrow();
  });

  it('deletes an existing', async () => {
    const stationLocation = await createStationLocation();
    const createdStationOfficer = await createStationOfficer(
      stationLocation.id
    );
    await client.mutation({
      deleteStationOfficer: {
        __scalar: true,
        __args: {
          id: createdStationOfficer.id,
        },
      },
    });
    expect(
      client.query({
        getStationOfficerById: {
          __scalar: true,
          __args: {
            id: createdStationOfficer.id,
          },
        },
      })
    ).rejects.toThrow();
  });

  it('throws error when get by wrong id', () => {
    expect(
      client.query({
        getStationOfficerById: {
          __scalar: true,
          __args: {
            id: v4(),
          },
        },
      })
    ).rejects.toThrow();
  });

  it('gets by id', async () => {
    const stationLocation = await createStationLocation();
    const createdStationOfficer = await createStationOfficer(
      stationLocation.id
    );

    const stationOfficerResponse = await client.query({
      getStationOfficerById: {
        __scalar: true,
        __args: {
          id: createdStationOfficer.id,
        },
      },
    });

    expect(stationOfficerResponse.getStationOfficerById).toEqual(
      createdStationOfficer
    );
  });

  it('throws error when update by wrong id', () => {
    expect(
      client.mutation({
        updateStationOfficer: {
          __scalar: true,
          __args: {
            id: v4(),
          },
        },
      })
    ).rejects.toThrow();
  });

  it('updates an existing', async () => {
    const stationLocation = await createStationLocation();
    const createdStationOfficer = await createStationOfficer(
      stationLocation.id
    );
    const firstname = nanoid();
    const lastname = nanoid();
    const msisdn = nanoid(10);

    const stationOfficerResponse = await client.mutation({
      updateStationOfficer: {
        __scalar: true,
        __args: {
          id: createdStationOfficer.id,
          firstname,
          lastname,
          msisdn,
        },
      },
    });

    const stationOfficer =
      stationOfficerResponse.updateStationOfficer;

    expect(stationOfficer.firstname).toEqual(firstname);
    expect(stationOfficer.lastname).toEqual(lastname);
    expect(stationOfficer.msisdn).toEqual(msisdn);
  });

  it('creates', async () => {
    const stationLocation = await createStationLocation();
    const firstname = nanoid();
    const lastname = nanoid();
    const msisdn = nanoid(10);

    const stationOfficerResponse = await client.mutation({
      createStationOfficer: {
        __scalar: true,
        __args: {
          stationLocationId: stationLocation.id,
          firstname,
          lastname,
          msisdn,
        },
      },
    });

    const stationOfficer =
      stationOfficerResponse.createStationOfficer;

    expect(stationOfficer.firstname).toEqual(firstname);
    expect(stationOfficer.lastname).toEqual(lastname);
    expect(stationOfficer.msisdn).toEqual(msisdn);
  });
});

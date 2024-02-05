import { AxiosInstance } from 'axios';
import { nanoid } from 'nanoid';
import { v4 } from 'uuid';
import { PROJECT_ID } from './utils/constants';
import { getAuthenticatedClient, prismaClient } from './utils/utils';

describe('Imei-configuration', () => {
  let client: AxiosInstance;
  let stationLocationId: string;

  beforeAll(async () => {
    client = await getAuthenticatedClient({ includeProjectId: true });
    const stationLocation = await prismaClient.stationLocation.create(
      {
        data: {
          id: v4(),
          title: nanoid(),
          projectId: PROJECT_ID,
          latitude: 0,
          longtitude: 0,
          createdBy: 'TEST',
        },
      }
    );

    stationLocationId = stationLocation.id;
  });

  async function createImei() {
    const imei = nanoid();
    const tags = [nanoid(), nanoid()];
    const imeiResponse = await client.post('/imei-configurations', {
      imei,
      label: 'NONE',
      priority: 'NORMAL',
      stationLocationId: stationLocationId,
      tags,
    });

    return imeiResponse.data.data;
  }

  it('deletes an existing', async () => {
    const imei = await createImei();
    const imeiResponse = await client.delete(
      `/imei-configurations/${imei.id}`
    );

    expect(imeiResponse.data.message).toEqual('deleted');
  });

  it('gets by tags', async () => {
    const imei = await createImei();
    const imeiResponse = await client.get(`/imei-configurations`, {
      params: {
        tags: [imei.tags[0]],
        limit: 20,
        skip: 0,
      },
    });
    const imeiResponseData = imeiResponse.data.data;
    expect(imeiResponseData[0]).toEqual(imei);
  });

  it('gets an existing', async () => {
    const imei = await createImei();
    const imeiResponse = await client.get(
      `/imei-configurations/${imei.id}`
    );
    const imeiResponseData = imeiResponse.data.data;
    expect(imeiResponseData.id).toEqual(imei.id);
  });

  it('throws error when get by wrong id', async () => {
    try {
      await client.get(`/imei-configurations/${v4()}`);
    } catch (error: any) {
      expect(error.response.status).toEqual(403);
      expect(error.response.data.message).toEqual(
        'forbidden operation'
      );
    }
  });

  it('throws error when update by wrong id', async () => {
    try {
      const newImei = nanoid();
      const tags = [nanoid(), nanoid()];
      await client.patch(`/imei-configurations/${v4()}`, {
        imei: newImei,
        label: 'NONE',
        priority: 'WARNING',
        tags,
      });
    } catch (error: any) {
      expect(error.response.status).toEqual(403);
      expect(error.response.data.message).toEqual(
        'forbidden operation'
      );
    }
  });

  it('updates an existing', async () => {
    const imei = await createImei();
    const newImei = nanoid();
    const tags = [nanoid(), nanoid()];
    const imeiResponse = await client.patch(
      `/imei-configurations/${imei.id}`,
      {
        imei: newImei,
        label: 'NONE',
        priority: 'WARNING',
        tags,
      }
    );

    const imeiResponseData = imeiResponse.data.data;
    expect(imeiResponseData.id).toEqual(imei.id);
    expect(imeiResponseData.tags).toEqual(tags);
    expect(imeiResponseData.imei).toEqual(newImei);
    expect(imeiResponseData.priority).toEqual('WARNING');
    expect(imeiResponseData.updatedBy).toBeDefined();
    expect(imeiResponseData.updatedAt).toBeDefined();
  });

  it('creates', async () => {
    const imei = nanoid();
    const tags = [nanoid(), nanoid()];
    const imeiResponse = await client.post('/imei-configurations', {
      imei,
      label: 'NONE',
      priority: 'NORMAL',
      stationLocationId: stationLocationId,
      tags,
    });

    const imeiResponseData = imeiResponse.data.data;
    expect(imeiResponseData.imei).toEqual(imei);
    expect(imeiResponseData.tags).toEqual(tags);
    expect(imeiResponse.data.message).toEqual('created');
  });
});

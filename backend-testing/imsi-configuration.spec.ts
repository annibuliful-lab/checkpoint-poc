import { AxiosInstance } from 'axios';
import { getAuthenticatedClient, prismaClient } from './utils/utils';
import { v4 } from 'uuid';
import { nanoid } from 'nanoid';
import { PROJECT_ID } from './utils/constants';

describe('Imsi-configuration', () => {
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

  async function createImsi() {
    const imsi = nanoid();
    const tags = [nanoid(), nanoid()];

    const imsiConfigurationResponse = await client.post(
      '/imsi-configurations',
      {
        imsi,
        label: 'NONE',
        priority: 'NORMAL',
        stationLocationId: stationLocationId,
        tags,
      }
    );

    return imsiConfigurationResponse.data.data;
  }

  it('gets by label', async () => {
    await createImsi();

    const response = await client.get('/imsi-configurations', {
      params: {
        label: 'NONE',
        limit: 1,
        skip: 0,
      },
    });

    const imsiConfigurations = response.data.data;
    expect(imsiConfigurations.length).toBeGreaterThanOrEqual(1);
    expect(
      imsiConfigurations.every((im: any) => im.label === 'NONE')
    ).toBeTruthy();
  });

  it('gets by mnc', async () => {
    const imsi = await createImsi();

    const response = await client.get('/imsi-configurations', {
      params: {
        mnc: imsi.mnc,
        limit: 20,
        skip: 0,
      },
    });

    const imsiConfigurations = response.data.data;
    expect(imsiConfigurations.length).toEqual(1);
    expect(imsi).toEqual(imsiConfigurations[0]);
  });

  it('gets by search', async () => {
    const imsi = await createImsi();

    const response = await client.get('/imsi-configurations', {
      params: {
        search: imsi.imsi,
        limit: 20,
        skip: 0,
      },
    });

    const imsiConfigurations = response.data.data;
    expect(imsiConfigurations.length).toEqual(1);
    expect(imsi).toEqual(imsiConfigurations[0]);
  });

  it('gets by mcc', async () => {
    const imsi = await createImsi();

    const response = await client.get('/imsi-configurations', {
      params: {
        mcc: imsi.mcc,
        limit: 20,
        skip: 0,
      },
    });

    const imsiConfigurations = response.data.data;
    expect(imsiConfigurations.length).toEqual(1);
    expect(imsi).toEqual(imsiConfigurations[0]);
  });

  it('gets by mnc', async () => {
    const imsi = await createImsi();

    const response = await client.get('/imsi-configurations', {
      params: {
        mnc: imsi.mnc,
        limit: 20,
        skip: 0,
      },
    });

    const imsiConfigurations = response.data.data;
    expect(imsiConfigurations.length).toEqual(1);
    expect(imsi).toEqual(imsiConfigurations[0]);
  });

  it('gets by tags', async () => {
    const imsi = await createImsi();

    const response = await client.get('/imsi-configurations', {
      params: {
        tags: [imsi.tags[0]],
        limit: 20,
        skip: 0,
      },
    });

    const imsiConfigurations = response.data.data;
    expect(imsiConfigurations.length).toBeGreaterThanOrEqual(1);
    expect(imsi).toEqual(imsiConfigurations[0]);
  });

  it('throws error when get by wrong id', async () => {
    try {
      await client.get(`/imsi-configurations/${v4()}`);
    } catch (error: any) {
      expect(error.response.status).toEqual(403);
      expect(error.response.data.message).toEqual(
        'forbidden operation'
      );
    }
  });

  it('gets an existing by id', async () => {
    const imsi = await createImsi();

    const imsiResponse = await client.get(
      `/imsi-configurations/${imsi.id}`
    );
    expect(imsiResponse.data.data).toEqual(imsi);
  });

  it('deletes an existing', async () => {
    const imsi = await createImsi();
    const deletedImsiResponse = await client.delete(
      `/imsi-configurations/${imsi.id}`
    );

    const deletedImsiResponseData = deletedImsiResponse.data;
    expect(deletedImsiResponseData.message).toEqual('deleted');
  });

  it('throws not found when update wrong id', async () => {
    const newImsi = nanoid();
    const newTags = [nanoid(), nanoid()];
    const newLabel = 'WHITELIST';
    const priority = 'WARNING';
    try {
      await client.patch(`/imsi-configurations/${v4()}`, {
        imsi: newImsi,
        label: newLabel,
        priority: priority,
        tags: newTags,
      });
    } catch (error: any) {
      expect(error.response.status).toEqual(403);
      expect(error.response.data.message).toEqual(
        'forbidden operation'
      );
    }
  });

  it('updates an existing', async () => {
    const imsi = nanoid();
    const tags = [nanoid(), nanoid()];

    const createdImsiConfigurationResponse = await client.post(
      '/imsi-configurations',
      {
        imsi,
        label: 'NONE',
        priority: 'NORMAL',
        stationLocationId: stationLocationId,
        tags,
      }
    );
    const newImsi = nanoid();
    const newTags = [nanoid(), nanoid()];
    const newLabel = 'WHITELIST';
    const priority = 'WARNING';
    const updatedImsiConfigurationResponse = await client.patch(
      `/imsi-configurations/${createdImsiConfigurationResponse.data.data.id}`,
      {
        imsi: newImsi,
        label: newLabel,
        priority: priority,
        tags: newTags,
      }
    );
    const imsiConfiguration =
      updatedImsiConfigurationResponse.data.data;

    expect(imsiConfiguration.id).toBeDefined();
    expect(imsiConfiguration.imsi).toEqual(newImsi);
    expect(imsiConfiguration.priority).toEqual(priority);
    expect(imsiConfiguration.label).toEqual(newLabel);
    expect(imsiConfiguration.tags).toEqual(newTags);
    expect(imsiConfiguration.updatedBy).toBeDefined();
    expect(imsiConfiguration.updatedAt).toBeDefined();
  });

  it('creates new', async () => {
    const imsi = nanoid();
    const tags = [nanoid(), nanoid()];

    const imsiConfigurationResponse = await client.post(
      '/imsi-configurations',
      {
        imsi,
        label: 'NONE',
        priority: 'NORMAL',
        stationLocationId: stationLocationId,
        tags,
      }
    );

    const imsiConfiguration = imsiConfigurationResponse.data.data;
    expect(imsiConfiguration.id).toBeDefined();
    expect(imsiConfiguration.imsi).toEqual(imsi);
    expect(imsiConfiguration.priority).toEqual('NORMAL');
    expect(imsiConfiguration.label).toEqual('NONE');
    expect(imsiConfiguration.tags).toEqual(tags);
  });
});

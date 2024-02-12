import axios from 'axios';

export const BACKEND_ENDPOINT = 'http://localhost:3030/';

export const httpClient = axios.create({
  baseURL: BACKEND_ENDPOINT,
});

export const PROJECT_ID = '246bb085-8ccc-4def-ac78-dc2ad5c7760b';
export const STATION_LOCATION_ID =
  'e1c6783c-e09c-43dd-b1e5-8041dcd2816e';

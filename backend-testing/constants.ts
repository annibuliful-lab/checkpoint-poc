import axios from 'axios';

export const BACKEND_ENDPOINT = 'http://localhost:3030/';

export const httpClient = axios.create({
  baseURL: BACKEND_ENDPOINT,
});

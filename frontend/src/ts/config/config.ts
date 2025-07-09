import { Staging } from './types';

type ConfigShape = {
  API_URL: string;
  STAGE: Staging;
};

export const Config: ConfigShape = {
  API_URL: 'http://localhost:8080',
  STAGE: window.STAGE as Staging,
};

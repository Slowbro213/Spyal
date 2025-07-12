import { Staging } from './types';

type ConfigShape = {
  API_URL: string;
  STAGE: Staging;
  LOG_ROUTE: string;
};

export const Config: ConfigShape = {
  API_URL: 'http://192.168.1.23:8080',
  STAGE: window.STAGE as Staging,
  LOG_ROUTE: 'http://192.168.1.23:8080/api/log',
};

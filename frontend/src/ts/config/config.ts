import { Staging } from './types';

type ConfigShape = {
  HOST: string;
  API_URL: string;
  STAGE: Staging;
  LOG_ROUTE: string;
  POKED_WS_SERVER: string;
  IS_DEVELOPMENT: boolean;
  IS_PRODUCTION: boolean;
};

// Environment detection
const getStage = (): Staging => {
  if (window.STAGE) return window.STAGE as Staging;
  if (import.meta.env?.MODE === Staging.Development) return Staging.Development;
  if (import.meta.env?.MODE === Staging.Production) return Staging.Production;
  return Staging.Development;
};

const stage = getStage();
const HOST = 'localhost:8080';

const PROTOCOL = stage === Staging.Production ? 'https' : 'http';

export const Config: ConfigShape = {
  HOST,

  API_URL: `${PROTOCOL}://${HOST}`,
  LOG_ROUTE: `${PROTOCOL}://${HOST}/api/log`,

  POKED_WS_SERVER: `${stage === Staging.Production ? 'wss' : 'ws'}://${HOST}/poked`,

  STAGE: stage,

  IS_DEVELOPMENT: stage === Staging.Development,
  IS_PRODUCTION: stage === Staging.Production,
};

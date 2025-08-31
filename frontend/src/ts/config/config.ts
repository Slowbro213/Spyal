import { Staging } from './types';

type ConfigShape = {
  HOST: string;
  API_URL: string;
  STAGE: Staging;
  LOG_ROUTE: string;
  IS_DEVELOPMENT: boolean;
  IS_PRODUCTION: boolean;
};

// Environment detection
const getStage = (): Staging => {
  if (window.STAGE) return window.STAGE as Staging;
  if (import.meta.env?.MODE === Staging.Development) return Staging.Development;
  if (import.meta.env?.MODE === Staging.Production) return Staging.Production;
  return Staging.Development; // default fallback
};

const stage = getStage();
const HOST = 'localhost:8080';

export const Config: ConfigShape = {
  HOST: HOST,
  API_URL: stage === Staging.Production ? `https://${HOST}` : `http://${HOST}`,

  STAGE: stage,

  LOG_ROUTE:
    stage === Staging.Production
      ? `https://${HOST}/api/log`
      : `http://${HOST}/api/log`,

  IS_DEVELOPMENT: stage === Staging.Development,
  IS_PRODUCTION: stage === Staging.Production,
};

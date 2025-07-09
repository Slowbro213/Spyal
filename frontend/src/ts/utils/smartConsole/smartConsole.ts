import { Staging, Config } from '@alspy/config';

export const log = (s: string) => {
  if (Config.STAGE === Staging.Development) console.log(s);
};

export { Config } from './config';
export { Staging } from './types';

declare global {
  interface Window {
    STAGE: string;
  }
}

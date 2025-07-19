import { pageRemoteCache, pageRemoteDestroy, pageRemoteInit } from './remote';
import type { Page } from '@alspy/pages/types';

export const remotePage: Page = {
  pageInit: pageRemoteInit,
  pageDestroy: pageRemoteDestroy,
  pageCache: pageRemoteCache,
};

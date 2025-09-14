import { pageLoginCache, pageLoginDestroy, pageLoginInit } from './login';
import type { Page } from '@alspy/pages/types';

export const loginPage: Page = {
  pageInit: pageLoginInit,
  pageDestroy: pageLoginDestroy,
  pageCache: pageLoginCache,
};

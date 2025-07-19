import { Page } from '@alspy/pages/types';
import { pageRoomCache, pageRoomInit, pageRoomDestroy } from './room';

export const roomPage: Page = {
  pageInit: pageRoomInit,
  pageDestroy: pageRoomDestroy,
  pageCache: pageRoomCache,
};

import { remotePage } from './remote';
import { roomPage } from './room';
import type { Page } from './types';

export const isCacheValid = (page: Page | undefined) => {
  if (!page) return true;
  const now = Date.now();
  const cacheExpired =
    typeof page.lastVisited === 'number' &&
    now > page.lastVisited + page.pageCache();

  return !cacheExpired;
};

const pages: Record<string, Page> = {
  '/create/remote': remotePage,
  '/room/*': roomPage,
};

export const GetPage = (location: string): Page | undefined => {
  // Try direct match first
  if (pages[location]) return pages[location];

  // Fallback to wildcard/regex-style matching
  for (const path in pages) {
    if (path.includes('*')) {
      const pattern = path.replace(/\*/g, '.*'); // turn '/room/*' into '/room/.*'
      const regex = new RegExp(`^${pattern}$`);
      if (regex.test(location)) return pages[path];
    }
  }

  return undefined; // No match
};

export const onPageChange = async () => {
  const location: string = window.location.pathname;
  const page = GetPage(location);
  if (!page) return;
  page.pageInit();
};

export const beforePageChange = async () => {
  const location: string = window.location.pathname;
  const page = GetPage(location);
  if (!page) return;
  page.pageDestroy();
};

import { pageCache } from './pageCache';
import { log } from './logger';
import { GetPage, isCacheValid } from '@alspy/pages';

export const fetchPage = async (href: string) => {
  const page = GetPage(href);
  const cacheValid = isCacheValid(page);
  if (pageCache.has(href) && cacheValid) {
    return;
  }
  try {
    const res = await fetch(href, {
      headers: {
        'X-Smart-Link': 'true',
      },
    });
    if (!res.ok) throw new Error(`Failed to fetch ${href}`);
    const html = await res.text();

    pageCache.set(href, html);
    if (page) page.lastVisited = Date.now();
  } catch (err) {
    log({
      level: 'error',
      msg: `Failed to fetch page: ${href}, Error: ${err}`,
    });
  }
};

import { pageCache } from './pageCache';
import { log } from './logger';
import { GetPage, isCacheValid } from '@alspy/pages';

export const fetchPage = async (href: string): Promise<string> => {
  const page = GetPage(href);
  const cacheValid = isCacheValid(page);
  if (pageCache.has(href) && cacheValid) {
    return href;
  }
  try {
    const res = await fetch(href, {
      headers: {
        'X-Smart-Link': 'true',
      },
    });
    if (!res.ok) throw new Error(`Failed to fetch ${href}`);
    let location;
    if (res.redirected) {
      location = res.url;
    } else {
      location = href;
    }
    const html = await res.text();

    pageCache.set(location, html);
    if (page) page.lastVisited = Date.now();
    return location;
  } catch (err) {
    log({
      level: 'error',
      msg: `Failed to fetch page: ${href}, Error: ${err}`,
    });
    return href;
  }
};

import { pageCache } from './pageCache';

export const fetchPage = async (href: string) => {
  if (pageCache.has(href)) {
    return;
  }
  try {
    const res = await fetch(href, {
      headers: { 'X-Smart-Link': 'true' },
    });
    if (!res.ok) throw new Error(`Failed to fetch ${href}`);
    const html = await res.text();

    pageCache.set(href, html);
  } catch (err) {
    console.error('SmartLink fetch failed:', err);
  }
};

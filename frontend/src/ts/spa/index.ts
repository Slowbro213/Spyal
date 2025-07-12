import { pageCache } from './pageCache';
import { prefetchRoutes } from './config';

export { loadPage } from './loadPage';
export { fetchPage } from './fetchPage';
export { log } from './logger';

window.addEventListener('popstate', async () => {
  const href = location.pathname;
  const target = document.querySelector('#content');
  if (!target) return;

  if (pageCache.has(href)) {
    target.innerHTML = pageCache.get(href)!;
  } else {
    const res = await fetch(href, { headers: { 'X-Smart-Link': 'true' } });
    const html = await res.text();
    pageCache.set(href, html);
    target.innerHTML = html;
  }
});

async function prefetchPages() {
  for (const route of prefetchRoutes) {
    try {
      if (pageCache.has(route)) continue;
      const res = await fetch(route, {
        headers: { 'X-Smart-Link': 'true' },
      });
      if (!res.ok) throw new Error(`Failed to prefetch ${route}`);
      const html = await res.text();
      pageCache.set(route, html);
    } catch (err) {
      console.warn('Prefetch error for', route, err);
    }
  }
}

document.addEventListener('DOMContentLoaded', () => {
  const path = location.pathname;
  const target = document.querySelector('#content');
  if (target && !pageCache.has(path)) {
    pageCache.set(path, target.innerHTML);
  }
  prefetchPages();
});

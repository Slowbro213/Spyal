import { pageCache } from './pageCache';
import { target } from './config';
import { serveErrorPage } from './error';
import { onPageChange } from '@alspy/pages';
export const loadPage = async (href: string) => {
  try {
    const html = pageCache.get(href);
    if (target === null || !html) {
      await serveErrorPage('huge');
      return;
    }
    target.innerHTML = html;
    history.pushState(null, '', href);
    onPageChange();
  } catch (err) {
    console.error('SmartLink fetch failed:', err);
    serveErrorPage('catastrofic');
  }
};

import { pageCache } from './pageCache';
import { target } from './config';
import { serveError } from './error';
import { onPageChange } from '../pages';
export const loadPage = async (href: string) => {
  try {
    const html = pageCache.get(href);
    if (target === null || !html) {
      await serveError();
      return;
    }
    target.innerHTML = html;
    history.pushState(null, '', href);
    onPageChange();
  } catch (err) {
    console.error('SmartLink fetch failed:', err);
    serveError();
  }
};

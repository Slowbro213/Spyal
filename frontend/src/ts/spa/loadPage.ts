import { pageCache } from './pageCache';
import { target } from './config';
import { serveErrorPage, Severity } from './error';
import { beforePageChange, onPageChange } from '@alspy/pages';
import { log } from '@alspy/spa';

export const loadPage = async (href: string) => {
  try {
    const html = pageCache.get(href);
    if (target === null || !html) {
      log({
        level: 'error',
        msg: `Failed to Load ${href}, Error: no target`,
      });
      await serveErrorPage(Severity.small);
      return;
    }
    beforePageChange();
    target.innerHTML = html;
    history.pushState(null, '', href);
    onPageChange();
  } catch (err) {
    log({
      level: 'error',
      msg: `Failed to Load ${href}, Error: ${err}`,
    });
    serveErrorPage(Severity.catastrofic);
  }
};

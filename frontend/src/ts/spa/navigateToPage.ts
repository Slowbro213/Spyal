import { fetchPage } from './fetchPage';
import { loadPage } from './loadPage';

export const navigateToPage = async (href: string) => {
  fetchPage(href).then(() => {
    loadPage(href);
  });
};

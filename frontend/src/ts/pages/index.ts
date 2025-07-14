import { pageRemoteInit, pageRemoteDestroy } from './remote';

const pageInits: Record<string, () => void> = {
  '/create/remote': pageRemoteInit,
};

const pageDestructions: Record<string, () => void> = {
  '/create/remote': pageRemoteDestroy,
};

export const onPageChange = async () => {
  const location: string = window.location.pathname;
  if (!pageInits[location]) return;
  pageInits[location]?.();
};

export const beforePageChange = async () => {
  const location: string = window.location.pathname;
  if (!pageDestructions[location]) return;
  pageDestructions[location]?.();
};

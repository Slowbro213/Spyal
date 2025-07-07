import { pageRemoteInit } from './remote';

const pageInits: Record<string, () => void> = {
  '/create/remote': pageRemoteInit,
};

export const onPageChange = async () => {
  const location: string = window.location.pathname;
  if (!pageInits[location]) console.log('No code for this yet');
  pageInits[location]?.();
};

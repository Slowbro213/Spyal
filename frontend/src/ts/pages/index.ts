import { pageRemoteInit } from './remote';

const pageInits: Record<string, () => void> = {
  '/create/remote': pageRemoteInit,
};

export const onPageChange = () => {
  const location: string = window.location.pathname;
  console.log(location);
  pageInits[location]?.();
};

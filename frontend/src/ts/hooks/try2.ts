import { test } from 'utils/try';

export const test2 = (): string => {
  const hello = test();
  const world = 'World!';

  return hello + ' ' + world;
};

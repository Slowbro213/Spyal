import { initButton } from './button';
export * from './smart-link';

export const componentMap = new Map<string, () => void>();
componentMap.set('button', initButton);

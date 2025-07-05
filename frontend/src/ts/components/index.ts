import { initButton } from './button';

export const componentMap = new Map<string, () => void>();
componentMap.set('button', initButton);

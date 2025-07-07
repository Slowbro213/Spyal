declare global {
  interface Window {
    STAGE: string;
  }
}

import './components';
import { onPageChange } from './pages';
import './spa';

onPageChange();

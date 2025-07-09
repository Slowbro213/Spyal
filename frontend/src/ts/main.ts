import { Importance, initToast, Level } from './services';
import './components';
import { onPageChange } from './pages';
import './spa';

onPageChange();
const toast = initToast();
toast?.show(Level.Success, Importance.Major, {
  title: 'Testing Toast',
  message: 'Testing Toast',
});

import { Toast } from '@alspy/components/toast';

let toast: Toast;

export const initToast = (): Toast => {
  if (toast) return toast;
  toast = document.querySelector('toast-service') as Toast;
  if (!toast) {
    console.log('No toast :(');
  }
  return toast;
};

import { Importance, initToast, Level } from '@alspy/services/toast';
import { navigateToPage } from '@alspy/spa';
import type { LoginForm } from './types';


async function handleFormSubmit(event: Event) {
  event.preventDefault();
  const toast = initToast();

  const { data } = (event as CustomEvent<{ data: Record<string, string> }>).detail;

  const formData: LoginForm = {
    username: data.username?.trim() ?? '',
    password: data.password?.trim() ?? '',
  };

  if (!formData.username || !formData.password) {
    toast.show(Level.Error, Importance.Major, { message: 'Plotësoni të gjitha fushat!' });
    return;
  }

  try {
    const body = new URLSearchParams();
    body.set('username', formData.username);
    body.set('password', formData.password);

    const res = await fetch('/login', {
      method: 'POST',
      body: body,
    });

    if (!res.ok) throw Error()

    const username = formData.username;

    document.dispatchEvent(
      new CustomEvent('auth:authenticated', {
        detail: { username },
        bubbles: true,
        composed: true
      })
    );


    const next = new URLSearchParams(window.location.search).get('next') || '/';
    navigateToPage(next);
  } catch (err: any) {
    toast.show(Level.Error, Importance.Major, { message: err.message || 'Gabim në hyrje' });
  }
}

export const pageLoginInit = () => {
  document.getElementById('login-form')?.addEventListener('smart-form:submit', handleFormSubmit);
};

export const pageLoginDestroy = () => {
  document.getElementById('login-form')?.removeEventListener('smart-form:submit', handleFormSubmit);
};

export const pageLoginCache = () => 20_000_000;

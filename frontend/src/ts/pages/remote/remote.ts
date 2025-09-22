import { client } from '@alspy/api';
import { Importance, initToast, Level } from '@alspy/services/toast';
import type { RemoteGameCreationResponse, RemoteGameForm } from './types';
import { navigateToPage } from '@alspy/spa';

async function handleFormSubmit(event: Event) {
  event.preventDefault();

  const toast = initToast();

  const { data } = (event as CustomEvent<{ data: Record<string, string> }>)
    .detail;

  const gameName = data['game-name']?.trim() || '';
  const spyNumber = Number(data['spy-number']);
  const maxPlayers = Number(data['max-players']);
  const isPrivate = data['game-visibility'] === 'private';

  if (!spyNumber || !maxPlayers || !('game-visibility' in data)) {
    toast.show(Level.Error, Importance.Major, {
      message: 'Plotësoni të gjitha fushat e detyrueshme!',
    });
    return;
  }

  const form: RemoteGameForm = {
    gameName,
    spyNumber,
    maxNumbers: maxPlayers,
    isPrivate,
  };

  try {
    const response: RemoteGameCreationResponse = await client.post(
      '/create/remote',
      form,
      {
        redirect: "manual",
      }
    );

    navigateToPage(`/room/${response.roomID}`);
  } catch (err: any) {
    toast.show(Level.Error, Importance.Major, {
      message: err.message
    })
  }
}

export const pageRemoteInit = () => {
  const formEl = document.getElementById('create-online-form');
  if (!formEl) return;

  formEl.addEventListener('smart-form:submit', handleFormSubmit);
};

export const pageRemoteDestroy = () => {
  const formEl = document.getElementById('create-online-form');
  if (!formEl) return;

  formEl.removeEventListener('smart-form:submit', handleFormSubmit);
};

export const pageRemoteCache = (): number => {
  const seconds = 20;
  return seconds * 1000;
};

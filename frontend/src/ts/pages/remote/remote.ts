import { client } from '@alspy/api';
import { Importance, initToast, Level } from '@alspy/services/toast';
import type { RemoteGameForm } from './types';

function handleFormSubmit(event: Event) {
  event.preventDefault();

  const toast = initToast();

  const { data } = (event as CustomEvent<{ data: Record<string, string> }>)
    .detail;

  const playerName = data['player-nickname']?.trim() || '';
  const gameNameRaw = data['game-name']?.trim() || '';
  const spyNumber = Number(data['spy-number']);
  const maxPlayers = Number(data['max-players']);
  const isPrivate = data['game-visibility'] === 'private';

  if (
    !playerName ||
    !spyNumber ||
    !maxPlayers ||
    !('game-visibility' in data)
  ) {
    toast.show(Level.Error, Importance.Major, {
      message: 'Plotësoni të gjitha fushat e detyrueshme!',
    });
    return;
  }

  const gameName =
    gameNameRaw || (playerName ? `Dhoma e ${playerName}` : 'Spyfall Dhoma');

  const form: RemoteGameForm = {
    playerName,
    gameName,
    spyNumber,
    maxNumbers: maxPlayers,
    isPrivate,
  };

  client.post('/create/remote', {
    params: { ...form },
  });
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

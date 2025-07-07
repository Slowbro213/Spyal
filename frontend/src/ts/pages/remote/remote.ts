import { client } from '@alspy/api';
import type { RemoteGameForm } from './types';

export const pageRemoteInit = () => {
  // Use the smart-form element, not the native form
  const formEl = document.getElementById(
    'create-online-form'
  ) as HTMLElement | null;
  if (!formEl) return;

  formEl.addEventListener('smart-form:submit', (event: Event) => {
    event.preventDefault();

    // Get plain data object from custom event
    const { data } = (event as CustomEvent<{ data: Record<string, string> }>)
      .detail;

    // Extract/parse and validate
    const playerName = data['player-nickname']?.trim() || '';
    const gameNameRaw = data['game-name']?.trim() || '';
    const gameDuration = Number(data['game-duration']);
    const maxPlayers = Number(data['max-players']);
    const isPrivate = data['game-visibility'] === 'private';

    // Validation: Required fields
    if (
      !playerName ||
      !gameDuration ||
      !maxPlayers ||
      !('game-visibility' in data)
    ) {
      alert('Plotësoni të gjitha fushat e detyrueshme!');
      return;
    }

    // Generate gameName if not set
    const gameName =
      gameNameRaw || (playerName ? `Dhoma e ${playerName}` : 'Spyfall Dhoma');

    // Build the form object (type safe)
    const form: RemoteGameForm = {
      playerName,
      gameName,
      time: gameDuration,
      maxNumbers: maxPlayers,
      isPrivate,
    };

    // Do what you want with form (log, send, etc)
    console.log(form);

    client.post('/create/remote', {
      params: { ...form },
    });

    // Example: submit with fetch, call API, etc.
    // fetch("/your/api", { method: "POST", body: JSON.stringify(form), headers: { "Content-Type": "application/json" } });
  });
};

import { poker } from '@alspy/poke';
import { EventName } from '@alspy/poke/types';

export const pageRoomInit = () => {
  const roomID = window.location.pathname.split('/').pop();
  if (!roomID) return;

  const channel = poker.channel('game', roomID);

  channel.spy(EventName.Userjoinedevent, (payload: any) => {
    const counter = document.querySelector('#game-room .feature-card h2 span');
    if (counter) {
      const [, max] = counter.textContent!.match(/\((\d+)\/(\d+)\)/)!.slice(1);
      const newCount = Number(counter.textContent!.match(/(\d+)/)![1]) + 1;
      counter.textContent = `Lojtarët (${newCount}/${max})`;
    }

    const listBox = document.querySelector<HTMLElement>('#game-room .space-y-3');
    if (!listBox) return;

    const div = document.createElement('div');
    div.className = 'flex items-center justify-between p-3 rounded-lg transition-colors hover:bg-gray-100 dark:hover:bg-gray-700';
    div.style.cssText = 'background: var(--card-bg); border: 1px solid var(--border-accent);';
    div.innerHTML = `
    <span class="py-2 px-4">${payload.username}</span>
    <span class="text-[var(--accent-green)]"><i class="fas fa-check-circle"></i></span>
  `;
    listBox.appendChild(div);
  });
};

export const pageRoomDestroy = () => {
  ((window as any).roomChannel)?.leave();
};

export const pageRoomCache = (): number => -1;

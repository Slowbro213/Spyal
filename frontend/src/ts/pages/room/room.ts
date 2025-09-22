import { client } from '@alspy/api';
import { poker } from '@alspy/poke';
import { EventName } from '@alspy/poke/types';
import { log, navigateToPage } from '@alspy/spa';

function getCookie(name: string) {
  return document.cookie
    .split("; ")
    .find(row => row.startsWith(name + "="))
    ?.split("=")[1];
}

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

  const box = document.getElementById('chat-messages') as HTMLDivElement;
  if (!box) return;
  channel.spy(EventName.Chatevent, (payload: any) => {
    const { username, text } = payload.msg;

    const empty = box.querySelector('.text-center');
    if (empty) empty.remove();

    const bubble = document.createElement('div');
    bubble.innerHTML = `<strong>${username}:</strong> ${text}`;

    box.appendChild(bubble);
    box.scrollTop = box.scrollHeight;
  });

  const username = getCookie("username");
  const chatInput = document.getElementById('chat-input') as HTMLInputElement;
  document.getElementById('send-chat-btn')?.addEventListener('click', () => {
    const text = chatInput.value.trim();
    if (!text) return;
    channel.poke(EventName.Chatevent, {
      msg: {
        text,
        username
      },
    });
    chatInput.value = '';
  });


  document.getElementById('leave-room')?.addEventListener('click', () => {
    try {
      client.post("/leave");
      navigateToPage("/");
      channel.close();
    } catch {
      log({
        level: 'error',
        msg: `Failed to Leave Room`,
      });
    }

  });
};

export const pageRoomDestroy = () => {
  ((window as any).roomChannel)?.leave();
};

export const pageRoomCache = (): number => -1;

import { poker } from '@alspy/poke';
import { EventName } from '@alspy/poke/types';
import { log } from '@alspy/spa';
import type { Channel } from '@alspy/poke';

let echoChann: Channel;
let running = true;

export const pageRoomInit = () => {
  echoChann = poker.channel('game', 'test');

  echoChann.spy(EventName.Echoevent, (event: any) => {
    log({ level: 'info', msg: `listened! msg is: ${event.msg}` });
  });

  running = true;
  (async () => {
    while (running) {
      echoChann.poke(EventName.Echoevent, {
        msg: 'Hello! This should be echoed',
      });
      await new Promise((resolve) => setTimeout(resolve, 1000));
    }
  })();
};

export const pageRoomDestroy = () => {
  running = false;
  echoChann.close();
};

export const pageRoomCache = (): number => {
  return -1;
};

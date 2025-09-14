import { poker } from '@alspy/poke';
import { EventName } from '@alspy/poke/types';
import { log } from '@alspy/spa';
import type { Channel } from '@alspy/poke';

let echoChann: Channel;

export const pageRoomInit = () => {
  echoChann = poker.channel('game', 'test');

  echoChann.spy(EventName.Echoevent, (event: any) => {
    log({ level: 'info', msg: `listened! msg is: ${event.msg}` });
  });

};

export const pageRoomDestroy = () => {
  echoChann.close();
};

export const pageRoomCache = (): number => {
  return -1;
};

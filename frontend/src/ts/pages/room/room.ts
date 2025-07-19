import { Config, Staging } from '@alspy/config';
import { log } from '@alspy/spa';

let socket: WebSocket;

export const pageRoomInit = () => {
  // Create a WebSocket connection
  const protocol = Config.STAGE === Staging.Production ? 'wss' : 'ws';
  const path = '/echo';
  socket = new WebSocket(
    `${protocol}://192.168.1.23:8080${path}`,
    ['poked'] // Subprotocol array
  );

  // Event listeners for connection events
  socket.addEventListener('open', (event) => {
    console.log('WebSocket connection opened:', event);
    socket.send('Hello, server!'); // Send a message to the server
  });

  socket.addEventListener('message', (event) => {
    console.log('Received message:', event.data);
  });

  socket.addEventListener('close', (event) => {
    console.log('WebSocket connection closed:', event);
  });

  socket.addEventListener('error', (event) => {
    console.error('WebSocket error:', event);
  });

  socket.onclose = () => {
    log({
      level: 'info',
      msg: 'Game: User Left',
    });
  };
};

export const pageRoomDestroy = () => {
  const normalClosureCode = 1000;

  socket.close(normalClosureCode, 'User Left the Game');
};

export const pageRoomCache = (): number => {
  return -1; // Never cache this page!!
};

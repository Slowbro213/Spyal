import { Config } from "@alspy/config";
import { EventName } from "./types";

type Listener = (payload: any) => void;

export class Channel {
  private listeners: Map<number, Listener[]> = new Map();
  private socket: WebSocket;
  private sendQueue: string[] = [];

  constructor(socket: WebSocket) {
    this.socket = socket;

    this.socket.onopen = () => {
      // flush queued messages
      for (const msg of this.sendQueue) {
        try {
          this.socket.send(msg);
        } catch (e) {
          console.error("Failed to send queued WS message:", e);
        }
      }
      this.sendQueue = [];
    };

    this.socket.onmessage = (e) => {
      try {
        const msg = JSON.parse(e.data);
        const { type, msg: payload } = msg;
        this._notify(type, payload);
      } catch {
        console.warn("⚠️ Invalid WS message:", e.data);
      }
    };

    this.socket.onerror = (err) => {
      console.error("WebSocket error:", err);
    };

    this.socket.onclose = (ev) => {
      // optionally notify listeners about closure if desired
      // console.log("WebSocket closed:", ev);
    };
  }

  spy(type: EventName, callback: Listener): this {
    if (!this.listeners.has(type)) {
      this.listeners.set(type, []);
    }
    this.listeners.get(type)!.push(callback);
    return this;
  }

  poke(type: EventName, msg: any): void {
    const data = JSON.stringify({ type, msg });
    if (this.socket.readyState === WebSocket.OPEN) {
      this.socket.send(data);
    } else {
      // queue the message so it will be sent when the socket opens
      this.sendQueue.push(data);
    }
  }

  close(code?: number, reason?: string): void {
    try {
      // Close unless already closed
      if (this.socket.readyState !== WebSocket.CLOSING && this.socket.readyState !== WebSocket.CLOSED) {
        // if the caller provided a code, use it; otherwise just call close()
        if (typeof code === "number") {
          this.socket.close(code, reason);
        } else {
          this.socket.close();
        }
      }
    } catch (e) {
      console.error("Error closing WebSocket:", e);
    }

    // cleanup handlers and queued messages
    this.socket.onopen = null;
    this.socket.onmessage = null;
    this.socket.onerror = null;
    this.socket.onclose = null;

    this.sendQueue = [];
    this.listeners.clear();
  }

  private _notify(type: number, payload: any) {
    const handlers = this.listeners.get(type);
    if (handlers) {
      handlers.forEach(fn => {
        try {
          fn(payload);
        } catch (e) {
          console.error("Channel listener error:", e);
        }
      });
    }
  }
}

class Poker {
  private baseUrl: string;
  private protocol: string;
  private channels: Map<string, Channel> = new Map();

  constructor() {
    this.baseUrl = Config.POKED_WS_SERVER;
    this.protocol = "poked";
  }

  /**
   * Return a Channel for `name`. If there is an existing channel whose socket
   * is OPEN or CONNECTING, it's returned. If the previous socket is CLOSED
   * (or CLOSING), a new socket/channel is created and returned.
   */
  channel(name: string): Channel {
    const existing = this.channels.get(name);
    if (existing) {
      const rs = existing as Channel & { socket?: WebSocket };
      // try to access the socket readyState; fallback to returning existing
      try {
        const readyState = (rs as any).socket?.readyState;
        if (readyState === WebSocket.OPEN || readyState === WebSocket.CONNECTING) {
          return existing;
        }
        // If CLOSED or CLOSING, fallthrough to recreate
      } catch {
        // If we can't inspect readyState, just return existing
        return existing;
      }
    }

    const ws = new WebSocket(`${this.baseUrl}?channel=${name}`, this.protocol);
    const chan = new Channel(ws);
    this.channels.set(name, chan);
    return chan;
  }

  closeChannel(name: string): void {
    const chan = this.channels.get(name);
    if (chan) {
      chan.close();
      this.channels.delete(name);
    }
  }

  closeAll(): void {
    for (const [name, chan] of Array.from(this.channels.entries())) {
      chan.close();
      this.channels.delete(name);
    }
  }
}

export const poker = new Poker();

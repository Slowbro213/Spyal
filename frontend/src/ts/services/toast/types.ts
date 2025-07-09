export enum Level {
  Success,
  Info,
  Warning,
  Error,
}

export enum Importance {
  Major,
  Minor,
}

export type Message = {
  title?: string;
  message: string;
};

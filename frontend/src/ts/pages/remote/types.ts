export type RemoteGameForm = {
  gameName?: string;
  spyNumber: number;
  maxNumbers: number;
  isPrivate: boolean;
};

export type RemoteGameCreationResponse = {
  roomID: string;
};

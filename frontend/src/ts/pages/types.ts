export type Page = {
  pageInit: () => void;
  pageDestroy: () => void;
  pageCache: () => number; // amount of time the page should be cached for in ms
  lastVisited?: number;
};

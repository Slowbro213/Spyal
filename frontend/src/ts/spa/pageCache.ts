import { MAX_NUM_PAGES } from './config';

export class PageCache {
  private cache = new Map<string, string>();
  private maxSize: number;

  constructor(maxSize = 100) {
    // You can adjust or remove max size
    this.maxSize = maxSize;
  }

  get(href: string): string | undefined {
    return this.cache.get(href);
  }

  set(href: string, value: string): void {
    // Optional: enforce max size (LRU style eviction)
    if (this.cache.size >= this.maxSize) {
      // Remove the first inserted (oldest) item
      const firstKey = this.cache.keys().next().value;
      if (firstKey) this.cache.delete(firstKey);
    }
    this.cache.set(href, value);
  }

  has(href: string): boolean {
    return this.cache.has(href);
  }

  delete(href: string): void {
    this.cache.delete(href);
  }

  clear(): void {
    this.cache.clear();
  }

  // Optional: expose entries/size for debugging
  get size(): number {
    return this.cache.size;
  }

  entries(): IterableIterator<[string, string]> {
    return this.cache.entries();
  }
}

export const pageCache = new PageCache(MAX_NUM_PAGES);

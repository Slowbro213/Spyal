import { fetchPage, loadPage } from '@alspy/spa';

export class SmartLink extends HTMLElement {
  fetched = false;
  constructor() {
    super();
    this.attachShadow({ mode: 'open' }).innerHTML = `
      <style>
        :host {
          color: var(--link-color, #0077cc);
          cursor: pointer;
          display: inline;
        }
      </style>
      <slot></slot>
    `;
    const prefetchOption = this.getAttribute('prefetch')?.toLowerCase();

    const fetchOnLoad = prefetchOption === 'load';

    const location = this.getAttribute('href');
    if (fetchOnLoad && location) {
      fetchPage(location);
    }
  }

  connectedCallback() {
    this.addEventListener('click', this.handleClick.bind(this));
  }

  async handleClick(e: MouseEvent) {
    e.preventDefault();
    const location = this.getAttribute('href');
    if (!location) return;
    if (!this.fetched) await fetchPage(location);
    loadPage(location);
  }
}

customElements.define('smart-link', SmartLink);

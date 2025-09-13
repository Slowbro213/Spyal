import { Importance } from '@alspy/services';
import { fetchPage, loadPage } from '@alspy/spa';
import { serveErrorPage, Severity } from '@alspy/spa/error';

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
    try {
    if (!this.fetched) await fetchPage(location);
    loadPage(location);
    } catch (err) {
      serveErrorPage(Severity.normal,"Akses i Pa Autorizuar","Ti nuk ke akses ne kete faqe!");
    }
  }
}

customElements.define('smart-link', SmartLink);

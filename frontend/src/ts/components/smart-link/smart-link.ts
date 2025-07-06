import { pageCache } from '../../pageCache';

export class SmartLink extends HTMLElement {
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
  }

  connectedCallback() {
    this.addEventListener('click', this.handleClick.bind(this));
  }

  async handleClick(e: MouseEvent) {
    e.preventDefault();
    const href = this.getAttribute('href');
    if (!href) return;

    const target = document.querySelector('#content');
    if (!target) return;

    if (pageCache.has(href)) {
      target.innerHTML = pageCache.get(href)!;
      history.pushState(null, '', href);
      return;
    }

    try {
      const res = await fetch(href, {
        headers: { 'X-Smart-Link': 'true' },
      });
      if (!res.ok) throw new Error(`Failed to fetch ${href}`);
      const html = await res.text();

      pageCache.set(href, html);
      target.innerHTML = html;
      history.pushState(null, '', href);
    } catch (err) {
      console.error('SmartLink fetch failed:', err);
    }
  }
}

customElements.define('smart-link', SmartLink);

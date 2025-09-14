import { Config } from "@alspy/config";

export class SmartLogout extends HTMLElement {
  constructor() {
    super();
    this.attachShadow({ mode: 'open' }).innerHTML = `
      <style>
        :host {
          display: inline-flex;
        }
        button {
          padding: 0.5rem 1rem;          /* 8px 16px  – same as old px-4 py-2 */
          font-size: 0.875rem;           /* text-sm */
          font-weight: 600;              /* font-semibold */
          color: #fff;
          background: linear-gradient(to right, #f43f5e, #e11d48); /* rose-500 → red-600 */
          border: none;
          border-radius: 0.5rem;         /* rounded-lg */
          cursor: pointer;
          transition: all 200ms;
        }
        button:hover {
          background: linear-gradient(to right, #fb7185, #f43f5e); /* rose-400 → red-500 */
          box-shadow: 0 0 0 2px rgba(244, 63, 94, 0.6);          /* ring-rose-400 */
        }
        button:disabled {
          opacity: 0.6;
          cursor: not-allowed;
        }
      </style>
      <button part="button" type="button">
        <slot>Dil</slot>
      </button>
    `;
  }

  connectedCallback() {
    this.shadowRoot?.querySelector('button')?.addEventListener('click', this._logout.bind(this));
  }

  async _logout() {
    const btn = this?.shadowRoot?.querySelector('button');
    if(btn) btn.disabled = true;

    try {
      const res = await fetch('/logout', {
        method: 'POST',
        credentials: 'same-origin',
        headers: { 'X-Requested-With': 'smart-logout' }
      });

      if (!res.ok) throw new Error(res.statusText || 'Logout failed');

      const data = await res.json().catch(() => ({}));

      document.dispatchEvent(new CustomEvent('auth:expire', {
        detail: { response: data, res },
        bubbles: true,
        composed: true
      }));

      console.log("dispatched");

    } catch (err) {
    } finally {
      if(btn) btn.disabled = false;
    }
  }
}

customElements.define('smart-logout', SmartLogout);

  import { fetchPage, loadPage } from '@alspy/spa'; // Adjust path if necessary

  export class SmartSearch extends HTMLElement {
    private input: HTMLInputElement;
    private form: HTMLFormElement;

    constructor() {
      super();

      // Basic styling for the input and form
      this.attachShadow({ mode: 'open' }).innerHTML = `
        <style>
          :host {
            display: block;
          }
          form {
            display: flex;
            gap: 0.5rem;
          }
          input {
            flex-grow: 1;
            padding: 0.75rem 1rem;
            background-color: var(--background-light, #1f2937); /* gray-800 */
            border: 1px solid var(--border-accent, #4b5563); /* gray-600 */
            border-radius: 0.5rem; /* rounded-lg */
            color: var(--text-primary, #f9fafb); /* gray-50 */
            font-size: 1rem;
          }
          input::placeholder {
            color: var(--text-secondary, #9ca3af); /* gray-400 */
          }
          input:focus {
            outline: none;
            border-color: var(--primary-teal, #14b8a6); /* teal-500 */
            box-shadow: 0 0 0 2px var(--primary-teal-dark, #0f766e); /* teal-600 */
          }
          button {
            padding: 0.75rem 1.25rem;
            background: var(--primary-teal, #14b8a6);
            color: white;
            border: none;
            border-radius: 0.5rem;
            font-weight: 600;
            cursor: pointer;
            transition: opacity 0.2s;
          }
          button:hover {
            opacity: 0.9;
          }
        </style>
        <form>
          <input type="text" name="search" placeholder="${this.getAttribute('placeholder') || 'Search...'}" value="${this.getAttribute('value') || ''}">
          <button type="submit">Kërko</button>
        </form>
      `;

      this.form = this.shadowRoot?.querySelector('form')!;
      this.input = this.shadowRoot?.querySelector('input')!;

      this.form.addEventListener('submit', this.handleSubmit.bind(this));
      this.input.addEventListener('input', this.handleInput.bind(this));
    }

    private handleSubmit(e: Event) {
      e.preventDefault();
      this.performSearch();
    }

    private handleInput() {
    }

    private async performSearch() {
      const query = this.input.value.trim();
      const url = new URL(window.location.href);

      if (query) {
        url.searchParams.set('q', query);
      } else {
        url.searchParams.delete('q');
      }

      try {
        await fetchPage(url.pathname + url.search);
        loadPage(url.pathname + url.search);
      } catch (err) {
      }
    }
  }

  customElements.define('smart-search', SmartSearch);

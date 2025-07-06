export class SmartForm extends HTMLElement {
  constructor() {
    super();
  }

  connectedCallback() {
    this.addEventListener('submit', this.handleSubmit as EventListener);
    // Enable submit on Enter
    this.addEventListener('keydown', (e) => {
      if (e instanceof KeyboardEvent && e.key === 'Enter') {
        const target = e.target as HTMLElement;
        if (
          target &&
          (target.tagName === 'INPUT' || target.tagName === 'SELECT') &&
          !(target instanceof HTMLTextAreaElement)
        ) {
          e.preventDefault();
          this.handleSubmit(e);
        }
      }
    });
  }

  async handleSubmit(e: Event) {
    e.preventDefault();
    const method = (this.getAttribute('method') || 'GET').toUpperCase();
    const action = this.getAttribute('action') || window.location.pathname;

    // Collect all inputs, selects, and textareas within this form
    const elements = Array.from(
      this.querySelectorAll('input, select, textarea')
    ) as (HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement)[];
    const formData = new FormData();

    for (const el of elements) {
      if (el.name && !el.disabled) {
        if (
          el instanceof HTMLInputElement &&
          (el.type === 'checkbox' || el.type === 'radio')
        ) {
          if (el.checked) formData.append(el.name, el.value);
        } else {
          formData.append(el.name, el.value);
        }
      }
    }

    let url = action;
    let body: FormData | undefined = undefined;
    const headers: Record<string, string> = { 'X-Smart-Link': 'true' };

    if (method === 'GET') {
      const params = new URLSearchParams();
      for (const [key, value] of formData.entries()) {
        params.append(key, typeof value === 'string' ? value : String(value));
      }
      url += (url.includes('?') ? '&' : '?') + params.toString();
    } else {
      body = formData;
    }

    try {
      const res = await fetch(url, {
        method,
        headers,
        body: method === 'GET' ? undefined : body,
      });
      this.dispatchEvent(
        new CustomEvent('smart-form:success', {
          detail: { response: res, ok: res.ok, url },
        })
      );
    } catch (err) {
      this.dispatchEvent(
        new CustomEvent('smart-form:error', {
          detail: { error: err },
        })
      );
    }
  }
}

customElements.define('smart-form', SmartForm);

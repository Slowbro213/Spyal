export class SmartForm extends HTMLElement {
  constructor() {
    super();
  }

  connectedCallback() {
    this.addEventListener('submit', this.handleFormSubmit as EventListener);

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
          this.handleFormSubmit(e);
        }
      }
    });
  }

  handleFormSubmit(e: Event) {
    e.preventDefault();

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

    // Serialize FormData to an object (optional, but usually convenient)
    const data: Record<string, string> = {};
    for (const [key, value] of formData.entries()) {
      data[key] = typeof value === 'string' ? value : String(value);
    }

    // Dispatch a custom event with form data (cancelable, bubbles up)
    this.dispatchEvent(
      new CustomEvent('smart-form:submit', {
        detail: { formData, data },
        bubbles: true,
        cancelable: true,
      })
    );
  }
}

customElements.define('smart-form', SmartForm);

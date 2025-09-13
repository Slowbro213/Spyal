import { Config } from '@alspy/config';
import { Importance, initToast, Level } from '@alspy/services/toast';

const toast = initToast();

export class LoginButton extends HTMLElement {
  private shadow: ShadowRoot;
  private loggedIn = false;
  private buttonEl!: HTMLButtonElement;
  private overlayEl!: HTMLElement;
  private userInput!: HTMLInputElement;
  private passInput!: HTMLInputElement;
  private submitBtn!: HTMLButtonElement;

  constructor() {
    super();
    this.shadow = this.attachShadow({ mode: 'open' });

    this.shadow.innerHTML = `
      <style>
        :host { display: inline-block; }
        .login-btn {
          display: inline-flex;
          align-items: center;
          justify-content: center;
          padding: 0.5rem 1rem;
          border-radius: 0.75rem;
          font-weight: 600;
          background: var(--btn-bg, #34d399);
          color: var(--btn-fg, #000);
          cursor: pointer;
          border: none;
        }
        .login-btn[aria-disabled="true"] { opacity: 0.6; cursor: not-allowed; }
        /* modal */
        .overlay {
          display: none;               /* hidden by default */
          position: fixed;             /* fixed relative to viewport */
          top: 0;
          left: 0;
          width: 100vw;
          height: 100vh;               /* full viewport height */
          background: rgba(0, 0, 0);
          z-index: 9999;
        
        }
        .overlay.show {
        display: flex;
        align-items: center;         /* vertical centering */
        justify-content: center;     /* horizontal centering */
      
        padding: 1rem;
        box-sizing: border-box;
 }
        .modal {
          width: 100%;
          max-width: 420px;
          background: var(--card-bg, #0b1220);
          color: var(--text, #fff);
          border-radius: 12px;
          padding: 1.25rem;
          box-shadow: 0 10px 30px rgba(0,0,0,0.4);
        }
        .modal h3 { margin: 0 0 1rem 0; font-size: 1.125rem; }
        .field { display: flex; flex-direction: column; gap: 0.25rem; margin-bottom: 0.75rem; }
        input[type="text"], input[type="password"] {
          padding: 0.5rem 0.75rem;
          border-radius: 8px;
          border: 1px solid rgba(255,255,255,0.06);
          background: rgba(255,255,255,0.02);
          color: inherit;
          outline: none;
        }
        .actions { display:flex; gap:0.5rem; justify-content: flex-end; margin-top: 0.5rem; }
        .btn {
          padding: 0.5rem 0.75rem;
          border-radius: 8px;
          font-weight: 600;
          cursor: pointer;
          border: none;
        }
        .btn.secondary { background: transparent; color: var(--text-secondary, #ccc); }
        .btn.primary { background: var(--btn-bg, #34d399); color: var(--btn-fg, #000); }
        @media (max-width: 480px) {
          .modal { padding: 1rem; border-radius: 10px; }
        }
      </style>

      <button class="login-btn" part="button" type="button">Login</button>

      <div class="overlay" id="overlay" role="dialog" aria-modal="true" aria-hidden="true">
        <div class="modal" role="document">
          <h3>Login</h3>
          <div class="field">
            <label>Username</label>
            <input id="username" type="text" autocomplete="username" />
          </div>
          <div class="field">
            <label>Password</label>
            <input id="password" type="password" autocomplete="current-password" />
          </div>
          <div class="actions">
            <button class="btn secondary" id="cancel">Anulo</button>
            <button class="btn primary" id="confirm">Login</button>
          </div>
        </div>
      </div>
    `;
  }

  connectedCallback() {
    this.buttonEl = this.shadow.querySelector(
      '.login-btn'
    ) as HTMLButtonElement;
    this.overlayEl = this.shadow.getElementById('overlay') as HTMLElement;
    this.userInput = this.shadow.getElementById('username') as HTMLInputElement;
    this.passInput = this.shadow.getElementById('password') as HTMLInputElement;
    this.submitBtn = this.shadow.getElementById('confirm') as HTMLButtonElement;

    this.buttonEl.addEventListener('click', this.onButtonClick.bind(this));
    this.shadow
      .getElementById('cancel')!
      .addEventListener('click', this.closeModal.bind(this));
    this.overlayEl.addEventListener('click', (e) => {
      if (e.target === this.overlayEl) this.closeModal();
    });

    this.submitBtn.addEventListener('click', this.onSubmit.bind(this));
    this.shadow.addEventListener('keydown', this.onKeyDown.bind(this));

    // initialize state from storage
    const prefix = Config.TOKEN;
    const storedUser = localStorage.getItem(`${prefix}username`);
    const storedToken = localStorage.getItem(`${prefix}token`);
    if (storedUser && storedToken) {
      this.setLoggedIn(storedUser);
    }
  }

  disconnectedCallback() {
    this.buttonEl.removeEventListener('click', this.onButtonClick.bind(this));
    this.submitBtn.removeEventListener('click', this.onSubmit.bind(this));
  }

  private onButtonClick(e: MouseEvent) {
    e.preventDefault();
    if (this.loggedIn) return;
    this.openModal();
  }

  private openModal() {
    this.overlayEl.classList.add('show');
    this.overlayEl.setAttribute('aria-hidden', 'false');
    this.userInput.value = '';
    this.passInput.value = '';
    this.userInput.focus();
    document.body.style.overflow = 'hidden';
  }

  private closeModal() {
    this.overlayEl.classList.remove('show');
    this.overlayEl.setAttribute('aria-hidden', 'true');
    document.body.style.overflow = 'auto';
  }

  private setLoggedIn(name: string) {
    this.loggedIn = true;
    this.buttonEl.textContent = name;
    this.buttonEl.setAttribute('aria-disabled', 'true');
  }

  private async onSubmit(e?: Event) {
    e?.preventDefault();
    const username = this.userInput.value.trim();
    const password = this.passInput.value;

    if (!username || !password) {
      toast.show(Level.Error, Importance.Major, {
        message: 'Plotësoni të gjitha fushat e detyrueshme!',
      });
      return;
    }

    this.submitBtn.disabled = true;
    this.submitBtn.textContent = 'Duke punuar...';

    try {
      const body = new URLSearchParams();
      body.set('username', username);
      body.set('password', password);

      const res = await fetch('/login', {
        method: 'POST',
        body: body,
        headers: {
          Accept: 'application/json',
        },
      });

      if (!res.ok) {
        const txt = await res.text().catch(() => '');
        toast.show(Level.Error, Importance.Major, {
          message: txt || 'Gabim gjatë autentikimit',
        });
        this.submitBtn.disabled = false;
        this.submitBtn.textContent = 'Login';
        return;
      }

      const data = await res.json().catch(() => null);
      const token = data?.token;
      if (!token) {
        toast.show(Level.Error, Importance.Major, {
          message: 'Token jo i vlefshëm nga serveri',
        });
        this.submitBtn.disabled = false;
        this.submitBtn.textContent = 'Login';
        return;
      }

      // Store token + username
      const prefix = Config.TOKEN;
      localStorage.setItem(`${prefix}username`, username);

      this.setLoggedIn(username);
      this.closeModal();
    } catch (err) {
      toast.show(Level.Error, Importance.Major, {
        message: 'Gabim rrjeti — provoni përsëri',
      });
    } finally {
      this.submitBtn.disabled = false;
      this.submitBtn.textContent = 'Login';
    }
  }

  private onKeyDown(evt: Event) {
    const e = evt as KeyboardEvent;
    if (e.key === 'Escape') {
      this.closeModal();
    } else if (e.key === 'Enter' && this.overlayEl.classList.contains('show')) {
      // submit on Enter
      this.onSubmit();
    }
  }
}

customElements.define('login-button', LoginButton);

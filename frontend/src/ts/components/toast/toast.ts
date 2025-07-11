import { Level, Importance, Message } from '@alspy/services/toast/types';

export class Toast extends HTMLElement {
  constructor() {
    super();
    this.setAttribute('id', 'toast-service');
    this.attachShadow({ mode: 'open' });
  }

  show(
    severity: Level,
    importance: Importance,
    message: Message,
    duration: number = 3000
  ) {
    const container = document.createElement('div');
    const isMajor = importance === Importance.Major;
    const position = isMajor ? 'top' : 'bottom';
    const colors = [
      '#4CAF50', // Success (green)
      '#2196F3', // Info (blue)
      '#FFEB3B', // Warning (yellow)
      '#F44336', // Error (red)
    ];

    container.innerHTML = `
      <style>
        @keyframes fadeIn {
          from {
            opacity: 0;
            transform: translateY(${isMajor ? '-10px' : '10px'});
          }
          to {
            opacity: 1;
            transform: translateY(0);
          }
        }

        @keyframes fadeOut {
          from {
            opacity: 1;
            transform: translateY(0);
          }
          to {
            opacity: 0;
            transform: translateY(${isMajor ? '-10px' : '10px'});
          }
        }

        @keyframes progress {
          from { width: 100%; }
          to { width: 0%; }
        }

        .toast {
          z-index: 9999;
          position: fixed;
          ${position}: 1.5rem;
          right: 1.5rem;
          padding: 0.75rem 1rem;
          background: #2c2c2c;
          border-left: 3px solid ${colors[severity]};
          color: white;
          border-radius: 4px;
          box-shadow: 0 3px 8px rgba(0,0,0,0.3);
          opacity: 0;
          transform: translateY(${isMajor ? '-10px' : '10px'});
          animation: 
            fadeIn 0.3s ease-out forwards, 
            fadeOut 0.3s ease-in forwards ${duration - 300}ms;
          width: 16rem;
          display: flex;
          align-items: center;
          font-family: 'Segoe UI', 'Roboto', 'Oxygen', 'Ubuntu', sans-serif;
          word-wrap: break-word;
          overflow-wrap: break-word;
        }

        .content {
          display: flex;
          flex-direction: column;
          gap: 0.4rem;
          width: 100%;
          position: relative;
          z-index: 2;
        }

        .title {
          font-weight: 600;
          font-size: 1rem;
          display: flex;
          align-items: center;
          gap: 0.5rem;
          letter-spacing: 0.4px;
          word-wrap: break-word;
          overflow-wrap: break-word;
        }

        .message {
          font-size: 0.85rem;
          line-height: 1.3;
          word-wrap: break-word;
          overflow-wrap: break-word;
        }

        .importance {
          font-size: 0.7rem;
          margin-top: 0.2rem;
          opacity: 0.9;
          font-weight: 500;
          display: inline-block;
          padding: 0.2rem 0.4rem;
          background: rgba(255,255,255,0.1);
          border-radius: 0.2rem;
          align-self: flex-start;
        }

        .icon {
          width: 1.2rem;
          height: 1.2rem;
          display: inline-flex;
          color: ${colors[severity]};
        }

        .progress-bar {
          position: absolute;
          bottom: 0;
          left: 0;
          height: 2px;
          background: ${colors[severity]};
          animation: progress ${duration}ms linear forwards;
        }

        .close-button {
          position: absolute;
          top: 0.4rem;
          right: 0.4rem;
          background: none;
          border: none;
          color: white;
          font-size: 1rem;
          cursor: pointer;
          opacity: 0.7;
          transition: opacity 0.2s, color 0.2s;
          line-height: 1;
          padding: 0.2rem;
        }

        .close-button:hover {
          opacity: 1;
          color: ${colors[severity]};
        }

        /* Responsive design */
        @media (max-width: 600px) {
          .toast {
            right: 0.75rem;
            left: 0.75rem;
            max-width: 80%;
            margin: 0 auto;
            padding: 0.6rem 0.8rem;
          }
          .title {
            font-size: 0.9rem;
          }
          .message {
            font-size: 0.8rem;
          }
          .importance {
            font-size: 0.65rem;
          }
          .icon {
            width: 1rem;
            height: 1rem;
          }
          .close-button {
            font-size: 0.9rem;
            top: 0.3rem;
            right: 0.3rem;
          }
        }

        @media (max-width: 400px) {
          .toast {
            ${position}: 1rem;
            max-width: 90%;
            padding: 0.5rem 0.7rem;
          }
          .title {
            font-size: 0.85rem;
          }
          .message {
            font-size: 0.75rem;
          }
          .importance {
            font-size: 0.6rem;
            padding: 0.15rem 0.3rem;
          }
          .icon {
            width: 0.9rem;
            height: 0.9rem;
          }
          .close-button {
            font-size: 0.85rem;
            top: 0.25rem;
            right: 0.25rem;
          }
        }
      </style>
      <div class="toast" role="alert">
        <button class="close-button" aria-label="Close">×</button>
        <div class="progress-bar"></div>
        <div class="content">
          <div class="title">
            ${this.getIcon(severity)}
            <span>${message.title || Level[severity]}</span>
          </div>
          ${message.message ? `<div class="message">${message.message}</div>` : ''}
        </div>
      </div>
    `;

    if (this.shadowRoot === null) return;
    this.shadowRoot.innerHTML = '';
    this.shadowRoot.appendChild(container);

    container.querySelector('.toast')?.addEventListener('click', () => {
      this.close();
    });

    // Keep the close button handler for accessibility
    container.querySelector('.close-button')?.addEventListener('click', (e) => {
      e.stopPropagation();
      this.close();
    });

    setTimeout(() => {
      this.close();
    }, duration);
  }

  private getIcon(severity: Level): string {
    const icons = [
      // Success (checkmark in circle)
      `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
        <circle cx="12" cy="12" r="10" stroke-width="1.5" fill="currentColor" fill-opacity="0.1"/>
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M9 12l2 2 4-4" />
      </svg>`,
      // Info (i in circle)
      `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
        <circle cx="12" cy="12" r="10" stroke-width="1.5" fill="currentColor" fill-opacity="0.1"/>
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 16v-4m0-4h.01" />
      </svg>`,
      // Warning (exclamation in triangle)
      `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
        <path fill="currentColor" fill-opacity="0.1" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
      </svg>`,
      // Error (x in circle)
      `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
        <circle cx="12" cy="12" r="10" stroke-width="1.5" fill="currentColor" fill-opacity="0.1"/>
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M6 18L18 6M6 6l12 12" />
      </svg>`,
    ];
    return icons[severity];
  }

  close() {
    if (this.shadowRoot === null) return;
    this.shadowRoot.innerHTML = '';
  }
}

customElements.define('toast-service', Toast);

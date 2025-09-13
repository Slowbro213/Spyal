import { initToast } from '@alspy/services';
import { Importance, Level } from '@alspy/services/toast';
import { Toast } from '@alspy/components/toast';
import { Config } from '@alspy/config';
import { Staging } from '@alspy/config/types';

export enum Severity {
  trivial,
  small,
  normal,
  huge,
  catastrofic,
}

const severityToEmoji: Record<Severity, string> = {
  [Severity.trivial]: '😅',
  [Severity.small]: '🙄',
  [Severity.normal]: '😱',
  [Severity.huge]: '🤯',
  [Severity.catastrofic]: '💥',
};

const severityToColor: Record<Severity, string> = {
  [Severity.trivial]: '#4ade80', // green
  [Severity.small]: '#60a5fa', // blue
  [Severity.normal]: '#f59e0b', // amber
  [Severity.huge]: '#ef4444', // red
  [Severity.catastrofic]: '#7e22ce', // purple
};

const toast: Toast = initToast();

export const serveErrorPage = async (
  severity: Severity = Severity.normal,
  customTitle = 'Gabim i Madh!',
  customMessage = 'Ndodhi një gabim i papritur. Ju lutemi rifilloni lojën.'
) => {
  if (Config.STAGE === Staging.Production && severity < Severity.normal) return;

  if (severity <= Severity.small) {
    toast.show(Level.Error, Importance.Minor, {
      title: 'Small Error Occurred',
      message: 'A small error has occurred, check logs',
    });
  }

  try {
    const res = await fetch('/public/html/error.html');
    let errorHtml = await res.text();

    errorHtml = errorHtml
      .replace(/{{\s*title[^}]*}}/g, customTitle)
      .replace(/{{\s*message[^}]*}}/g, customMessage)
      .replace(/{{\s*emoji[^}]*}}/g, severityToEmoji[severity])
      .replace(/{{\s*color[^}]*}}/g, severityToColor[severity]);

    document.body.innerHTML = errorHtml;
  } catch (error) {
    console.error(error);
    document.body.innerHTML = `
      <div style="text-align:center;padding:2rem">
        <h1>${customTitle}</h1>
        <p>${customMessage}</p>
      </div>
    `;
  }
};

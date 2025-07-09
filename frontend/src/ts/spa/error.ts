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

// Map each severity to its own file
const severityToFile: Record<Severity, string> = {
  [Severity.trivial]: 'public/html/error-trivial.html',
  [Severity.small]: 'public/html/error-small.html',
  [Severity.normal]: 'public/html/error.html',
  [Severity.huge]: 'public/html/error-huge.html',
  [Severity.catastrofic]: 'public/html/error-catastrofic.html',
};

const toast: Toast = initToast();

export const serveErrorPage = async (severity: Severity = Severity.normal) => {
  if (Config.STAGE === Staging.Production && severity < Severity.normal) return;

  if (Severity.small >= severity)
    toast.show(Level.Error, Importance.Minor, {
      title: 'Small Error Occured',
      message: 'A small error has occured, check logs',
    });

  const file = severityToFile[severity] || severityToFile[Severity.normal];
  const res = await fetch(file);
  const errorHtml = await res.text();
  document.body.innerHTML = errorHtml;
};

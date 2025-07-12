import { Config, Staging } from '@alspy/config';
import { serveErrorPage, Severity } from './error';

type Log = {
  level: 'info' | 'warn' | 'error';
  msg: string;
};

export const log = async (l: Log) => {
  const prod = Config.STAGE === Staging.Production;
  if (!prod) console.log(l);

  const url = Config.LOG_ROUTE;

  try {
    const res = await fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(l),
    });

    if (!res.ok) {
      if (!prod)
        serveErrorPage(
          Severity.huge,
          'Logging route error',
          'The logging route failed to be reached ' + (await res.text())
        );
      console.error(`Failed to log to ${url}`, await res.text());
    }
  } catch (err) {
    if (!prod)
      serveErrorPage(
        Severity.huge,
        'Logging route error',
        `The logging route ${url} failed to be reached ${err}`
      );
    console.error('Log request failed:', err);
  }
};

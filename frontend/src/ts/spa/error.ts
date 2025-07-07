type Severity = 'trivial' | 'small' | 'normal' | 'huge' | 'catastrofic';

// Map each severity to its own file
const severityToFile: Record<Severity, string> = {
  trivial: 'public/html/error-trivial.html',
  small: 'public/html/error-small.html',
  normal: 'public/html/error.html', // Default
  huge: 'public/html/error-huge.html',
  catastrofic: 'public/html/error-catastrofic.html',
};

export const serveError = async (severity: Severity = 'normal') => {
  const file = severityToFile[severity] || severityToFile.normal;
  const res = await fetch(file);
  const errorHtml = await res.text();
  document.body.innerHTML = errorHtml;
};

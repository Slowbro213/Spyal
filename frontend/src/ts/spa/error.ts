export const serveError = async () => {
  const res = await fetch('public/html/error.html');
  const errorHtml = await res.text();
  document.body.innerHTML = errorHtml;
};

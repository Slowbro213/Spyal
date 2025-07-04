// Select all elements where id starts with "component/"
const components = document.querySelectorAll<HTMLElement>('[id^="component/"]');

components.forEach((el) => {
  const componentPath = el.id; // e.g., "component/room"

  fetch(`/${componentPath}`) // becomes "/component/room"
    .then(res => {
      if (!res.ok) throw new Error(`Failed to fetch ${componentPath}`);
      return res.text();
    })
    .then(html => {
      el.outerHTML = html;
    })
    .catch(err => {
      el.textContent = `❌ Error loading ${componentPath}`;
      console.error(err);
    });
});

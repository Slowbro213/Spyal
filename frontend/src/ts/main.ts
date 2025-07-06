import { componentMap } from './components';
import './components';
import { pageCache } from './pageCache';

const renderComponents = async () => {
  const components = document.querySelectorAll<HTMLElement>(
    '[id^="components/"]'
  );
  if (components.length === 0) return;

  let error = false;

  await Promise.all(
    Array.from(components).map(async (el) => {
      const componentPath = el.id;
      const component = componentPath.slice('/components/'.length - 1);

      const params = new URLSearchParams(
        Object.entries(el.dataset).reduce(
          (acc, [key, val]) => {
            if (val !== undefined) acc[key] = val;
            return acc;
          },
          {} as Record<string, string>
        )
      ).toString();

      try {
        const res = await fetch(`/${componentPath}?${params}`);
        if (!res.ok) throw new Error(`Failed to fetch ${componentPath}`);

        const html = await res.text();

        const temp = document.createElement('div');
        temp.innerHTML = html;
        const newNode = temp.firstElementChild;

        if (newNode) {
          el.replaceWith(newNode); // replace element safely
          componentMap.get(component)?.(); // run component logic
        }
      } catch (err) {
        el.textContent = `❌ Error loading ${componentPath}`;
        console.error(err);
        error = true;
      }
    })
  );

  if (!error) {
    await renderComponents(); // Recursively render any nested components
  }
};

renderComponents();

window.addEventListener('popstate', async () => {
  const href = location.pathname;
  const target = document.querySelector('#content');
  if (!target) return;

  if (pageCache.has(href)) {
    target.innerHTML = pageCache.get(href)!;
  } else {
    const res = await fetch(href, { headers: { 'X-Smart-Link': 'true' } });
    const html = await res.text();
    pageCache.set(href, html);
    target.innerHTML = html;
  }
});

const prefetchRoutes = ['/', '/create']; // your important routes

async function prefetchPages() {
  for (const route of prefetchRoutes) {
    try {
      if (pageCache.has(route)) continue;
      const res = await fetch(route, {
        headers: { 'X-Smart-Link': 'true' },
      });
      if (!res.ok) throw new Error(`Failed to prefetch ${route}`);
      const html = await res.text();
      pageCache.set(route, html);
    } catch (err) {
      console.warn('Prefetch error for', route, err);
    }
  }
}

document.addEventListener('DOMContentLoaded', () => {
  const path = location.pathname;
  const target = document.querySelector('#content');
  if (target && !pageCache.has(path)) {
    pageCache.set(path, target.innerHTML);
  }
  prefetchPages();
});

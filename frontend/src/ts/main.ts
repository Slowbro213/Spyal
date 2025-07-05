import { componentMap } from './components';

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

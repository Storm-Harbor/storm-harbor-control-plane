(() => {
  const storageKey = "dr-docs-primary-navigation-collapsed";
  const desktop = window.matchMedia("(min-width: 76.25em)");
  let preference = null;
  try {
    const stored = localStorage.getItem(storageKey);
    if (stored === "true" || stored === "false") preference = stored === "true";
  } catch (_) { /* Navigation remains usable without storage. */ }

  const elements = () => ({
    button: document.querySelector("[data-dr-nav-toggle]"),
    primary: document.querySelector(".md-sidebar--primary"),
    drawer: document.querySelector("#__drawer"),
  });

  const apply = () => {
    const {button, primary, drawer} = elements();
    if (!button || !primary || !drawer) return;
    if (!primary.dataset.drDefaultCollapsed) {
      primary.dataset.drDefaultCollapsed = String(primary.hidden);
    }
    // Material hides navigation on the landing page. Keep that initial layout,
    // but let the same menu reveal it instead of leaving the hidden attribute.
    primary.hidden = false;
    primary.id = "dr-primary-navigation";
    const collapsed = preference ?? primary.dataset.drDefaultCollapsed === "true";
    const expanded = desktop.matches ? !collapsed : drawer.checked;
    document.body.classList.toggle("dr-nav-collapsed", desktop.matches && collapsed);
    primary.inert = !expanded;
    button.setAttribute("aria-expanded", String(expanded));
    button.setAttribute("aria-label", expanded ? "Close navigation" : "Open navigation");
    button.title = expanded ? "Close navigation" : "Open navigation";
  };

  const setDrawer = (open) => {
    const {drawer} = elements();
    if (!drawer) return;
    drawer.checked = open;
    drawer.dispatchEvent(new Event("change", {bubbles: true}));
  };

  const save = (collapsed) => {
    preference = collapsed;
    try { localStorage.setItem(storageKey, String(collapsed)); } catch (_) { /* Optional storage. */ }
  };

  document.addEventListener("click", (event) => {
    if (!event.target.closest("[data-dr-nav-toggle]")) return;
    const {button, drawer} = elements();
    if (desktop.matches) save(button.getAttribute("aria-expanded") === "true");
    else setDrawer(!drawer.checked);
    apply();
  });

  document.addEventListener("change", (event) => {
    if (event.target.id === "__drawer") apply();
  });

  document.addEventListener("keydown", (event) => {
    const {button} = elements();
    if (event.key !== "Escape" || !button || button.getAttribute("aria-expanded") !== "true") return;
    if (desktop.matches) save(true);
    else setDrawer(false);
    apply();
    button.focus();
  });

  desktop.addEventListener("change", () => {
    setDrawer(false);
    apply();
  });

  const install = () => {
    setDrawer(false);
    apply();
  };
  if (document.readyState === "loading") document.addEventListener("DOMContentLoaded", install, {once: true});
  else install();
  if (typeof document$ !== "undefined") document$.subscribe(install);
})();

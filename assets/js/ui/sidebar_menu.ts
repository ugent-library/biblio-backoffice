const LOCAL_STORAGE_KEY = "detail-sidebar-state";
const DEFAULT_SIDEBAR_STATE = "sidebar-message";

export default function initSidebarMenu(el) {
  const menuButtons = el.querySelectorAll("[data-sidebar-menu] button");
  const closeButtons = el.querySelectorAll("[data-sidebar-close]");

  menuButtons.forEach((button) => {
    button.addEventListener("click", () => {
      const targetId = button.getAttribute("data-target-id");
      const targetContent = document.getElementById(targetId);

      // Check if the target content is currently displayed.
      if (targetContent.classList.contains("open")) {
        // If so, hide it.
        targetContent.classList.remove("open");
        button.classList.remove("active");

        setActiveSidebar(null);
      } else {
        // If not, hide all contents and show the target content.
        document
          .querySelectorAll("[data-sidebar-content]")
          .forEach((c) => c.classList.remove("open"));
        menuButtons.forEach((btn) => btn.classList.remove("active"));

        targetContent.classList.add("open");
        button.classList.add("active");

        setActiveSidebar(targetId);
      }
    });
  });

  closeButtons.forEach((button) => {
    button.addEventListener("click", () => {
      const parentContent = button.closest("[data-sidebar-content]");
      if (parentContent) {
        parentContent.classList.remove("open");

        // Find the corresponding menu button and remove active class.
        getSidebarNavLink(parentContent.id)?.classList.remove("active");

        setActiveSidebar(null);
      }
    });
  });

  // Set initial active sidebar
  const activeSidebar = getActiveSidebar();
  if (activeSidebar) {
    document.getElementById(activeSidebar)?.classList.add("open");
    getSidebarNavLink(activeSidebar)?.classList.add("active");
  }
}

function getSidebarNavLink(targetId: string): HTMLAnchorElement | null {
  return document.querySelector(
    `[data-sidebar-menu] .nav-link[data-target-id="${targetId}"]`,
  );
}

function getActiveSidebar(): string | undefined {
  const key = getPageKey();
  if (key) {
    const state = getDetailSidebarState();

    if (state.has(key)) {
      return state.get(key);
    } else {
      return DEFAULT_SIDEBAR_STATE;
    }
  }
}

function setActiveSidebar(activeSidebar: string | null) {
  const key = getPageKey();
  if (key) {
    const state = getDetailSidebarState();
    if (activeSidebar !== DEFAULT_SIDEBAR_STATE) {
      state.set(key, activeSidebar);
    } else {
      state.delete(key);
    }

    localStorage.setItem(
      LOCAL_STORAGE_KEY,
      JSON.stringify(Object.fromEntries(state)),
    );
  }
}

let sidebarState: Map<string, string> = null;

function getDetailSidebarState(): Map<string, string> {
  if (!sidebarState) {
    const savedState = localStorage.getItem(LOCAL_STORAGE_KEY);
    if (savedState) {
      try {
        const state = JSON.parse(savedState) as Record<string, string>;

        sidebarState = new Map(Object.entries(state));
      } catch (err) {
        console.warn(
          `There was an error parsing the saved state with key "${LOCAL_STORAGE_KEY}": `,
          err,
        );

        sidebarState = new Map();
      }
    } else {
      sidebarState = new Map();
    }
  }

  return sidebarState;
}

function getPageKey() {
  const regex = /^\/(publication|dataset)\/(?<id>[A-Z0-9]{10,})$/;

  const match = document.location.pathname.match(regex);

  return match?.groups?.id;
}

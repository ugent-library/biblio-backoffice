/**
 * Collapsible header.
 *
 * Collapsible header is used for hiding / showing the summary of a
 * publication on the edit forms. The state of the header (hidden / shown)
 * is persisted across page reloads via the browsers local storage.
 */
export default function () {
  class Header {
    constructor(el) {
      this.el = el;
      this.collapsers = el.querySelectorAll('.c-header-collapse-trigger');
      this.collapsed = false;
      this.init();
    }

    init() {
      try {
        const savedState = JSON.parse(
          localStorage.getItem(HEADER_STATE_STORAGE_KEY)
        );
        headerState = Object.assign({}, headerState, savedState);
      } catch (err) {
        console.warn(
          'There was an error parsing the saved state for the prototype navigation.'
        );
      }

      // Handle state on page load: open/close nav and close saved modules
      if (headerState.isExpanded) {
        this.show();
      } else {
        this.hide();
      }

      this.collapsers?.forEach((collapser) =>
        collapser.addEventListener('click', this)
      );
    }

    handleEvent(event) {
      this.toggle(event);
    }

    toggle(event) {
      if (this.collapsed) this.show();
      else if (!this.collapsed) this.hide();
    }

    show() {
      this.el.dataset.collapsed = false;
      this.collapsed = false;
      headerState.isExpanded = true;

      this.saveNavState();
    }

    saveNavState() {
      localStorage.setItem(
        HEADER_STATE_STORAGE_KEY,
        JSON.stringify(headerState)
      );
    }

    hide() {
      this.el.dataset.collapsed = true;
      this.collapsed = true;
      headerState.isExpanded = false;
      this.saveNavState();
    }
  }

  let headerState = {
    isExpanded: false,
  };

  const collapsibleHeaders = document.querySelectorAll('.c-header-collapsible');
  const HEADER_STATE_STORAGE_KEY = `bedrockheaderState`;

  if (collapsibleHeaders.length) [...collapsibleHeaders].map((collapsibleHeader) => new Header(collapsibleHeader));
}
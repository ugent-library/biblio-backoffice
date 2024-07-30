import htmx from "htmx.org/dist/htmx.esm.js";
import Popover from "bootstrap.native/popover";

export default function () {
  htmx.onLoad((rootEl) => {
    rootEl.querySelectorAll("[data-bs-toggle=popover-custom]").forEach((el) => {
      const container = document.querySelector(el.dataset.popoverContent);
      const content = container.querySelector(".popover-body");
      const heading = container.querySelector(".popover-heading");
      const title = heading?.innerHTML ?? "";

      new Popover(el, {
        content: content.innerHTML,
        title,
        delay: 1000,
      });
    });
  });
}

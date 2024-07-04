import htmx from "htmx.org/dist/htmx.esm.js";
import { Popover } from "bootstrap.native";

export default function () {
  let addEvents = function (rootEl) {
    rootEl
      .querySelectorAll("[data-bs-toggle=popover-custom]")
      .forEach(function (el) {
        let container = document.querySelector(el.dataset.popoverContent);
        let content = container.querySelector(".popover-body");
        let heading = container.querySelector(".popover-heading");
        let title = "";
        if (heading) {
          title = heading.innerHTML;
        }
        new Popover(el, {
          content: content.innerHTML,
          title: title,
          delay: 1000,
        });
      });
  };

  htmx.onLoad(addEvents);
}

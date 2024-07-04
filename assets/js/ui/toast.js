import htmx from "htmx.org/dist/htmx.esm.js";

export default function () {
  let addEvents = function (rootEl) {
    rootEl
      .querySelectorAll('.btn-close[data-bs-dismiss="toast"]')
      .forEach(function (btn) {
        btn.addEventListener("click", () => {
          btn.closest(".toast").remove();
        });
      });
  };

  htmx.onLoad(addEvents);
}

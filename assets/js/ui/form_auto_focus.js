import htmx from "htmx.org/dist/htmx.esm.js";

export default function () {
  /**
   * Due to the use of HTMX and the lack of element form
   * this is now the only way to recognize a focusable element
   * on the page
   * Form now only present on htmx load
   */
  htmx.onLoad(function (el) {
    let input = el.querySelector(".form-control-auto-focus");
    if (input == null) return;
    input.focus();
  });
}

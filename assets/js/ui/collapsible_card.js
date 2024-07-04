// /* ==========================================================================
//     Collapsible card content
//    ========================================================================== */

import htmx from "htmx.org/dist/htmx.esm.js";

export default function () {
  let toggler = function (evt) {
    let showCardInfo = evt.currentTarget.closest("[data-collapsible-card]");
    let cardContent = showCardInfo.querySelector(
      "[data-collapsible-card-content]",
    );
    cardContent.classList.toggle("show");
  };

  let addEvents = function () {
    document
      .querySelectorAll("[data-collapsible-card-toggle]")
      .forEach((cardContent) => cardContent.addEventListener("click", toggler));
  };

  addEvents();

  htmx.onLoad(function (el) {
    addEvents();
  });
}

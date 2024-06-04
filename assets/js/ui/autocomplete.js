import htmx from "htmx.org";

export default function () {
  let addEvents = function (rootEl) {
    rootEl.querySelectorAll(".autocomplete").forEach(function (ac) {
      let target = document.querySelector(ac.dataset.target);
      target.addEventListener("keydown", function (evt) {
        if (evt.key === "Escape") {
          ac.querySelector(".autocomplete-hits").innerHTML = "";
        }
      });
    });

    rootEl.querySelectorAll(".autocomplete-hit").forEach(function (hit) {
      hit.addEventListener("click", function () {
        let target = hit.closest(".autocomplete").dataset.target;
        let v = hit.dataset.value;
        document.querySelector(target).value = v;
        hit.closest(".autocomplete-hits").innerHTML = "";
      });
    });
  };

  htmx.onLoad(function (el) {
    addEvents(el);
  });
  // listen for inputs added with multiple.js
  document.addEventListener("form-value-add", function (evt) {
    addEvents(evt.target);
  });
}

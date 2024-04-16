export default function () {
  document.querySelectorAll("[data-facet-dropdown]").forEach((dropdown) => {
    const facet = dropdown.getAttribute("data-facet-dropdown");
    const form = dropdown.querySelector("form");

    if (form) {
      let beginState = null;

      dropdown.addEventListener("show.bs.dropdown", function (evt) {
        beginState = getCurrentState(form, facet);
      });

      dropdown.addEventListener("hidden.bs.dropdown", function (evt) {
        const currentState = getCurrentState(form, facet);

        if (beginState !== null && beginState !== currentState) {
          form.submit();
        }
      });
    } else {
      console.log(`Could not find form element for facet '${facet}'.`);
    }
  });
}

function getCurrentState(form, facet) {
  const entries = Array.from(new FormData(form).entries());

  return entries
    .filter(([key]) => key === `f[${facet}]`)
    .map(([_, value]) => value)
    .join(",");
}

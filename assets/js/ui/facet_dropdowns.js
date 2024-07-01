export default function () {
  document.querySelectorAll("[data-facet-dropdown]").forEach((dropdown) => {
    const facet = dropdown.getAttribute("data-facet-dropdown");
    const form = dropdown.querySelector("form");

    if (form) {
      let beginState = null;

      dropdown.addEventListener("show.bs.dropdown", function () {
        beginState = getCurrentState(form, facet);
      });

      dropdown.addEventListener("shown.bs.dropdown", function () {
        // Auto-focus first input field if it exists
        form.querySelector("input[type=text]")?.focus();
      });

      dropdown.addEventListener("hidden.bs.dropdown", function () {
        const currentState = getCurrentState(form, facet);

        // Apply filter when dropdown fields have been altered
        if (beginState !== null && beginState !== currentState) {
          form.submit();
        }

        beginState = null;
      });
    } else {
      console.log(`Could not find form element for facet '${facet}'.`);
    }
  });
}

function getCurrentState(form, facet) {
  // TODO: test with: new URLSearchParams(new FormData(form)).toString()
  const entries = Array.from(new FormData(form).entries());

  return entries
    .filter(([key]) => key === `f[${facet}]`)
    .map(([_, value]) => value)
    .join(",");
}

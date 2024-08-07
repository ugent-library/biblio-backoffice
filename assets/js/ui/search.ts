export default function initSearchFields(el: HTMLElement) {
  el.querySelectorAll<HTMLInputElement>("input[type=search]").forEach(
    (input) => {
      const form = input.closest("form");
      if (form && !hasHtmxAttributes(input) && !hasHtmxAttributes(form)) {
        input.addEventListener("search", () => {
          if (input.value === "") {
            form.submit();
          }
        });
      }
    },
  );
}

function hasHtmxAttributes(el: HTMLElement) {
  return el.getAttributeNames().some((attr) => attr.startsWith("hx-"));
}

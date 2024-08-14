export default function initSearchFields(el: HTMLElement) {
  el.querySelectorAll<HTMLInputElement>("input[data-submit-on-clear]").forEach(
    (input) => {
      const form = input.closest("form");
      if (form) {
        input.addEventListener("search", () => {
          if (input.value === "") {
            form.submit();
          }
        });
      }
    },
  );
}

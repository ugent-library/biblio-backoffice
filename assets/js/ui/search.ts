export default function initSearchFields(el: HTMLElement) {
  el.querySelectorAll<HTMLInputElement>("input[data-submit-on-clear]").forEach(
    (input) => {
      const form = input.closest("form");
      if (form) {
        input.dataset.previousValue = input.value;

        input.addEventListener("change", () => {
          input.dataset.previousValue = input.value;
        });

        input.addEventListener("input", () => {
          if (input.dataset.previousValue && input.value === "") {
            form.submit();
          }
        });
      }
    },
  );
}

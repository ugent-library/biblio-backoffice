export default function (rootEl) {
  rootEl.querySelectorAll("[data-clipboard]").forEach((btn) => {
    btn.addEventListener("click", () => {
      navigator.clipboard.writeText(btn.dataset.clipboard).then(() => {
        let icon = btn.querySelector(".if");
        let text = btn.querySelector(".btn-text");
        let origBtnClass = btn.className;
        let origIconClass = icon.className;
        let origTextClass = text.className;
        let origText = text.innerText;

        btn.classList.remove("btn-outline-secondary");
        btn.classList.add("btn-outline-success");

        icon.classList.remove("if-copy", "text-muted");
        icon.classList.add("if-check", "text-success");

        text.classList.remove("text-muted");
        text.classList.add("text-success");
        text.setAttribute("aria-live", "polite");
        text.innerText = "Copied";

        setTimeout(function () {
          btn.className = origBtnClass;
          icon.className = origIconClass;
          text.className = origTextClass;
          text.innerText = origText;
        }, 1500);
      });
    });
  });
}

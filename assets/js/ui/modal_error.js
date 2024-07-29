import Modal from "bootstrap.native/modal";

export default function (error) {
  let modals = document.querySelector("#modals");
  if (modals) {
    /*
     * Expects somewhere in the document ..
     *
     * <template class="template-modal-error"></template>
     *
     * .. a template that encapsulates the modal body
     * */
    const templateModalError = document.querySelector(
      "template.template-modal-error",
    );

    if (templateModalError) {
      const modalEl = templateModalError.content
        .cloneNode(true)
        .querySelector(".modal");

      modalEl.querySelector(".msg").textContent = error;

      modals.innerHTML = "";
      modals.appendChild(modalEl);

      initModal(modalEl);
    }
  }
}

function initModal(modalEl) {
  const modal = new Modal(modalEl, {
    backdrop: "static",
    keyboard: false,
  });

  modal.show();

  modalEl.addEventListener("hidden.bs.modal", function () {
    modalEl.remove();
  });
}

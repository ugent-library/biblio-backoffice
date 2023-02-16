import htmx from 'htmx.org';

export default function (error) {

  let modals = document.querySelector("#modals")

  if (!modals) return

  /*
   * Expects somewhere in the document ..
   *
   * <template class="template-modal-error"></template>
   *
   * .. a template that encapsulates the modal body
   * */
  let templateModalError = document.querySelector("template.template-modal-error")

  if (!templateModalError) return

  let modal = templateModalError.content.cloneNode(true)

  // modal-close not triggered for dynamically added modals
  modal.querySelector(".modal-close").addEventListener("click", function(){
    modals.innerHTML = ""
  })

  modal.querySelector(".msg").textContent = error

  modals.innerHTML = ""
  modals.appendChild(modal)
}

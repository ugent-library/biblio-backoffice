// Handles fields with multiple values
export default function() {
    let multiple = function() {

        // Delete a value from the field
        let deleteFormValue = function (e) {
            let formField = e.target.closest("div.form-values")
            e.target.closest("div.form-value").remove();
            let length = Array.from(formField.children).length

            for (var i = 0; i < length; i++) {
                let input = formField.children[i].querySelector(".form-control")
                let name = input.getAttribute("name")
                name = name.replace(/\[.*\]/, "")
                input.setAttribute("name", name + "[" + i + "]")
            }

        }

        // Add a new value to the field
        let addFormValue = function (e) {
            let formField = e.target.closest("div.form-values")
            let formValue = formField.lastElementChild

            let length = Array.from(formField.children).length

            let node = formValue.cloneNode(true)

            formValue.querySelector(".form-control")
                    .value = ""

            let input = formValue.querySelector(".form-control")
            let inputName = input.getAttribute("name")
            inputName = inputName.replace(/\[.*\]/, "")
            input.setAttribute("name", inputName + "[" + length + "]")

            let classList = node.querySelector("button.form-value-add").classList
            classList.remove("form-value-add")
            classList.remove("btn-outline-primary")
            classList.add("btn-link-muted")
            classList.add("form-value-delete")

            classList = node.querySelector("i.if-add").classList
            classList.remove("if-add")
            classList.add("if-delete")

            node.querySelector("div.sr-only").textContent = "Delete"

            node.querySelector("button.form-value-delete").addEventListener("click", deleteFormValue)

            formValue.before(node)
        }

        document.querySelectorAll("button.form-value-delete").forEach( el =>
            el.addEventListener("click", deleteFormValue)
        )

        document.querySelectorAll("button.form-value-add").forEach( el =>
            el.addEventListener("click", addFormValue)
        )
    }

    // Required for HTMX. Ensures DOM has been fully updated.
    setTimeout(multiple, 40)
}
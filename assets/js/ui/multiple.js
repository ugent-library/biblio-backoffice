import htmx from 'htmx.org';

// Handles fields with multiple values
export default function() {
    const reTmpl = /^data-tmpl-(.+)/;

    const setValueIndex = (formValue, valueIndex) => {
        Array.from(formValue.getElementsByTagName('*')).forEach(function(el) {
            if (el.hasAttributes()) {
                let attrs = el.attributes;
                for (var i = 0; i < attrs.length; i++) {
                    let m = attrs[i].name.match(reTmpl);
                    if (m) {
                        el.setAttribute(m[1], attrs[i].value.replace('{i}', valueIndex))                            
                    }
                }
            }
        })                       
    }

    const addEvents = (rootEl) => {
        // Delete a value from the field
        let deleteFormValue = function (e) {
            let formField = e.target.closest("div.form-values")
            e.target.closest("div.form-value").remove()
            let length = Array.from(formField.children).length

            for (var valueIndex = 0; valueIndex < length; valueIndex++) {
                setValueIndex(formField.children[valueIndex], valueIndex)
            }
        }

        // Add a new value to the field
        let addFormValue = function (e) {
            let formField = e.target.closest("div.form-values")
            let lastValue = formField.lastElementChild
            let valueIndex = Array.from(formField.children).length

            let newValue = lastValue.cloneNode(true)
            newValue.querySelector(".form-control").value = ""
            newValue.querySelectorAll(".is-invalid").forEach(
                item => {
                    item.classList.remove("is-invalid")
                }
            )

            // set html attrs from their templates
            setValueIndex(newValue, valueIndex)

            // switch last value button to delete
            let lastBtn = lastValue.querySelector("button.form-value-add")
            let classList = lastBtn.classList
            classList.remove("form-value-add")
            classList.remove("btn-outline-primary")
            classList.add("btn-link-muted")
            classList.add("form-value-delete")
            classList = lastValue.querySelector("i.if-add").classList
            classList.remove("if-add")
            classList.add("if-delete")
            lastValue.querySelector("div.sr-only").textContent = "Delete"
            lastBtn.removeEventListener("click", addFormValue)
            lastBtn.addEventListener("click", deleteFormValue)

            // insert new value
            lastValue.after(newValue)
            // activate htmx on new element
            htmx.process(newValue)
            // activate add button
            newValue.querySelector("button.form-value-add").addEventListener("click", addFormValue)
            // fire added event
            newValue.dispatchEvent(new CustomEvent('form-value-add', {bubbles: true}))
        }

        rootEl.querySelectorAll("button.form-value-delete").forEach(el =>
            el.addEventListener("click", deleteFormValue)
        )

        rootEl.querySelectorAll("button.form-value-add").forEach(el =>
            el.addEventListener("click", addFormValue)
        )
    };

    htmx.onLoad(addEvents);
}
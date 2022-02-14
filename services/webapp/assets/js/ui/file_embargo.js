import htmx from 'htmx.org';

export default function () {
    let addEvents = function() {
        document.querySelectorAll('#publication-files-modal input[type=radio][name=access_level]').forEach(function(radio) {
            radio.addEventListener('change', function(evt) {
                const container = document.getElementById('file-embargo')
                const embargo   = container.querySelector('input[name=embargo]')
                const embargo_to= container.querySelector('select[name=embargo_to]')
                if (this.value === "open_access") {
                    container.classList.add("d-none")
                    embargo.disabled = true
                    embargo_to.disabled = true
                } else if (this.value === "local") {
                    container.classList.remove("d-none")
                    const opt = embargo_to.querySelector("option[value=local]")
                    opt.selected = false
                    opt.disabled = true
                    embargo.disabled = false
                    embargo_to.disabled = false
                } else if (this.value === "closed") {
                    container.classList.remove("d-none")
                    const opt = embargo_to.querySelector("option[value=local]")
                    opt.disabled = false
                    embargo.disabled = false
                    embargo_to.disabled = false
                }
            })
        })
    }

    // addEvents()

    // TODO don't use afterSettle
    htmx.on("htmx:afterSettle", function(evt) {
        addEvents()
    });
}

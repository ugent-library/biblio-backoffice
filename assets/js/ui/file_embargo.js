import htmx from 'htmx.org';

export default function () {
    let addEvents = function() {
        document.querySelectorAll('#publication-files-modal form input[type=radio][name=access_level]').forEach(function(radio) {
            radio.addEventListener('change', function(evt) {
                const container = document.getElementById('file-embargo')
                if (this.value === "open_access") {
                    container.classList.add("d-none")
                    container.querySelector("input[name=embargo]").value = ""
                    container.querySelector("select[name=embargo_to] option").forEach(opt => opt.selected = false)
                } else if (this.value === "local") {
                    container.classList.remove("d-none")
                    const opt = container.querySelector("select[name=embargo_to] option[value=local]")
                    opt.selected = false
                    opt.disabled = true
                } else if (this.value === "closed") {
                    container.classList.remove("d-none")
                    const opt = container.querySelector("select[name=embargo_to] option[value=local]")
                    opt.disabled = false
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
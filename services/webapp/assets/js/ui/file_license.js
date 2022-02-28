import htmx from "htmx.org";

export default function () {
    let addEvents = function(rootEl) {
        rootEl.querySelectorAll(".file-attributes").forEach(function(fileAttrs) {
            let ccLicense = fileAttrs.querySelector("[name=cc_license]");
            let noLicense = fileAttrs.querySelector("[name=no_license]");
            // let otherLicense = fileAttrs.querySelector("[name=other_license]");

            ccLicense.addEventListener('change', function() {
                if (this.value !== "") {
                    noLicense.checked = false
                }
            });
            noLicense.addEventListener('change', function() {
                if (this.checked) {
                    ccLicense.value = ""
                }
            });

        })
    }

    htmx.onLoad(addEvents)
}

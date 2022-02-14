import htmx from "htmx.org";

export default function () {
    let addEvents = function() {

        const container = document.getElementById("publication-files-modal");
        if(container === null) return;

        const access_levels = container.querySelectorAll("input[type=radio][name=access_level]");
        if(access_levels.length === 0) return;

        const cc_license = document.querySelector("select[name=cc_license]");
        if(cc_license === null) return;

        // default setting (no access level selected)
        cc_license.disabled = true; //prevents input from being submitted
        cc_license.parentNode.classList.add("d-none");

        access_levels.forEach(function(radio) {

            // handling at loading time (could be done in template also)
            if(radio.checked){
                cc_license.disabled = radio.value !== "open_access";
                if(cc_license.disabled){
                    cc_license.parentNode.classList.add("d-none");
                }
                else {
                    cc_license.parentNode.classList.remove("d-none");
                }
            }

            // handling when changed
            radio.addEventListener("change", function(){
                cc_license.disabled = this.value !== "open_access";
                if(cc_license.disabled){
                    cc_license.parentNode.classList.add("d-none");
                }
                else {
                    cc_license.parentNode.classList.remove("d-none");
                }
            });

        });

    };

    htmx.on("htmx:afterSettle", function(evt) {
        addEvents();
    });

};

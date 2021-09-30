import htmx from 'htmx.org';
// import BSN from "bootstrap.native/dist/bootstrap-native-v4";

export default function() {

    // Compose a spinner element
    function createSpinner() {
        const spinner = document.createElement("div")
        spinner.classList.add('spinner-border')

        const text = document.createElement("span")
        text.classList.add("sr-only")
        let cta = document.createTextNode("Loading...")
        text.appendChild(cta)

        spinner.appendChild(text)

        return spinner
    }

    // On submit, disable the cancel / save buttons & set the spinner
    function formSubmit(form) {
        const submitButton = form.querySelector('.btn-save')
        const cancelButton = form.querySelector('.btn-cancel')

        // Load the spinner when the button is clicked
        htmx.on(submitButton, "click", function(evt) {
            const spinner = createSpinner()
            submitButton.after(spinner)
        })

        // Disable the buttons after HTMX has started, but before the XHR request is
        // dispatched. Doing this on the 'click' event blocks triggering the HTMX lifecycle.
        //
        // See: https://github.com/bigskysoftware/htmx/issues/394
        htmx.on("htmx:beforeRequest", function(evt) {
            submitButton.setAttribute("disabled", "")
            cancelButton.setAttribute("disabled", "")
        })
    }

    // After submission, auto-dismiss all alerts after 10 seconds.
    // TODO: if 2 consecutive save actions happen within the 10 second interval,
    //    the first displayed alert will be destroyed by HTMX, causing the setTimeout
    //    to trigger a runtime error as it tries to apply a BSN.alert on a non-existing
    //    element.
    // function closeAlerts() {
    //     let alerts = document.querySelectorAll('.alert')
    //     alerts.forEach((el) => {
    //         setTimeout(() => {
    //             let alert = new BSN.Alert(el)
    //             alert.close()
    //         }, 10000)
    //     })
    // }

    // Init event listeners whenever HTMX swaps in a card-collapsible having a form element.
    htmx.on("htmx:afterSettle", function(evt) {
        let item = evt.detail.target.children.item(0)
        if (item && item.nodeName && (item.nodeName.toLowerCase() == "form")) {
            formSubmit(item)
        } else {
            // closeAlerts()
        }
    });
}
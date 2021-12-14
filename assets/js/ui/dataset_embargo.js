import htmx from 'htmx.org';

export default function () {
    let addEvents = function() {
        const form = document.querySelector('.dataset-details form');
        if (form !== null) {
            const accessLevel = form.querySelector('select[name=access_level]');

            const embargo = form.querySelector('input[name=embargo]');
            const embargoTo = form.querySelector('select[name=embargo_to');

            // Defaults until embargo gets selected
            if (accessLevel.value != 'info:eu-repo/semantics/embargoedAccess') {
                embargo.setAttribute('readonly', true);
                embargoTo.setAttribute('readonly', true);
            }

            let handler = function(evt) {
                switch (accessLevel.value) {
                    case 'info:eu-repo/semantics/embargoedAccess':
                        embargo.removeAttribute('readonly');
                        embargoTo.removeAttribute('readonly');
                        break;
                    default:
                        embargo.setAttribute('readonly', true);
                        embargoTo.setAttribute('readonly', true);
                        break;
                }
            }

            accessLevel.addEventListener('change',handler)
        }
    }

    htmx.on("htmx:afterSettle", function (evt) {
        addEvents();
    });
}
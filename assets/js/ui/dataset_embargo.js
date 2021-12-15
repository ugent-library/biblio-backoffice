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
                embargo.setAttribute('disabled', true);
                embargoTo.setAttribute('disabled', true);
            }

            let hidden = function(el, name) {
                const parent = el.parentNode;
                if (parent.querySelector('*[name=' + name + '][type=hidden]') == null) {
                    const hidden = document.createElement('input');
                    hidden.type = "hidden";
                    hidden.name = name;
                    hidden.value = "";
                    hidden.classList.add('hidden-embargo');
                    el.parentNode.appendChild(hidden);
                }
            }

            let handler = function(evt) {
                switch (accessLevel.value) {
                    case 'info:eu-repo/semantics/embargoedAccess':
                        embargo.removeAttribute('disabled');
                        embargoTo.removeAttribute('disabled');
                        document.querySelectorAll('.hidden-embargo').forEach(function(hidden) {
                            hidden.remove();
                        });
                        break;
                    default:
                        embargo.setAttribute('disabled', true);
                        embargo.value = "";
                        hidden(embargo, "embargo");
                        embargoTo.setAttribute('disabled', true);
                        embargoTo.value = "";
                        hidden(embargoTo, "embargo_to");
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
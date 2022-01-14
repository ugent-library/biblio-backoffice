import htmx from 'htmx.org';

export default function () {
    let addEvents = function(rootEl) {
        rootEl.querySelectorAll('.autocomplete-hit').forEach(function(hit) {
            hit.addEventListener('click', function () {
                let target = hit.closest(".autocomplete").dataset.target;
                let v = hit.dataset.value;
                document.querySelector(target).value = v;
                hit.closest(".autocomplete-hits").innerHTML = "";                               
            })
        })
    }

    addEvents(document);
    htmx.onLoad(function(el) {
        addEvents(el)
    });
}

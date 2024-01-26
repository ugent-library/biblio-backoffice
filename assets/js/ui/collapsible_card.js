// /* ==========================================================================
//     Collapsible card content
//    ========================================================================== */

import htmx from 'htmx.org';

export default function() {
    let toggler = function(evt) {
        let showCardInfo = evt.currentTarget.closest('.hstack-md-responsive')
        let cardContent = showCardInfo.querySelector('[card-content]')
        cardContent.classList.toggle('show')
    }

    let addEvents = function() {
        document.querySelectorAll('[card-content-toggle]').forEach(cardContent =>
            cardContent.addEventListener('click', toggler)
        )
    }

    addEvents()

    htmx.onLoad(function(el) {
        addEvents()
    });
}
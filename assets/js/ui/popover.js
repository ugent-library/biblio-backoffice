import htmx from 'htmx.org';
import {Popover} from "bootstrap.native";

export default function () {
    let addEvents = function(rootEl) {
        rootEl.querySelectorAll('[data-toggle=popover-custom]').forEach(function(el) {
            let container = document.querySelector(el.dataset.popoverContent)
            let content = container.querySelector('.popover-body')
            let heading = container.querySelector('.popover-heading')
            let title = ""
            if (heading) {
                title = heading.innerHTML
            }
            new Popover(el, {
                content: content.innerHTML,
                title: title,
                delay: 1000,
            })
        })
    }

    htmx.onLoad(addEvents);
}

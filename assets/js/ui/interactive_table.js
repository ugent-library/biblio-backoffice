import htmx from 'htmx.org';

export default function() {
    htmx.on("htmx:beforeRequest", function(evt) {
        let item = evt.detail

        if (item.elt.classList.contains("btn-swap")) {
            for (let i = 0; i < item.target.parentElement.children.length; i++) {
                let row = item.target.parentElement.children[i]
                let buttons = row.getElementsByTagName("button")
                Array.from(buttons).forEach(function (button) {
                    if (! button.classList.contains("create")) {
                        button.classList.add("d-none")
                    }
                });
            }
        }
    });
}
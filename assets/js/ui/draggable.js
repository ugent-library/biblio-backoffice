import htmx from 'htmx.org';
import { DraggableTable } from '../classes/draggableTable.js';

// TODO: We likely want to turn this into a separate class too.
//  we just want to set the callback and ensure everything gets hooked
//  to specific DOM elements, rather then generic selectors as we do now.
export function draggable() {

    // Disable buttons when we edit / add a row.
    let disableRowButtons = function(evt) {
        let rows = document.querySelector('table.inline-editing tbody').children;

        for (let i = 0; i < rows.length; i++) {
            let row = rows[i];

            if ( ! row.classList.contains("row-new") && ! row.classList.contains("row-edit")) {
                let buttons = row.getElementsByTagName("button");
                Array.from(buttons).forEach(function (button) {
                    button.classList.add("d-none")
                });
            }
        }
    }

    let currentUrl = new URL(window.location.href);
    let callback = "";
    let tableSelector;

    // Generate the callback for the authors
    // TODO: Turn the pattern matching into something that incorporates the basePath of the Go app.
    if (currentUrl.pathname.match(new RegExp(`^\.*\/publication\/[0-9]*$`, 'gm'))) {
        if (currentUrl.hash == "#contributors-content") {
            callback = currentUrl.pathname + "/htmx/authors/order/:start/:end"
            tableSelector = "foobar";
        }
    }

    if (tableSelector !== undefined) {
        // Init the Draggable table
        const table = new DraggableTable({ tableSelector: tableSelector, callback: callback });
        table.init();

        // Init the addAuthor button
        let addAuthor = document.querySelector("button.btn-outline-primary.add-author");

        // ... This is where the magic starts to happen ...

        // After the table is refreshed
        htmx.on("ITListAfterSwap", function(evt) {
            // Make the table draggable again.
            table.init();
            // We can click the 'add author' button from the top menu.
            addAuthor.removeAttribute("disabled");
        });

        // After the order was changed thru drag n' droppin'
        htmx.on("ITOrderAuthorsAfterSwap", function(evt) {
            table.init();
        });

        // A new empty row form was added. Disable all add / edit / delete buttons
        htmx.on("ITAddRowAfterSwap", function(evt) {
            // Make the table static non-draggable.
            table.reset();
            // Disable the 'Add author' button from the top menu.
            addAuthor.setAttribute("disabled", "true");
            // Remove all buttons except the one's on the active form.
            disableRowButtons(evt);
        });

        // An row is being edited. Disable all add / edit / delete buttons
        htmx.on("ITEditRowAfterSwap", function(evt) {
            table.reset();
            addAuthor.setAttribute("disabled", "true");
            disableRowButtons(evt);
        });
    }
}
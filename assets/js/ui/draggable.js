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

    // Compose a spinner element
    let createSpinner = function() {
        const spinner = document.createElement("div")
        spinner.classList.add('spinner-border')

        const text = document.createElement("span")
        text.classList.add("sr-only")
        let cta = document.createTextNode("Loading...")
        text.appendChild(cta)

        spinner.appendChild(text)

        return spinner
    }

    let currentUrl = new URL(window.location.href);
    let callback = "";
    let tableSelector;

    // Generate the callback for the authors
    // TODO: Turn the pattern matching into something that incorporates the basePath of the Go app.
    if (currentUrl.pathname.match(new RegExp(`^\.*\/publication\/[0-9]*$`, 'gm'))) {
        if (currentUrl.hash == "#contributors-content") {
            callback = currentUrl.pathname + "/htmx/authors/order/:start/:end"
            tableSelector = "authors-table";
        }
    }

    if (tableSelector !== undefined) {
        // Init the Draggable table
        const table = new DraggableTable({ tableSelector: tableSelector, callback: callback });
        table.init();

        // Init the addAuthor button
        let addAuthor = document.querySelector("button.btn-outline-primary.add-author");

        // Init a spinner (loading...) object
        let spinner;

        // ... This is where the magic starts to happen ...

        // After the table is refreshed
        htmx.on("ITListAfterSwap", function(evt) {
            // Make the table draggable again.
            table.init();
            // We can click the 'add author' button from the top menu.
            addAuthor.removeAttribute("disabled");

            // We remove the spinner if there is any active
            if (spinner !== undefined) {
                spinner.remove();
            }
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
            // Add spinner handler to the 'create' button.
            let updateButton = document.querySelector("table#authors-table button.create-author");
            updateButton.addEventListener("click", function(e) {
                spinner = createSpinner();
                addAuthor.after(spinner);
            })
        });

        // An row is being edited. Disable all add / edit / delete buttons
        htmx.on("ITEditRowAfterSwap", function(evt) {
            table.reset();
            addAuthor.setAttribute("disabled", "true");
            disableRowButtons(evt);
            let createButton = document.querySelector("table#authors-table button.update-author");
            createButton.addEventListener("click", function(e) {
                spinner = createSpinner();
                addAuthor.after(spinner);
            })
        });

        // A row is being deleted. Show a spinner next to the 'delete' button in the popup
        htmx.on("ITConfirmRemoveFromPublicationAfterSwap", function(evt) {
            let removeButton = document.querySelector("div.modal-confirm-author-removal button.delete-author");
            removeButton.addEventListener("click", function(e) {
                spinner = createSpinner();
                removeButton.after(spinner);
            })
        })
    }
}
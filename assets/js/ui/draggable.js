import htmx from 'htmx.org';
import { DraggableTable } from '../classes/draggableTable.js';

export function draggable() {

    let currentUrl = new URL(window.location.href);
    let callback = "";
    let tableSelector;

    // Generate the callback for the authors
    if (currentUrl.pathname.match(new RegExp(`^\\/publication\/[0-9]*$`, 'gm'))) {
        if (currentUrl.hash == "#contributors-content") {
            callback = currentUrl.pathname + "/htmx/authors/order/:start/:end"
            tableSelector = "foobar";
        }
    }

    if (tableSelector !== undefined) {
        const table = new DraggableTable({ tableSelector: tableSelector, callback: callback });

        table.init();

        htmx.on("ITListAfterSwap", function(evt) {
            table.init();
        });

        htmx.on("ITOrderAuthorsAfterSwap", function(evt) {
            table.init();
        });

        htmx.on("ITAddRow", function(evt) {
            table.reset();
        });

        htmx.on("ITEditRow", function(evt) {
            table.reset();
        });
    }
}
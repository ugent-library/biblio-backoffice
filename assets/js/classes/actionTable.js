import htmx from 'htmx.org';

export class ActionTable {

    constructor({ tableSelector, addButtonSelector }) {
        this._tableSelector = tableSelector;
        this._addButtonSelector = addButtonSelector;
        this._table = undefined;
        this._addButton = undefined;
        this._spinner = undefined;
    }

    _disableRowButtons() {
        let rows = this._table.querySelector('tbody').children;

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

    _createSpinner() {
        const spinner = document.createElement("div")

        spinner.classList.add('spinner-border')
        const text = document.createElement("span")
        text.classList.add("sr-only")
        let cta = document.createTextNode("Loading...")
        text.appendChild(cta)
        spinner.appendChild(text)

        return spinner
    }

    _list() {
        // We can click the 'add author' button from the top menu.
        if (this._addButton !== undefined) {
            this._addButton.removeAttribute("disabled");
        }

        // We remove the spinner if there is any active
        if (this._spinner !== undefined) {
            this._spinner.remove();
        }
    }

    _addRow() {
        // Remove all buttons except the one's on the active form.
        this._disableRowButtons();

        let createButton = this._table.querySelector("button.create-contributor");

        if (createButton !== undefined && this._addButton !== undefined) {
            // Disable the 'Add author' button from the top menu.
            this._addButton.setAttribute("disabled", "true");

            // Init spinner
            createButton.addEventListener("click", (function(e) {
                this._spinner = this._createSpinner();
                this._addButton.after(this._spinner);
            }).bind(this));
        }
    }

    _editRow() {
        // Remove all buttons except the one's on the active form.
        this._disableRowButtons();

        let updateButton = this._table.querySelector("button.update-contributor");

        if (updateButton !== undefined && this._addButton !== undefined) {
            // Disable the 'Add author' button from the top menu.
            this._addButton.setAttribute("disabled", "true");

            // Init spinner
            updateButton.addEventListener("click", function(e) {
                this._spinner = this._createSpinner();
                this._addButton.after(this._spinner);
            }.bind(this));
        }
    }

    _removeRow() {
        let removeButton = document.querySelector(".modal-dialog button.delete-contributor");

        if (removeButton !== undefined) {
            removeButton.addEventListener("click", (function(e) {
                this._spinner = this._createSpinner();
                removeButton.after(this._spinner);
            }).bind(this))
        }
    }

    init() {
        this._table = document.querySelector(this._tableSelector);
        this._addButton = document.querySelector(this._addButtonSelector);

        if (this._table !== undefined) {
            htmx.on("ITListAfterSwap", (function(evt) {
                this._list();
            }).bind(this));

            htmx.on("ITAddRowAfterSwap", (function(evt) {
                this._addRow();
            }).bind(this));

            htmx.on("ITEditRowAfterSwap", (function(evt) {
                this._editRow();
            }).bind(this));

            htmx.on("ITConfirmRemoveFromPublicationAfterSwap", (function(evt) {
                this._removeRow();
            }).bind(this));
        }
    }
}
import htmx from 'htmx.org';

/**
 * DraggableTable
 *
 * Turn a table into a table with rows which can be drag and dropped.
 * This allows users to (re)order the rows in the table. An optional API callback
 * can be passed to the constructor. This allows sending the the state of the table
 * to the backend.
 *
 * Code repurposed from: https://htmldom.dev/drag-and-drop-table-row/
 */
export class DraggableTable {
    _table;
    _tableSelector;
    _callback;
    _draggingEle;
    _draggingRowIndex;
    _placeholder;
    _list;
    _isDraggingStarted;
    _x;
    _y;
    _mouseUp;
    _mouseMove;
    _mouseDown;

    /**
     * Dragable Table constructor.
     * @param { string } DOM element to be selected. It must be a HTML Table tag - https://developer.mozilla.org/en-US/docs/Web/HTML/Element/table
     * @param { callback } The backend callback to notify after swapping the element.
     */
    constructor({ tableSelector, callback }) {
        this._tableSelector = tableSelector;
        this._callback = callback ?? "";
        this._isDraggingStarted = false;

        // The current position of mouse relative to the dragging element
        this._x = 0;
        this._y = 0;
    }

    _swap(nodeA, nodeB) {
        const parentA = nodeA.parentNode;
        const siblingA = nodeA.nextSibling === nodeB ? nodeA : nodeA.nextSibling;

        // Move `nodeA` to before the `nodeB`
        nodeB.parentNode.insertBefore(nodeA, nodeB);

        // Move `nodeB` to before the sibling of `nodeA`
        parentA.insertBefore(nodeB, siblingA);
    }

    _isAbove(nodeA, nodeB) {
        // Get the bounding rectangle of nodes
        const rectA = nodeA.getBoundingClientRect();
        const rectB = nodeB.getBoundingClientRect();

        return rectA.top + rectA.height / 2 < rectB.top + rectB.height / 2;
    }

    _cloneTable() {
        // const rect = this._table.getBoundingClientRect();
        const width = parseInt(window.getComputedStyle(this._table).width);

        this._list = document.createElement('div');
        this._list.classList.add('clone-list');
        // list.style.position = 'absolute';
        // list.style.left = `${rect.left}px`;
        // list.style.top = `${rect.top}px`;
        this._table.parentNode.insertBefore(this._list, this._table);

        // Hide the original table
        this._table.style.visibility = 'hidden';
        this._table.classList.add("d-none");

        this._table.querySelectorAll('tr').forEach(function (row) {
            // Create a new table from given row
            const item = document.createElement('div');
            item.classList.add('draggable');

            const newTable = document.createElement('table');
            newTable.setAttribute('class', 'clone-table table');
            newTable.style.width = `${width}px`;

            const newRow = document.createElement('tr');
            const cells = [].slice.call(row.children);
            cells.forEach(function (cell) {
                const newCell = cell.cloneNode(true);
                newRow.appendChild(newCell);
            });

            newTable.appendChild(newRow);
            item.appendChild(newTable);
            this._list.appendChild(item);
        }, this);
    }

    _mouseDownHandler(e) {
        // Get the original row
        const originalRow = e.target.closest('tr');
        this._draggingRowIndex = [].slice.call(this._table.querySelectorAll('tr')).indexOf(originalRow);

        // Determine the mouse position
        this._x = e.clientX;
        this._y = e.clientY;

        // Attach the listeners to `document`
        this._mouseMove = this._mouseMoveHandler.bind(this);
        document.addEventListener('mousemove', this._mouseMove);

        this._mouseUp = this._mouseUpHandler.bind(this);
        document.addEventListener('mouseup', this._mouseUp);
    }

    _mouseMoveHandler(e) {
        if (!this._isDraggingStarted) {
            this._isDraggingStarted = true;

            this._cloneTable();
            this._draggingEle = [].slice.call(this._list.children)[this._draggingRowIndex];
            this._draggingEle.classList.add('dragging');


            // Let the placeholder take the height of dragging element
            // So the next element won't move up
            this._placeholder = document.createElement('div');
            this._placeholder.classList.add('placeholder');
            this._draggingEle.parentNode.insertBefore(this._placeholder, this._draggingEle.nextSibling);
            this._placeholder.style.height = `${this._draggingEle.offsetHeight}px`;
        }

        // Set position for dragging element
        this._draggingEle.style.position = 'absolute';
        this._draggingEle.style.top = `${this._draggingEle.offsetTop + e.clientY - this._y}px`;
        this._draggingEle.style.left = `${this._draggingEle.offsetLeft + e.clientX - this._x}px`;

        // Reassign the position of mouse
        this._x = e.clientX;
        this._y = e.clientY;

        // The current order
        // prevEle
        // draggingEle
        // placeholder
        // nextEle
        const prevEle = this._draggingEle.previousElementSibling;
        const nextEle = this._placeholder.nextElementSibling;

        // The dragging element is above the previous element
        // User moves the dragging element to the top
        // We don't allow to drop above the header
        // (which doesn't have `previousElementSibling`)
        if (prevEle && prevEle.previousElementSibling && this._isAbove(this._draggingEle, prevEle)) {
            // The current order    -> The new order
            // prevEle              -> placeholder
            // draggingEle          -> draggingEle
            // placeholder          -> prevEle
            this._swap(this._placeholder, this._draggingEle);
            this._swap(this._placeholder, prevEle);
            return;
        }

        // The dragging element is below the next element
        // User moves the dragging element to the bottom
        if (nextEle && this._isAbove(nextEle, this._draggingEle)) {
            // The current order    -> The new order
            // draggingEle          -> nextEle
            // placeholder          -> placeholder
            // nextEle              -> draggingEle
            this._swap(nextEle, this._placeholder);
            this._swap(nextEle, this._draggingEle);
        }
    }

    _mouseUpHandler(e) {
        // Remove the placeholder
        this._placeholder && this._placeholder.parentNode.removeChild(this._placeholder);

        this._draggingEle.classList.remove('dragging');
        this._draggingEle.style.removeProperty('top');
        this._draggingEle.style.removeProperty('left');
        this._draggingEle.style.removeProperty('position');

        // Get the end index
        const endRowIndex = [].slice.call(this._list.children).indexOf(this._draggingEle);

        this._isDraggingStarted = false;

        // Remove the `list` element
        this._list.parentNode.removeChild(this._list);

        // Move the dragged row to `endRowIndex`
        let rows = [].slice.call(this._table.querySelectorAll('tr'));
        this._draggingRowIndex > endRowIndex
            ? rows[endRowIndex].parentNode.insertBefore(rows[this._draggingRowIndex], rows[endRowIndex])
            : rows[endRowIndex].parentNode.insertBefore(
                    rows[this._draggingRowIndex],
                    rows[endRowIndex].nextSibling
                );

        // Bring back the table
        this._table.style.removeProperty('visibility');
        this._table.classList.remove('d-none');

        if (this._callback !== "") {
            let start = this._draggingRowIndex -1
            let end = endRowIndex -1

            let callback = this._callback;

            callback = callback.replace(':start', start);
            callback =  callback.replace(':end', end);
            htmx.ajax('PUT', callback, this._table.querySelector("tbody"))
        }

        // Remove the handlers of `mousemove` and `mouseup`
        document.removeEventListener('mousemove', this._mouseMove);
        document.removeEventListener('mouseup', this._mouseUp);
    }

    init() {
        this._table = document.getElementById(this._tableSelector);

        this._table.querySelectorAll('tr').forEach(function (row, index) {
            // Ignore the header
            // We don't want user to change the order of header
            if (index === 0) {
                return;
            }

            const firstCell = row.firstElementChild;
            firstCell.classList.add('draggable');

            this._mouseDown = this._mouseDownHandler.bind(this);
            firstCell.addEventListener('mousedown', this._mouseDown);
        }, this);
    }

    reset() {
        this._table = document.getElementById(this._tableSelector);

        this._table.querySelectorAll('tr').forEach(function (row, index) {
            // Ignore the header
            // We don't want user to change the order of header
            if (index === 0) {
                return;
            }

            const firstCell = row.firstElementChild;
            firstCell.classList.remove('draggable');

            firstCell.removeEventListener('mousedown', this.__mouseDown);
        }, this);
    }
}
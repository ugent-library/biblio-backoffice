import { ActionTable } from '../classes/actionTable.js';

export function actionTables() {
    document.querySelectorAll("table.contributor-table").forEach(function(table) {
        const actionTable = new ActionTable({
            tableSelector: "table#" + table.id,
            addButtonSelector: "button.btn-outline-primary.add-contributor"
        })
        actionTable.init()
    })
}

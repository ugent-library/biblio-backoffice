import { ActionTable } from '../classes/actionTable.js';

export function actionTables() {

    // Authors table
    const authorsTable = new ActionTable({
        tableSelector: "table#authors-table",
        addButtonSelector: "button.btn-outline-primary.add-author"
    })

    authorsTable.init();
}

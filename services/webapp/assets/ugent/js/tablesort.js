const Tablesort = require('tablesort');

if (document.querySelector('.table-sortable')) {
    new Tablesort(document.querySelector('.table-sortable'));
}

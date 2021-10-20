import check from './ui/check.js'
import bootstrap from './ui/bootstrap.js'
import { draggable } from './ui/draggable.js'
import multiple from './ui/multiple.js'
import changeSubmit from './ui/form_change_submit.js'
import submit from './ui/form_submit.js'
import modalClose from './ui/modal_close.js'
import modalPopper from './ui/modal_popper.js'
import multipleSelect from './ui/multi_select.js'
import tabs from './ui/tabs.js'

document.addEventListener('DOMContentLoaded', function () {
    tabs();
    check();
    bootstrap();
    draggable();
    multiple();
    changeSubmit();
    submit();
    modalClose();
    modalPopper();
    multipleSelect();
});
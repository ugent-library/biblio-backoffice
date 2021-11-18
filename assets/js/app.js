import csrf from './ui/csrf.js'
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
import radioCard from './ui/radio_card.js'
import promote from './ui/promote.js'
import toast from './ui/toast.js'

document.addEventListener('DOMContentLoaded', function () {
    csrf()
    tabs()
    check()
    bootstrap()
    draggable()
    multiple()
    changeSubmit()
    submit()
    modalClose()
    modalPopper()
    multipleSelect()
    radioCard()
    promote()
    toast()
});
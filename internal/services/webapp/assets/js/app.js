import htmx from 'htmx.org';
import csrf from './ui/csrf.js'
import checkbox from './ui/checkbox.js'
import bootstrap from './ui/bootstrap.js'
import bootstrapPopper from './ui/bootstrap_popper.js'
import popover from './ui/popover.js'
import header from './ui/header.js'
import multiple from './ui/multiple.js'
import changeSubmit from './ui/form_change_submit.js'
import autocomplete from './ui/autocomplete.js';
import submit from './ui/form_submit.js'
import modalClose from './ui/modal_close.js'
// import tabs from './ui/tabs.js'
import radioCard from './ui/radio_card.js'
import toast from './ui/toast.js'
import sortable from './ui/sortable.js';

// configure htmx
htmx.config.defaultFocusScroll = true
// load htmx extensions
window.htmx = htmx
require('htmx.org/dist/ext/remove-me.js');

// initialize everyting
document.addEventListener('DOMContentLoaded', function () {
    csrf()
    // tabs()
    checkbox()
    bootstrap()
    bootstrapPopper()
    popover()
    header()
    multiple()
    changeSubmit()
    autocomplete()
    submit()
    modalClose()
    radioCard()
    toast()
    sortable()
});

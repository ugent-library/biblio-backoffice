import htmx from 'htmx.org';
import csrf from './ui/csrf.js'
import checkbox from './ui/checkbox.js'
import bootstrap from './ui/bootstrap.js'
import bootstrapPopper from './ui/bootstrap_popper.js'
import popover from './ui/popover.js'
import header from './ui/header.js' // TODO is this still needed?
import multiple from './ui/multiple.js'
import changeSubmit from './ui/form_change_submit.js'
import autocomplete from './ui/autocomplete.js';
import modalClose from './ui/modal_close.js'
import radioCard from './ui/radio_card.js'
import toast from './ui/toast.js'
import sortable from './ui/sortable.js';
import collapseSubSidebar from './ui/collapsible_sub_sidebar.js'
import formAutoFocus from './ui/form_auto_focus.js'
import fileUpload from './ui/file_upload.js'

// configure htmx
htmx.config.defaultFocusScroll = true
// load htmx extensions
window.htmx = htmx
require('htmx.org/dist/ext/remove-me.js');

// initialize everyting
document.addEventListener('DOMContentLoaded', function () {
    csrf()
    checkbox()
    bootstrap()
    bootstrapPopper()
    popover()
    header()
    multiple()
    changeSubmit()
    autocomplete()
    modalClose()
    radioCard()
    toast()
    sortable()
    collapseSubSidebar()
    formAutoFocus()
    fileUpload()
});

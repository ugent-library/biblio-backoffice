import htmx from "htmx.org/dist/htmx.esm.js";

import { initCallback } from "bootstrap.native";
import csrf from "./ui/csrf.js";
import checkbox from "./ui/checkbox.js";
import clipboard from "./ui/clipboard.js";
import popover from "./ui/popover.js";
import header from "./ui/header.js"; // TODO is this still needed?
import multiple from "./ui/multiple.js";
import changeSubmit from "./ui/form_change_submit.js";
import autocomplete from "./ui/autocomplete.js";
import modalClose from "./ui/modal_close.js";
import radioCard from "./ui/radio_card.js";
import toast from "./ui/toast.js";
import sortable from "./ui/sortable.js";
import collapseSubSidebar from "./ui/collapsible_sub_sidebar.js";
import fileUpload from "./ui/file_upload.js";
import tags from "./ui/tags.js";
import facetDropdowns from "./ui/facet_dropdowns.js";

// configure htmx
htmx.config.defaultFocusScroll = true;

htmx.onLoad((el) => {
  // apply Bootstrap JS to new DOM content (not on initial load)
  if (el !== document.body) {
    initCallback(el);
  }
});

// load htmx extensions (uncomment if any HTMX extensions are needed)
// window.htmx = htmx;

// initialize everything
document.addEventListener("DOMContentLoaded", function () {
  csrf();
  checkbox();
  popover();
  header();
  multiple();
  changeSubmit();
  autocomplete();
  modalClose();
  radioCard();
  toast();
  sortable();
  collapseSubSidebar();
  fileUpload();
  tags();
  facetDropdowns();
});

htmx.onLoad(function (el) {
  clipboard(el);
});

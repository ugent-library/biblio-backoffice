// Collapse.js
//   Init Bootstrap Native on the collapsible fieldset and children
//   whenever htmx:afterSwap event is triggered. This ensures all
//   The elements in the HTMX fragments get their BSN event listeners re-applied.
//
export default function(collapsible) {
    // re-init BSN on the element that just got swapped via HTMX
    BSN.initCallback(collapsible)
}
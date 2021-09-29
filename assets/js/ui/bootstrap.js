import htmx from 'htmx.org';
import BSN from "bootstrap.native/dist/bootstrap-native-v4";

/**
 * Initialize Bootstrap Native after HTMX has settled the DOM.
 *
 * When HTMX is executed, the updated parts of the DOM won't be
 * registered with Bootstrap Native. Elements like i.e. popovers,
 * alerts, tooltips,... won't work passed via HTMX won't work.
 * This function re-initializes Bootstrap Native on the updated DOM.
 */
export default function() {
    htmx.on("htmx:afterSettle", function(evt) {
        BSN.initCallback()
    });
}
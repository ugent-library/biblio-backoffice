import htmx from 'htmx.org';
import BSN from "bootstrap.native/dist/bootstrap-native-v4";

 // Reinitialize Bootstrap Native after HTMX has udated the DOM.
 export default function() {
    htmx.onLoad(function(el) {
        BSN.initCallback()
    });
}
import htmx from 'htmx.org';
import _BSN from "bootstrap.native";

//import '../ugent/js/index';
import multipleValues from './multiple.js'
import popovers from './popover.js'
import collapsible from './collapse.js'

(function main () {
    htmx.on("htmx:afterSwap", function(evt) {
        // TODO
        //   Only execute on add / edit forms of publication / dataset

        multipleValues()
        popovers()
        collapsible(evt.detail.target)
    });
})();
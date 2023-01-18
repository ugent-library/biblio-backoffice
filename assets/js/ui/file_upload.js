import htmx from "htmx.org";

export default function() {

  htmx.onLoad(function(doc) {
    let formUpload = doc.querySelector("#form-upload")

    // function is called several times, but not always with form-upload present
    if(formUpload == null) return;

    formUpload.addEventListener("htmx:xhr:loadstart", function(evt) {
        doc.querySelector("input[name=file]").disabled = true;
        doc.querySelector("#file-upload-start").setAttribute("style", "display: none");
        doc.querySelector("#file-upload-busy").setAttribute("style", "display: block");
    });

    formUpload.addEventListener("htmx:xhr:progress", function(evt) {
        // at the end of the request ajax returns an event with evt.detail.total == 0
        if(!evt.detail.lengthComputable) return;

        let pct = Math.floor(evt.detail.loaded/evt.detail.total * 100);
        let pb = doc.querySelector("#progress-bar");
        pb.setAttribute("style", "width: "+pct+"%");
        pb.setAttribute("aria-valuenow", pct);
        let pi = doc.querySelector("#progress-bar-info");
        pi.innerHTML = ""+pct+"%";
    });

  });

}

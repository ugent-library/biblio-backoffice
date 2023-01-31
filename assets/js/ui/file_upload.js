import htmx from "htmx.org";

// TODO handle error, abort
export default function() {
  htmx.onLoad(function(rootEl) {
    rootEl.querySelectorAll("input[type=file].upload-progress").forEach(input => {
        let form = input.closest("form")

        form.addEventListener("htmx:xhr:loadstart", function(evt) {
            input.disabled = true;
            form.querySelector(".file-upload-start").classList.add("d-none");
            form.querySelector(".file-upload-busy").classList.remove("d-none");
        });

        form.addEventListener("htmx:xhr:progress", function(evt) {
            // at the end of the request ajax returns an event with evt.detail.total == 0
            if(!evt.detail.lengthComputable) return;
    
            let pct = Math.floor(evt.detail.loaded/evt.detail.total * 100);
            let pb = form.querySelector(".progress-bar");
            pb.setAttribute("style", "width: "+pct+"%");
            pb.setAttribute("aria-valuenow", pct);
            let pi = form.querySelector(".progress-bar-percent");
            pi.innerText = pct;
        });
    })
  });
}
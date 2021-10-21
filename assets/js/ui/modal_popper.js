import Popper from 'popper.js'

// Wire in popper.js support. This ensure popups stay within the viewport.
//
// Bootstrap Native doesn't incorporate Popper.js. We have to wire everything ourselves.
// See: https://github.com/thednp/bootstrap.native/issues/211
export default function() {
    document.querySelectorAll("div.dropdown > button").forEach(function(button) {
        button.addEventListener("click", function(evt) {
            let menu = button.parentElement.children.item(1)

            if (menu.classList.contains("show")) {
                menu.removeAttribute("x-placement")
                menu.removeAttribute("style")

                let popper = new Popper(button, menu, {
                    modifiers: {
                        preventOverflow: { enabled: true },
                        flip: { enabled: true},
                        hide: { enabled: false}
                    }
                })
            }
        })
    })
}
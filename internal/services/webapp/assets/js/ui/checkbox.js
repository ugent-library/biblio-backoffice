export default function() {
    let formCheckAll = function (evt) {
        let form = evt.target.closest("form")
        let chks = form.querySelectorAll("input[type='checkbox']")
        chks.forEach(el =>
            el.checked = true
        )
    }

    document.querySelectorAll("button.form-check-all").forEach(el =>
        el.addEventListener("click", formCheckAll)
    )
}
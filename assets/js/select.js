export default function() {
    document.querySelectorAll("select.form-change-submit").forEach(el =>
        el.addEventListener("change", evt =>
            evt.target.closest("form").submit()
        )
    )
}
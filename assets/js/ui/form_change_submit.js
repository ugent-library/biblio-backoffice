export default function() {
    document.querySelectorAll("form.form-change-submit").forEach(function(el) {
        el.addEventListener("change", function(evt) {
            const spinner = el.querySelector(".spinner-border")
            if (spinner !== null) {
                spinner.style.display = "inline-block";
                spinner.style.opacity = "1";
            }

            el.submit()
        });
    });

    document.querySelectorAll("form .form-change-submit").forEach(function(el) {
        el.addEventListener("change", function(evt) {
            evt.target.closest("form").submit()
        });
    });
}
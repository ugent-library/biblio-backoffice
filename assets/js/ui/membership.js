import htmx from 'htmx.org';

export default function () {
    let modalCloseSecondary = function(modal) {
        modal.classList.remove("show")
        // Timeout gives us a fluid animation
        setTimeout(function() {
            if (modal) {
                modal.remove();
            }
        }, 100)
    }

    let membershipMember = function(target, id, firstName, lastName) {
        const targetEl = document.querySelector(target);
        targetEl.querySelector('form input[name="id"]').value = id;
        targetEl.querySelector('form input[name="first_name"]').value = firstName;
        targetEl.querySelector('form input[name="last_name"]').value = lastName;
        targetEl.querySelector('form input[name="first_name"]').setAttribute("readonly", "readonly")
        targetEl.querySelector('form input[name="last_name"]').setAttribute("readonly", "readonly")
    }

    let membershipExternal = function(target) {
        const targetEl = document.querySelector(target);
        targetEl.querySelector('form input[name="id"]').value = "";
        targetEl.querySelector('form input[name="first_name"]').removeAttribute("readonly")
        targetEl.querySelector('form input[name="last_name"]').removeAttribute("readonly")
    }

    let addEvents = function() {
        document.querySelectorAll('button.membership-member').forEach(btn =>
            btn.addEventListener('click', function(evt) {
                membershipMember(btn.dataset.target, btn.dataset.id, btn.dataset.firstName, btn.dataset.lastName)
                modalCloseSecondary(btn.closest(".modal"))
            })
        )
        document.querySelectorAll('input[type="radio"].membership-external').forEach(function(radio) {
            radio.addEventListener('change', function(evt) {
                if (!radio.checked) {
                    return
                }
                console.log('demote')
                membershipExternal(radio.dataset.target)
            })
        })

        document.querySelectorAll('button.membership-external').forEach(btn =>
            btn.addEventListener('click', function(evt){
                membershipExternal(btn.dataset.target)
                modalCloseSecondary(btn.closest(".modal"))
            })
        )
    }

    addEvents()

    // TODO don't use afterSettle
    htmx.on("htmx:afterSettle", function(evt) {
        addEvents()
    });
}
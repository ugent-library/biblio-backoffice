import htmx from 'htmx.org';

export default function () {

    // Promote an author to "UGent author"
    let promoteAuthor = function(evt) {
        let contributors = document.querySelectorAll("#modal-contributors-list li a")

        contributors.forEach(function (el) {
            el.addEventListener("click", function(evt) {
                let rowDelta = el.dataset.row;
                let item = el.closest("li.list-group-item");

                let row = document.querySelectorAll("table#author-table tbody tr").item(rowDelta)

                // Set the values on the input fields
                row.querySelector("input[name=first_name]").value = item.dataset.firstname;
                row.querySelector("input[name=last_name]").value = item.dataset.lastname;
                row.querySelector("input[name=ID]").value = item.dataset.id;

                // Disable the input fields (name shouldn't be editable if UGent author)
                row.querySelector("input[name=first_name]").setAttribute("disabled", "disabled");
                row.querySelector("input[name=last_name]").setAttribute("disabled", "disabled");

                // Flip the promote / degrade buttons
                row.querySelector('div.external').classList.remove('d-display');
                row.querySelector('div.external').classList.add('d-none');
                row.querySelector('div.ugent-author').classList.remove('d-none');
                row.querySelector('div.ugent-author').classList.add('d-display');

                // Close the modal
                let modal = document.querySelectorAll(".modal").item(0)
                let backdrop = document.querySelectorAll(".modal-backdrop").item(0)

                if (modal) {
                    modal.classList.remove("show")
                }

                if (backdrop) {
                    backdrop.classList.remove("show")
                }

                // Timeout gives us a fluid animation
                setTimeout(function() {
                    if (backdrop) {
                        backdrop.remove();
                    }

                    if (modal) {
                        modal.remove();
                    }
                }, 100)

                evt.preventDefault();
            }, false);
        });
    }

    // Demote an UGent author to "external member"
    let demoteAuthor = function (evt) {
        let demoteButtons = document.querySelectorAll("table#author-table tbody tr button.demote-external-member");

        demoteButtons.forEach(function (el) {
            el.addEventListener("click", function(evt) {
                let row = el.closest("tr")

                // Remove the UGent ID: this equals demoting the user
                row.querySelector("input[name=ID]").value = "";

                // Enable the input fields (editable as external member)
                row.querySelector("input[name=first_name]").removeAttribute("disabled");
                row.querySelector("input[name=last_name]").removeAttribute("disabled");

                // Flip the promote / degrade buttons
                row.querySelector('div.external').classList.remove('d-none');
                row.querySelector('div.external').classList.add('d-display');
                row.querySelector('div.ugent-author').classList.remove('d-display');
                row.querySelector('div.ugent-author').classList.add('d-none');
            })
        })
    }

    // Hook up functions
    htmx.on("ITPromoteModalAfterSettle", promoteAuthor);
    htmx.on("ITAddRowAfterSettle", demoteAuthor);
    htmx.on("ITEditRowAfterSettle", demoteAuthor);

}
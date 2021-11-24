import htmx from 'htmx.org';
import BSN from "bootstrap.native/dist/bootstrap-native-v4";

export default function () {
    let addEvents = function() {
        document.querySelectorAll('input[data-member-choose]').forEach(function(radio) {
            const modalSelector = '#'+radio.dataset.memberChoose
            const modal = new BSN.Modal(modalSelector)

            radio.addEventListener("change", function(evt) {
                if (!this.checked) {
                    return
                }

                modal.show()
            })
        })
    }

    // addEvents()

    // TODO don't use afterSettle
    htmx.on("htmx:afterSettle", function(evt) {
        // addEvents()
    });

    // // Promote an author to "UGent author"
    // let promoteAuthor = function(evt) {
    //     let list = document.querySelector("#modal-contributors-list")
    //     let tableID = list.dataset.contributorList
    //     let rowDelta = list.dataset.row
    //     let row = table.querySelectorAll("tbody tr").item(rowDelta)
    //     let contributors = list.querySelectorAll("li a")

    //     contributors.forEach(function (el) {
    //         el.addEventListener("click", function(evt) {
    //             let item = el.closest("li.list-group-item");

    //             // Set the values on the input fields
    //             row.querySelector("input[name=first_name]").value = item.dataset.firstname;
    //             row.querySelector("input[name=last_name]").value = item.dataset.lastname;
    //             row.querySelector("input[name=ID]").value = item.dataset.id;

    //             // Disable the input fields (name shouldn't be editable if UGent author)
    //             row.querySelector("input[name=first_name]").setAttribute("disabled", "disabled");
    //             row.querySelector("input[name=last_name]").setAttribute("disabled", "disabled");

    //             // Flip the promote / degrade buttons
    //             // row.querySelector('div.external').classList.remove('d-display');
    //             // row.querySelector('div.external').classList.add('d-none');
    //             // row.querySelector('div.ugent-author').classList.remove('d-none');
    //             // row.querySelector('div.ugent-author').classList.add('d-display');

    //             // Close the modal
    //             let modal = document.querySelectorAll(".modal").item(0)
    //             let backdrop = document.querySelectorAll(".modal-backdrop").item(0)

    //             if (modal) {
    //                 modal.classList.remove("show")
    //             }

    //             if (backdrop) {
    //                 backdrop.classList.remove("show")
    //             }

    //             // Timeout gives us a fluid animation
    //             setTimeout(function() {
    //                 if (backdrop) {
    //                     backdrop.remove();
    //                 }

    //                 if (modal) {
    //                     modal.remove();
    //                 }
    //             }, 100)

    //             evt.preventDefault();
    //         }, false);
    //     });
    // }

    // // Demote an UGent author to "external member"
    // let demoteAuthor = function (evt) {
    //     let demoteRadios = document.querySelectorAll("table#contributor-table tbody tr input.demote-external-member");

    //     demoteRadios.forEach(function (el) {
    //         // if (this.checked) {
    //         //     let row = el.closest("tr")
    //         //     // Remove the UGent ID: this equals demoting the user
    //         //     row.querySelector("input[name=ID]").value = "";
    //         //     // Enable the input fields (editable as external member)
    //         //     row.querySelector("input[name=first_name]").removeAttribute("disabled");
    //         //     row.querySelector("input[name=last_name]").removeAttribute("disabled");                
    //         // } else {
    //         //     let row = el.closest("tr")
    //         //     row.querySelector("input[name=first_name]").setAttribute("disabled", true);
    //         //     row.querySelector("input[name=last_name]").setAttribute("disabled", true);                
    //         // }

    //         el.addEventListener("changed", function(evt) {
    //             if (!this.checked) {
    //                 return
    //             }

    //             let row = el.closest("tr")
                
    //             // Remove the UGent ID: this equals demoting the user
    //             row.querySelector("input[name=ID]").value = "";

    //             // Enable the input fields (editable as external member)
    //             row.querySelector("input[name=first_name]").removeAttribute("disabled");
    //             row.querySelector("input[name=last_name]").removeAttribute("disabled");

    //             // Flip the promote / degrade buttons
    //             // row.querySelector('div.external').classList.remove('d-none');
    //             // row.querySelector('div.external').classList.add('d-display');
    //             // row.querySelector('div.ugent-author').classList.remove('d-display');
    //             // row.querySelector('div.ugent-author').classList.add('d-none');
    //         })
    //     })
    // }

    // // Hook up functions
    // htmx.on("ITPromoteModalAfterSettle", promoteAuthor);
    // htmx.on("ITAddRowAfterSettle", demoteAuthor);
    // htmx.on("ITEditRowAfterSettle", demoteAuthor);

}
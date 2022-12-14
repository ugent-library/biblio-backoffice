export default function() {
    const label_select_all = "Select all"
    const label_deselect_all = "Deselect all"

    let formSwitchCheck = function (evt) {
        let btn  = evt.target
        let form = btn.closest("form")
        let chks = form.querySelectorAll("input[type='checkbox']")
        let isSelectAll = btn.textContent == label_select_all
        chks.forEach(el => el.checked = isSelectAll)
        btn.textContent = isSelectAll ? label_deselect_all : label_select_all
    }

    document.querySelectorAll("button.form-check-all").forEach(btn => {

        btn.addEventListener("click", formSwitchCheck)

        //determine startup button value
        let form = btn.closest("form")
        let chks = form.querySelectorAll("input[type='checkbox']")
        let isSelectAll = true
        for(let i = 0;i < chks.length;i++){
            if(chks[i].checked) {
                isSelectAll = false
                break
            }
        }
        btn.textContent = isSelectAll ? label_select_all : label_deselect_all
    })

}

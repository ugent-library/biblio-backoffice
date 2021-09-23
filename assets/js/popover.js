import BSN from 'bootstrap.native';

export default function() {
    // Simple popovers
    let targets = document.querySelectorAll('[data-toggle=popover')
    Array.from(targets).forEach(
        target => new BSN.Popover(target, {})
    )

    // TODO: implement when we do file uploads
    let customTargets = document.querySelectorAll('[data-toggle=popover-custom')
    Array.from(customTargets).forEach(
        target => BSN.Popover(target, {
            content: function() {
                // var content = $(this).attr("data-popover-content");
                // return $(content).children(".popover-body").html();
            },
            title: function() {
                // var content = $(this).attr("data-popover-content");
                // if ($(content).children(".popover-heading").length > 0) {
                //   return $(content).children(".popover-heading").html();
                // } else {
                //   return '';
                // }
            }
        })
    )
}
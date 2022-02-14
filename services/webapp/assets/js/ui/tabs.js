export default function () {
    let tabs = document.querySelectorAll(".bc-toolbar ul a")

    // Set the anchor in the browser nav to the active tab.
    Array.from(tabs).forEach(function (link) {
        link.addEventListener('show.bs.tab', function (evt) {
            window.location.hash = evt.target.hash;
        })
    })

    // Read the hash from browser nav and show the active tab.
    let hash = location.hash.replace(/^#/, "")
    if (hash) {
        let activeTab =  document.querySelector('a[href="#' + hash + '"]');
        activeTab.Tab.show();
    }
}
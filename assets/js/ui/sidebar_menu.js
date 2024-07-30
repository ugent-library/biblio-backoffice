export default function initSidebarMenu() {
  const menuButtons = document.querySelectorAll("[data-sidebar-menu] button");
  const contents = document.querySelectorAll("[data-sidebar-content]");
  const closeButtons = document.querySelectorAll("[data-sidebar-close]");

  menuButtons.forEach((button) => {
    button.addEventListener("click", () => {
      const targetId = button.getAttribute("data-target-id");
      const targetContent = document.getElementById(targetId);

      // Check if the target content is currently displayed.
      if (targetContent.classList.contains("open")) {
        // If so, hide it.
        targetContent.classList.remove("open");
        button.classList.remove("active");
      } else {
        // If not, hide all contents and show the target content.
        contents.forEach((content) => {
          content.classList.remove("open");
        });
        menuButtons.forEach((btn) => {
          btn.classList.remove("active");
        });

        targetContent.classList.add("open");
        button.classList.add("active");
      }
    });
  });

  closeButtons.forEach((button) => {
    button.addEventListener("click", () => {
      const parentContent = button.closest("[data-sidebar-content]");
      if (parentContent) {
        parentContent.classList.remove("open");

        // Find the corresponding menu button and remove active class.
        const correspondingButton = document.querySelector(
          `[data-sidebar-menu] button[data-target="${parentContent.id}"]`,
        );
        if (correspondingButton) {
          correspondingButton.classList.remove("active");
        }
      }
    });
  });
}

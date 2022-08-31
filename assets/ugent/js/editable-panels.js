/* ==========================================================================
    Editable panels
   ========================================================================== */

const panelStateTriggers = document.querySelectorAll('[data-panel-toggle]');

const handleClick = (e) => {
  // Find target Id
  const targetId = e.currentTarget.getAttribute('data-panel-toggle')
  // Find editable content and controls
  const content = document.getElementById(targetId);
  const controls = e.currentTarget.closest('.bc-toolbar').querySelectorAll('.c-button-toolbar[data-panel-state]');

  // Toggle state
  content.querySelectorAll('[data-panel-state]').forEach((el) => (el.classList.toggle('u-hidden')));
  controls.forEach((el) => (el.classList.toggle('u-hidden')));

};

for (let i = 0; i <panelStateTriggers.length; i += 1) {
  // Attach event listeners
  panelStateTriggers[i].addEventListener('click', handleClick, false);
}

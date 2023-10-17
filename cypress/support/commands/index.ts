// Parent commands
import login from './login'
import loginAsResearcher from './login-as-researcher'
import loginAsLibrarian from './login-as-librarian'
import switchMode from './switch-mode'
import ensureModal from './ensure-modal'
import ensureNoModal from './ensure-no-modal'
import visitPublication from './visit-publication'
import ensureToast from './ensure-toast'
import ensureNoToast from './ensure-no-toast'
import setFieldByLabel from './set-field-by-label'

// Child commands
import finishLog from './finish-log'
import closeToast from './close-toast'
import setField from './set-field'

// Dual commands
import extractBiblioId from './extract-biblio-id'
import closeModal from './close-modal'

// Parent commands
Cypress.Commands.addAll({
  login,

  loginAsResearcher,

  loginAsLibrarian,

  switchMode,

  ensureModal,

  ensureNoModal,

  visitPublication,

  ensureToast,

  ensureNoToast,

  setFieldByLabel,
})

// Child commands
Cypress.Commands.addAll(
  { prevSubject: true },
  {
    finishLog,

    closeToast,

    setField,
  }
)

// Dual commands
Cypress.Commands.addAll(
  {
    prevSubject: 'optional',
  },
  {
    extractBiblioId,

    closeModal,
  }
)

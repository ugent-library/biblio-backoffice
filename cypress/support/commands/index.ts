// Parent commands
import ensureModal from './ensure-modal'
import ensureNoModal from './ensure-no-modal'
import finishLog from './finish-log'
import login from './login'
import loginAsLibrarian from './login-as-librarian'
import loginAsResearcher from './login-as-researcher'
import switchMode from './switch-mode'

// Child commands
import visitPublication from './visit-publication'

// Dual commands
import extractBiblioId from './extract-biblio-id'

// Parent commands
Cypress.Commands.addAll({
  login,

  loginAsResearcher,

  loginAsLibrarian,

  switchMode,

  ensureModal,

  ensureNoModal,

  visitPublication,
})

// Child commands
Cypress.Commands.addAll(
  { prevSubject: true },
  {
    finishLog,
  }
)

// Dual commands
Cypress.Commands.addAll(
  {
    prevSubject: 'optional',
  },
  {
    extractBiblioId,
  }
)

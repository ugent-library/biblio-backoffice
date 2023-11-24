const mix = require('laravel-mix')

mix.js('./assets/js/app.js', 'js')
mix.sass('./assets/css/app.scss', 'css')
mix.setPublicPath('./static')

// copy images
mix.copy('assets/images/**/*', 'static/images');
mix.copy('assets/ugent/images/**/*', 'static/images');
mix.copy('assets/ugent/favicon.ico', 'static/favicon.ico');

// copy fonts
mix.copy('assets/ugent/fonts/**', 'static/fonts');
// set the resourceroot for fonts so it points to the static assets path
mix.setResourceRoot('../')

if (mix.inProduction()) {
  mix.version()
} else {
  mix.sourceMaps()
}

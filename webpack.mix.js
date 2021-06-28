const mix = require('laravel-mix')

mix.js('./assets/js/app.js', 'js')
  .sass('./assets/css/app.scss', 'css')
  .setPublicPath('./static')

if (mix.inProduction()) {
  mix.version()
} else {
  mix.sourceMaps()
}

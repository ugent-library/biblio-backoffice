vTempl=github.com/a-h/templ/cmd/templ@v0.2.747
vWgo=github.com/bokwoon95/wgo@latest

live/templ:
	@go run $(vTempl) generate \
	--path views \
	--watch \
	--proxy="http://localhost:3011" \
	--proxyport=3001 \
	--proxybind="localhost" \
	--open-browser=false

live/esbuild:
	@node esbuild.mjs --watch

live/server:
	@BIBLIO_BACKOFFICE_PORT=3011 go run $(vWgo) run \
	-xdir etc -xdir node_modules -xdir assets -xdir cypress \
	-file .go -file static/manifest.json \
	main.go server start

live/sync_assets:
	go run $(vWgo) \
	-xdir etc -xdir app/node_modules -xdir app/assets -xdir cypress \
	-file static/manifest.json \
	go run $(vTempl) generate --notify-proxy --proxyport=3001 --proxybind="localhost"

live:
	@make -j4 live/templ live/esbuild live/server live/sync_assets

dev:
	@go run $(vWgo) run \
	-xdir etc -xdir node_modules -xdir assets -xdir cypress \
	-file .go -file static/manifest.json \
	main.go server start \
	:: wgo -dir assets node esbuild.mjs \
	:: wgo -dir views -file .templ go generate views/generate.go
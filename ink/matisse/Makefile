SEPTEMBER := ../september/bin/main.ink

# compile + output Ink sources
run:
	echo '' > static/ink/bundle.js
	${SEPTEMBER} translate src/canvas.ink | tee /dev/stderr >> static/ink/bundle.js
	${SEPTEMBER} translate src/000.ink | tee /dev/stderr >> static/ink/bundle.js

# compile + output whenever Ink sources change
watch:
	ls src/* | entr make

# run formatter and save
fmt:
	inkfmt fix src/*.ink

fmt-check:
	inkfmt src/*.ink

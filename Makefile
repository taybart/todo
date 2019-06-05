
run:
	mkdir -p client/todo.app/Contents/MacOS
	cd client; go build -o todo.app/Contents/MacOS/todo
	open client/todo.app

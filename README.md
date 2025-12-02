# onionmanager
Portable Tor Onion Service Manager
## Build
	$ git clone https://github.com/aeeravsar/onionmanager
	$ cd onionmanager
	$ go build
## Install
	$ sudo cp onionmanager /usr/local/bin
## Usage
Custom Tor binary (Optional)

	$ TOR_BINARY=<path_to_tor_binary> onionmanager ...
	
Create new Onion service

	$ onionmanager new ~/Documents/myservices/someservice
	Enter virtual port: 80  
	Enter local port: 8080
	Add another port mapping? (y/n): y
	Enter virtual port: 7777
	Enter local port: 5000
	Add another port mapping? (y/n): n
	Generating onion service keys...

	Created new Onion service at: /home/eravsar/Documents/myservices/someservice
	Onion address: xwueyr5r5kmtztyk3uek6wz4zsc637lhgbkftk3sdyk6darvlh6hopqd.onion

	Configuration saved to: /home/eravsar/Documents/myservices/someservice
	Edit manager.conf to customize Tor settings.

	Service ready. Run with: onionmanager run /home/eravsar/Documents/myservices/someservice
	
Edit config (Optional)

	$ cat /home/eravsar/Documents/myservices/someservice/manager.conf
	# Edit torrc options except HiddenServiceDir and DataDirectory here
	
	HiddenServicePort 80 127.0.0.1:8080
	HiddenServicePort 7777 127.0.0.1:5000
	
Run the service

	$ onionmanager run /home/eravsar/Documents/myservices/someservice

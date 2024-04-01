TODO: Move to github.com/issues

test install guide from beginning

e.g. expert page deleting growtable entries does delete, but firefox does not refresh

Remove nearly all javascript, render page for every week, render other page for measurements and yet another for growcontrol,.... -> w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

refactor handler: DRY, move paths to const, 

git pull from script seems to change file ownerships: sudo chown -R pi:pi

find/install git plugin to compare with working tree, enable javascript vscode debug and documentation

Expert: Calibrabrate with select one by one, wait for signal, not 3secs,

get remotly thorugh fritzbox

create REMOVE.sh to uninstall all

Add Docker Image or Makefile to run/test localy for github repo
https://earthly.dev/blog/golang-makefile/

-> set 24h format and date in chrome, bei m stimmt alles

Google js listeners instead of call function, do not use classes and ids, use tags, see todos below:

Thumbnails? first blurry loading?

compare growsniffer/sender with: https://tutorials-raspberrypi.de/raspberry-pis-ueber-433mhz-funk-kommunizieren-lassen/

Remove most of javascript, profile

create growmpage.html version without go template (like expert, handle growganizer as js object)
growcontrol auf growmpage per javascript copie + use form array (same imput names) -> ohne zwischenspeichern grwocontrols und weeks hinzufügen/löschen

download all dependencies(css/js/go) to local folder or extract only necessary code -> rebuild in 10 years possible

divide in growganizer/growtable templates

https://go.dev/doc/articles/wiki/final.go?m=text

evtl. vergangene startdates inkl. wochenname abrufbar machen? -> Kalender farbig markieren? einfaches select "Seeds-29-03.195, ..."

google best practice to refactor, look for prominent go code

check all append() for if slice is nil/empty

maybe only edit raw and remove all other "+, -, ..." in growcontrol and weeks, keep only activate?!

launch.json
{
	"name": "Launch Package",
	"type": "go",
	"request": "launch",
	"mode": "auto",
	"program": "${workspaceRoot}/page"
  },

# TODO: install a watchdog?
# sudo apt-get update && sudo apt-get -y upgrade
# sudo apt install -y watchdog
# TODO: create /home/pi/health in grow.go every 3 seconds
# modprobe bcm2835_wdt
# echo „bcm2835_wdt“ | sudo tee -a /etc/modules
# sudo update-rc.d watchdog defaults
# sudo bash -c 'echo "max-load-1             = 24
# watchdog-device        = /dev/watchdog
# file = /home/pi/health
# change = 700" > /etc/watchdog.conf'
# (crontab -l 2>/dev/null; echo "*/5 * * * * echo ok > /home/pi/health") | crontab -
# sudo service watchdog start


# sudo git config --system --add safe.directory '/home/pi/growmpage'


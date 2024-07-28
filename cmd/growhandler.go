package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/growmpage/growhelper"
)

func (g *Growmpage) handleGrowmpage() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/week#"+growhelper.ToString(g.growganizer.ActiveWeekIndex), http.StatusSeeOther)
	})
	http.HandleFunc("/week", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, g.growcache)
	})
	http.HandleFunc("/Snapshot", func(w http.ResponseWriter, r *http.Request) {
		growhelper.Get("MEASURE") //TODO: use function -> internet-down-protected
		pictureName := growhelper.Get("PICTURE")
		// g.SaveToDatabase()
		http.Redirect(w, r, "/week#"+growhelper.ToString(g.growganizer.ActiveWeekIndex)+"?snapshot="+pictureName, http.StatusSeeOther)
	})
	http.HandleFunc("/Update", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(os.Stdout, "ParseForm() err: %v", err)
			return
		}
		control := r.FormValue("control")
		switch control {
		case "Save": //TODO: normalize control and all other like urls
			g.growganizer.UpdateWeeksByForm(r)
			http.Redirect(w, r, "/", http.StatusSeeOther)
		case "rawUpdateGrowganizer":
			g.growganizer.RawUpdateGrowganizer(r.FormValue("json"))
			http.Redirect(w, r, growhelper.Filename_expertpage, http.StatusSeeOther)
		case "rawUpdateGrowtable":
			g.growtable.RawUpdateGrowtable(r.FormValue("json"), true)
			http.Redirect(w, r, growhelper.Filename_expertpage+"?table=true", http.StatusSeeOther)
		default:
			g.growganizer.UpdateWeek(w, r)
		}
		g.updateGrowtableColors()
		g.SaveToDatabase()
		g.updateGrowmpage()
		g.restartGrowcontroller()
		growhelper.GitPrivateBackup() //TODO: cleanup refreshing things
	})
	http.HandleFunc("/DeleteAlert", func(w http.ResponseWriter, r *http.Request) {
		g.growganizer.DeleteAlert(growhelper.ReadString(r.Body))
		g.updateGrowmpage()
		g.growganizer.SaveToFile()
	})
}

func (g *Growmpage) handleExpertpage() {
	http.HandleFunc("/SIM", func(w http.ResponseWriter, r *http.Request) { //for url-checker too
		fmt.Printf("\"SIM\": %v\n", "SIM")
		go g.growswitcher.SwitchSim(g.growganizer.PlugControls[0], g.growganizer.PlugControls[1])
	})
	http.HandleFunc("/gitPrivateBackup", func(w http.ResponseWriter, r *http.Request) {
		growhelper.GitPrivateBackup()
		fmt.Fprint(w, "backup finished \n")
	})
	http.HandleFunc("/gitPublicUpgrade", func(w http.ResponseWriter, r *http.Request) {
		growhelper.GitPublicUpgrade(w, g.growganizer.Version)
	})
	http.HandleFunc("/gitPublicInstall", func(w http.ResponseWriter, r *http.Request) {
		growhelper.GithubPublicInstall(w, g.growganizer.Version)
	})
	http.HandleFunc("/gitPrivateReset", func(w http.ResponseWriter, r *http.Request) {
		growhelper.GitPrivateReset(w)
		go func() {
			time.Sleep(2 * time.Second)
			os.Exit(0)
		}()
	})
	http.HandleFunc("/gitPrivateVersionLocal", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, growhelper.GitPrivateVersion(false))
	})
	http.HandleFunc("/gitPublicVersionOrigin", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, growhelper.GitPublicVersion())
	})
	http.HandleFunc("/PlugControl", func(w http.ResponseWriter, r *http.Request) { // TODO: also an internal function... -> do not sort/group in methods?!
		plugIndex := g.growganizer.PlugIndex(growhelper.ReadString(r.Body))
		go g.growswitcher.Switch(g.growganizer.PlugControls[plugIndex])
	})
	http.HandleFunc("/Calibrate", func(w http.ResponseWriter, r *http.Request) {
		pinNumberReceive := growhelper.ReadString(r.Body)
		command := exec.Command("bash", "-c", "python3 "+growhelper.Filename_growsniffer+" -g "+pinNumberReceive)
		out, err := command.CombinedOutput()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Fprint(w, string(out))
	})
	http.HandleFunc("/SAVETODATABASE", func(w http.ResponseWriter, r *http.Request) {
		g.SaveToDatabase()
	})
}

func (g *Growmpage) handleInternal() {
	http.HandleFunc("/UPDATE", func(w http.ResponseWriter, r *http.Request) {
		g.updateGrowmpage()
	})
	http.HandleFunc("/RESTARTGROWCONTROLLER", func(w http.ResponseWriter, r *http.Request) {
		g.restartGrowcontroller()
	})
	http.HandleFunc("/ALERT", func(w http.ResponseWriter, r *http.Request) {
		message := growhelper.ReadString(r.Body)
		if message != "" {
			g.growganizer.Alerts = append(g.growganizer.Alerts, message)
		}
		go g.growswitcher.SwitchSim(g.growganizer.PlugControls[0], g.growganizer.PlugControls[1]) //TODO: magic number
	})
	http.HandleFunc("/MEASURE", func(w http.ResponseWriter, r *http.Request) {
		climateObserving := g.growobserver.Measure()
		g.growtable.AddMeasurement(climateObserving.Temperature, climateObserving.Humidity, g.currentWeek()) //TODO: divide function: always meassure on picture!
		g.SaveToDatabase()                                                                                   //TODO: only save to sd if needed, e.g. opening expert page, once a day...
	})
	http.HandleFunc("/PICTURE", func(w http.ResponseWriter, r *http.Request) {
		cameraObserving := g.growobserver.Picture()
		g.growtable.AddPicture(cameraObserving.Minutes, g.currentWeek())
		fmt.Fprint(w, cameraObserving.Minutes)
		g.SaveToDatabase()
	})
	http.HandleFunc("/GROWTABLEREPORT", func(w http.ResponseWriter, r *http.Request) {
		sinceHours := growhelper.ToInt(growhelper.ReadString(r.Body))
		report := g.growtable.Report(sinceHours)
		responseBytes, _ := json.Marshal(report)
		fmt.Fprint(w, string(responseBytes))
	})
	http.HandleFunc("/DAYSSINCEACTIVESTARTDATE", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, g.GetDaysSinceActiveStartDate())
	})

}

func (g *Growmpage) handleFiles() {
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, growhelper.Filename_favicon)
	})
	http.HandleFunc("/html/expert.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, growhelper.Filename_expertpage)
	})
	http.HandleFunc("/data/growganizer.json", func(w http.ResponseWriter, r *http.Request) {
		disableCache(w)
		http.ServeFile(w, r, "../data/growganizer.json")
	})
	http.HandleFunc("/data/growtable.json", func(w http.ResponseWriter, r *http.Request) {
		disableCache(w)
		http.ServeFile(w, r, "../data/growtable.json")
	})
	http.Handle("/data/camera/", http.StripPrefix("/data/camera/", http.FileServer(http.Dir("../data/camera"))))
}

func disableCache(w http.ResponseWriter) {
	w.Header().Add("X-Frame-Options", "DENY")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
}

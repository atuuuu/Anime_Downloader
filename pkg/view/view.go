package view

import (
	"../downloader_utils"
	"fmt"
	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
	"log"
	"os"
	"strings"
	"sync"
	"syscall"
	"time"
)

const (
	EsSystemRequired  = 0x00000001
	EsDisplayRequired = 0x00000002
)

var wg = sync.WaitGroup{}

var w *astilectron.Window

var L = log.New(log.Writer(), log.Prefix(), log.Flags())

var App = bootstrap.Options{
	AstilectronOptions: astilectron.Options{
		AppIconDarwinPath: "resources/app/icon.ico",
	},
	Debug:  false,
	Logger: L,
	OnWait: func(_ *astilectron.Astilectron, ws []*astilectron.Window, _ *astilectron.Menu, _ *astilectron.Tray, _ *astilectron.Menu) error {
		w = ws[0]

		log.Println("Starting keep alive poll... (silence)")

		w.On(astilectron.EventNameAppClose, func(e astilectron.Event) (deleteListener bool) {
			wg.Wait()
			os.Exit(0)
			return true
		})

		w.OnMessage(handleMessage)

		go noSleep()

		return nil
	},
	Windows: []*bootstrap.Window{{
		Homepage:       "view.html",
		MessageHandler: nil,
		Options: &astilectron.WindowOptions{
			BackgroundColor: astikit.StrPtr("#FFF"),
			Center:          astikit.BoolPtr(true),
			Height:          astikit.IntPtr(600),
			Width:           astikit.IntPtr(800),
		},
	}},
}

func noSleep() {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	setThreadExecStateProc := kernel32.NewProc("SetThreadExecutionState")

	for {
		time.Sleep(40 * time.Second)
		setThreadExecStateProc.Call(uintptr(EsSystemRequired))
		setThreadExecStateProc.Call(uintptr(EsDisplayRequired))
		fmt.Println("No sleep !")

	}
}

func handleMessage(m *astilectron.EventMessage) interface{} {
	// Unmarshal
	var s string
	if err := m.Unmarshal(&s); err != nil {
		L.Println(err)
	}
	var message = strings.Split(s, "<")
	// Process message
	if message[0] == "download" {
		for i := 1; i < len(message); i++ {
			if message[i] != " " && message[i] != "" && message[i] != "undefined" {
				fmt.Println("lien : " + message[i])
				wg.Add(1)
				go downloader_utils.InitDownload(message[i], "./download/"+message[i]+".mp4", w, i)
			}
		}
		go trackDownload()
	}
	if message[0] == "end" {
		wg.Done()
		os.Exit(0)
	}
	return nil
}

func trackDownload() {
	wg.Wait()
}

var SourceHTML = `<!doctype html>
<html lang="fr">

<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Anime downloader 2nd Edition Ver. 2.31</title>
  <style>
	a.button1{
		display:inline-block;
		padding:0.35em 1.2em;
		border:0.1em solid #FFFFFF;
		margin:0 0.3em 0.3em 0;
		border-radius:0.12em;
		box-sizing: border-box;
		text-decoration:none;
		font-family:'Roboto',sans-serif;
		font-weight:300;
		background-color:#581D15;
		text-align:center;
		transition: all 0.2s;
		cursor: pointer;
	}
	a.button1:hover{
		font-size: 140%;
		color:#FFFFFF;
		background-color: transparent;
	}
	@media all and (max-width:30em){
		a.button1{
			display:block;
			margin:0.4em auto;
		}
	}

	.input {
	  font-size: 16px;
	  font-size: max(16px, 1em);
	  font-family: inherit;
	  background-color: #fff;
	  border-radius: 5px;
	  margin-left: 2%;
	  width: 30%;
	  resize: none;
	}
	progress.green{
	  background-color:green;
	  height: 10px;
	}
	progress.green::-webkit-progress-value {background-color:green}

	progress.red{
	  background-color:red;
	  height: 10px;
	}
	progress.red::-webkit-progress-value {background-color:red}
  </style>
</head>

<body style="background-color: #682420;">
	<div style="border: solid; border-color: #E4B8C5; background-color: #9E372B;>
		<header style="margin-left: 2%;">
		    <h3>Anime Downloader</h3>
   			<p>Entrez ici vos liens vers les épisodes</p>
		</header>

		<div id="main" role="main" class="main" style="margin-left: 2%; background-image: url(https://media.senscritique.com/media/000019754762/960/Les_meilleurs_animes_japonais.jpg); background-repeat: no-repeat;">
			<aside style="float: right; margin-right: 40%; margin-top: 15%">
				<a class="button1" id="go" onclick="startDownload()" style="height: 50px; width: 350%">Go !</a>
			</aside>
			<div>
				<progress id="p1" value="0" max="100" style="margin-left: 2%;"></progress><br/>
				<textarea id="lien1" class="input"></textarea>
			</div>
			<div>
				<progress id="p2" value="0" max="100" style="margin-left: 2%;"></progress><br/>
				<textarea id="lien2" class="input"></textarea>
			</div>
			<div>
				<progress id="p3" value="0" max="100" style="margin-left: 2%;"></progress><br/>
				<textarea id="lien3" class="input"></textarea>
			</div>
			<div>
				<progress id="p4" value="0" max="100" style="margin-left: 2%;"></progress><br/>
				<textarea id="lien4" class="input"></textarea>
			</div>
		</div>

		<div style="margin-left: 2%; font-size: 80%;">
			<aside style="float: right; margin-right: 5%">
				<!-- pub ici si besoin-->
			</aside>
			<p>Réalisé par Atu</p>
			<p>Discord : Atu#2765</p>
			<p>Mail : arthur.boyreau@gmail.com</p>
			<p>Pour vos projets web/application/bot, ou en cas de bug, n'hésitez pas à me contacter</p>
			<p>Tarifs à négocier</p>
		</div>
	</div>
</body>

    <script pkg="index.js"></script>

</html>`

var SourceJS = `// This will wait for the astilectron namespace to be ready
document.addEventListener('astilectron-ready', function() {
    // This will listen to messages sent by GO
    astilectron.onMessage(function(message) {
        // Process message
		var result = message.split('>');
        if (result[0] == "progress") {
			var val = parseFloat(result[2]);
			val = val * 100;
			document.getElementById('p'+result[1]).value = val;
        }
		if (result[0] == "error") {
			alert("Lien " + result[1] + " : téléchargement échoué\n" + result[2])
			document.getElementById('p'+result[1]).classList.add("red");
			document.getElementById('p'+result[1]).value = 50;
			document.getElementById('lien'+result[1]).disabled = false;
		}
		if (result[0] == "success") {
			alert("Lien " + result[1] + " : téléchargement terminé")
			document.getElementById('p'+result[1]).classList.add("green");
			document.getElementById('lien'+ result[1]).disabled = false;
			document.getElementById('lien'+ result[1]).value = " ";
		}
    });
})

function startDownload() {
	var link = [];
	for(i = 1; i < 5; i++) {
		if(!document.getElementById('lien'+i).disabled) {
			link[i-1] = document.getElementById('lien'+i).value;
			if(link[i-1] != undefined && link[i-1] != "" && link[i-1] != " ") {
				document.getElementById('lien'+i).setAttribute("disabled", true);
				document.getElementById('p'+ i).classList.remove("red");
				document.getElementById('p'+ i).classList.remove("green");
				document.getElementById('p'+ i).value = 0;
			}
		}
	}

    astilectron.sendMessage("download<" + link[0] + "<" + link[1] + "<" + link[2] + "<" + link[3], function(message) {
        console.log("envoyé " + message);
    });
	if(link[0] != "" || link[1] != "" || link[2] != "" || link[3] != "") {
		alert("Téléchargement lancé");
	}
}`

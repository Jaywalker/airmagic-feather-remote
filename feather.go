package main

import "github.com/nsf/termbox-go"
import "net/http"
import "time"
import "flag"

type key struct {
	x  int
	y  int
	ch rune
}

var DeviceIP string
var DevicePort string

// Our arrow and select buttons
var K_ARROW_LEFT = []key{{2, 3, '('}, {3, 3, 0x2190}, {4, 3, ')'}}
var K_ARROW_UP = []key{{6, 1, '('}, {7, 1, 0x2191}, {8, 1, ')'}}
var K_ARROW_DOWN = []key{{6, 5, '('}, {7, 5, 0x2193}, {8, 5, ')'}}
var K_ARROW_RIGHT = []key{{10, 3, '('}, {11, 3, 0x2192}, {12, 3, ')'}}
var K_SELECT = []key{{6, 2, 0x2591}, {7, 2, 0x2591}, {8, 2, 0x2591},
	{5, 3, 0x2591}, {6, 3, 0x2591}, {7, 3, 0x21B5}, {8, 3, 0x2591}, {9, 3, 0x2591},
	{6, 4, 0x2591}, {7, 4, 0x2591}, {8, 4, 0x2591}}

//TV/Menu buttons
var K_MENU = []key{{2, 8, '('}, {3, 8, 'M'}, {4, 8, ')'}}
var K_TV = []key{{10, 8, '('}, {11, 8, 'T' /* `📺` 0x1F4FA */}, {12, 8, ')'}}

//Siri/PlayPause buttons
var K_SIRI = []key{{2, 12, '('}, {3, 12, 'S' /* 0x1F399 */}, {4, 12, ')'}}
var K_PLAYPAUSE = []key{{1, 18, '('}, {2, 18, 0x23F5}, {3, 18, 0x23F8}, {4, 18, ')'}} // Looks better than the single 0x23EF, i think

//Volume buttons
var K_VOLUME_TOP = []key{{10, 12, 0x256D}, {11, 12, 0x2500}, {12, 12, 0x256E},
	{10, 13, 0x2502}, {11, 13, '+'}, {12, 13, 0x2502},
	{10, 14, 0x2502}, {11, 14, 0x2591}, {12, 14, 0x2502}}

var K_VOLUME_MID = []key{{10, 15, 0x2502}, {11, 15, 0x2591}, {12, 15, 0x2502}}

var K_VOLUME_BOTTOM = []key{{10, 16, 0x2502}, {11, 16, 0x2591}, {12, 16, 0x2502},
	{10, 17, 0x2502}, {11, 17, '-'}, {12, 17, 0x2502},
	{10, 18, 0x2570}, {11, 18, 0x2500}, {12, 18, 0x256F}}

const (
	ARROW_LEFT = iota
	ARROW_UP
	ARROW_DOWN
	ARROW_RIGHT
	SELECT
	MENU
	TV
	SIRI
	PLAYPAUSE
	VOLUME_UP
	VOLUME_DOWN
)

var currently_pressed = -1

func draw_key(k []key, fg, bg termbox.Attribute) {
	w, h := termbox.Size()
	w /= 2
	w -= 9
	h /= 2
	h -= 14
	for _, k := range k {
		termbox.SetCell(k.x+w+2, k.y+h+2, k.ch, fg, bg)
	}
}

func draw_remote() {
	w, h := termbox.Size()
	w /= 2
	w -= 9
	h /= 2
	h -= 14
	//Draw our corners
	termbox.SetCell(0+w, 0+h, 0x256D /* 0x250C*/, termbox.ColorWhite, termbox.ColorBlack)    //top left
	termbox.SetCell(18+w, 0+h, 0x256E /* 0x2510*/, termbox.ColorWhite, termbox.ColorBlack)   //top right
	termbox.SetCell(0+w, 27+h, 0x2570 /* 0x2514 */, termbox.ColorWhite, termbox.ColorBlack)  //bottom left
	termbox.SetCell(18+w, 27+h, 0x256F /* 0x2518 */, termbox.ColorWhite, termbox.ColorBlack) //bottom right

	//Now draw our top and bottom lines
	for i := 1; i < 18; i++ {
		termbox.SetCell(i+w, 0+h, 0x2500, termbox.ColorWhite, termbox.ColorBlack)
		termbox.SetCell(i+w, 27+h, 0x2500, termbox.ColorWhite, termbox.ColorBlack)
	}

	//Now draw our side lines
	for i := 1; i < 27; i++ {
		termbox.SetCell(0+w, i+h, 0x2502, termbox.ColorWhite, termbox.ColorBlack)
		termbox.SetCell(18+w, i+h, 0x2502, termbox.ColorWhite, termbox.ColorBlack)
	}

	//Last, draw our keys
	if currently_pressed == ARROW_LEFT {
		draw_key(K_ARROW_LEFT, termbox.ColorWhite, termbox.ColorBlue)
	} else {
		draw_key(K_ARROW_LEFT, termbox.ColorBlack, termbox.ColorWhite)
	}
	if currently_pressed == ARROW_UP {
		draw_key(K_ARROW_UP, termbox.ColorWhite, termbox.ColorBlue)
	} else {
		draw_key(K_ARROW_UP, termbox.ColorBlack, termbox.ColorWhite)
	}
	if currently_pressed == ARROW_DOWN {
		draw_key(K_ARROW_DOWN, termbox.ColorWhite, termbox.ColorBlue)
	} else {
		draw_key(K_ARROW_DOWN, termbox.ColorBlack, termbox.ColorWhite)
	}
	if currently_pressed == ARROW_RIGHT {
		draw_key(K_ARROW_RIGHT, termbox.ColorWhite, termbox.ColorBlue)
	} else {
		draw_key(K_ARROW_RIGHT, termbox.ColorBlack, termbox.ColorWhite)
	}
	if currently_pressed == SELECT {
		draw_key(K_SELECT, termbox.ColorWhite, termbox.ColorBlue)
	} else {
		draw_key(K_SELECT, termbox.ColorBlack, termbox.ColorWhite)
	}
	if currently_pressed == MENU {
		draw_key(K_MENU, termbox.ColorWhite, termbox.ColorBlue)
	} else {
		draw_key(K_MENU, termbox.ColorBlack, termbox.ColorWhite)
	}
	if currently_pressed == TV {
		draw_key(K_TV, termbox.ColorWhite, termbox.ColorBlue)
	} else {
		draw_key(K_TV, termbox.ColorBlack, termbox.ColorWhite)
	}
	if currently_pressed == SIRI {
		draw_key(K_SIRI, termbox.ColorWhite, termbox.ColorBlue)
	} else {
		draw_key(K_SIRI, termbox.ColorBlack, termbox.ColorWhite)
	}
	if currently_pressed == PLAYPAUSE {
		draw_key(K_PLAYPAUSE, termbox.ColorWhite, termbox.ColorBlue)
	} else {
		draw_key(K_PLAYPAUSE, termbox.ColorBlack, termbox.ColorWhite)
	}
	if currently_pressed == VOLUME_UP {
		draw_key(K_VOLUME_TOP, termbox.ColorWhite, termbox.ColorBlue)
	} else {
		draw_key(K_VOLUME_TOP, termbox.ColorBlack, termbox.ColorWhite)
	}
	draw_key(K_VOLUME_MID, termbox.ColorBlack, termbox.ColorWhite)
	if currently_pressed == VOLUME_DOWN {
		draw_key(K_VOLUME_BOTTOM, termbox.ColorWhite, termbox.ColorBlue)
	} else {
		draw_key(K_VOLUME_BOTTOM, termbox.ColorBlack, termbox.ColorWhite)
	}
}

func remote_command(cmd string) {
	//http://10.0.0.3/remoteCommand=left
	//enterText    =    sends text to a text entry view (ie enterText='my search')
	//remoteCommand    =    sends a remote command to the AppleTV (menu, up, down, play, left, right, select)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://"+DeviceIP+":"+DevicePort+"/remoteCommand="+cmd, nil)
	// q := req.URL.Query()
	// q.Add("remoteCommand", cmd)
	// req.URL.RawQuery = q.Encode()
	//resp, err := client.Do(req)
	client.Do(req)
}

func dispatch_press(ev *termbox.Event) {
	if ev.Ch == 'n' || ev.Ch == 'j' || ev.Key == termbox.KeyArrowDown {
		currently_pressed = ARROW_DOWN
		remote_command("down")
	} else if ev.Ch == 'e' || ev.Ch == 'k' || ev.Key == termbox.KeyArrowUp {
		currently_pressed = ARROW_UP
		remote_command("up")
	} else if ev.Ch == 'o' || ev.Ch == 'l' || ev.Key == termbox.KeyArrowRight {
		currently_pressed = ARROW_RIGHT
		remote_command("right")
	} else if ev.Ch == 'y' || ev.Ch == 'h' || ev.Key == termbox.KeyArrowLeft {
		currently_pressed = ARROW_LEFT
		remote_command("left")
	} else if ev.Ch == 'm' {
		currently_pressed = MENU
		remote_command("menu")
	} else if ev.Ch == 't' {
		currently_pressed = TV
		//TODO: remote_command("tv")
	} else if ev.Ch == 's' {
		currently_pressed = SIRI
		//TODO: remote_command("siri")
	} else if ev.Ch == 'u' {
		currently_pressed = VOLUME_UP
		//TODO: remote_command("vol_up")
	} else if ev.Ch == 'd' {
		currently_pressed = VOLUME_DOWN
		//TODO: remote_command("vol_down")
	} else if ev.Ch == 'p' {
		currently_pressed = PLAYPAUSE
		remote_command("play")
	} else if ev.Key == termbox.KeyEnter {
		currently_pressed = SELECT
		remote_command("select")
	}
}

func redraw() {
	//This is a dirty hax for making the buttons not appear sticky. Works for me. If you don't like the flicker, you deal with it.
	for {
		time.Sleep(time.Millisecond * 500)
		currently_pressed = -1
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		draw_remote()
		termbox.Flush()
	}
}

func main() {
	flag.StringVar(&DeviceIP, "ip", "127.0.0.1", "The IP of the device to control")
	flag.StringVar(&DevicePort, "port", "3073", "The TCP port to use on the device to control")
	flag.Parse()
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)

	go redraw()
loop:
	for {
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		draw_remote()
		termbox.Flush()
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyCtrlQ {
				break loop
			}
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			draw_remote()
			dispatch_press(&ev)
			termbox.Flush()
		case termbox.EventResize:
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			draw_remote()
			termbox.Flush()
			/*
				case termbox.EventMouse:
					termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
					draw_remote()
					termbox.Flush()
			*/
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}

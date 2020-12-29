# Sox-WebInterface
A quick web interface to start and stop recording using SOX

I need a way to record sound from a Vinyl record player to a headless Raspberry PI.
This Go Lang code provides a very simple web interface to start and stop recording using SOX audio tool.

I'm using a Raspberry Pi Zero and a cheap 7.1 Sound adapter with my Vinyl record player plugged into the Line Input.
With ALSA the default is Mic input. To check this do a amixer contents:
numid=16,iface=MIXER,name='PCM Capture Source'
  ; type=ENUMERATED,access=rw------,values=1,items=4
  ; Item #0 'Mic'
  ; Item #1 'Line'
  ; Item #2 'IEC958 In'
  ; Item #3 'Mixer'
  : values=0
  
 To change this for my card which is the second card according to ALSA:
 amixer -c 2 cset numid=16 1
 
To run the webserver use:
go run main.go

And goto http://localhost:8081

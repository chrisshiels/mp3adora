# 'Makefile'.
# Chris Shiels.


all: file.wav file.mp3 fileid3v12.mp3


file.wav:
	sox --rate 44100 --bits 24 --channels 2 -n file.wav synth 3 sine 300-600


file.mp3: file.wav
	lame -h file.wav file.mp3


fileid3v12.mp3: file.mp3
	cp file.mp3 fileid3v12.mp3
	id3tag \
		--artist=Artist \
		--album=Album \
		--song=Song \
		--comment=Comment \
		--year=1970 \
		--track=1 \
		fileid3v12.mp3


clean:
	rm -f file.wav file.mp3 fileid3v12.mp3

# goMediaServer
A simple http-flv media server by golang

publish:
ffmpeg.exe -re -i cuc_ieschool.flv -c copy -f flv  http://localhost:8888/live/cuc_ieschool.flv

play:
ffplay.exe http://localhost:8888/live/cuc_ieschool.flv

reference:
go-media-server:https://github.com/songshenyi/go-media-server
go-oryx: https://github.com/ossrs/go-oryx

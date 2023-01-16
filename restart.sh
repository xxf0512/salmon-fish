cd fixtures && docker-compose start -d
cd ..
nohup ./salmon-fish >> log &
cd explorer && docker-compose start -d
cd ..
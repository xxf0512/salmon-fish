cd fixtures && docker-compose start -d
cd ..
nohup ./salmon-fish >> log &
cd exploer && docker-compose start -d
cd ..
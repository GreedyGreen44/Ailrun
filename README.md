# Overview
Pretty straightforward app to get info from [globe.adsbexchange.com](https://globe.adsbexchange.com/)

For now supports only getting aircrafts information from defined map box and saving them to csv file periodically

## Config file example

Here is example of config file used by application

```
#Essential fields
Host=globe.adsbexchange.com
BoxBot=33.521887
BoxTop=40.672535
BoxLeft=-92.697437
BoxRight=-63.154660
OutputType=csv
OutputDirectory=/home/User/AilrunOutput
TimerValueSecs=300
#Optional fields
UserAgent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36
```

## Docker image

Additionally you can pull a docker image:
```
docker pull greedygreen44/ailrun
```

Note, that running inside a container requires sharing local directories with it, like this:

```bash
docker run -t -v local/directory/with/configs:/Ailrun/configs:ro -v local/output/directory:Ailrun/output:rw greedygreen44/ailrun /Ailrun/configs/config.txt
```
Also you have to change OutputDirectory parameter in config file to match container internal output directory:

```
OutputDirectory=/Ailrun/output
```

## Dependencies
For decompressing zstd response [klauspost/compress](https://github.com/klauspost/compress/tree/master/zstd) is used

## Note

Note, that according to adsbExchange [Terms and Conditions](https://www.adsbexchange.com/legal-and-privacy/) commercial (for profit or non-profit organization) use requires written authorization from ADS-B Exchange.
You can read about other conditions of using provided data [here](https://rapidapi.com/adsbx/api/adsbexchange-com1)
